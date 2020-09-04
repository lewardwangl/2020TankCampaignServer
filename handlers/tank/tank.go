// 这个文件主要是 坦克大战使用的接口处理逻辑
package tank

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"gorm.io/gorm"
	"net"
	"net/http"
	"server/config"
	db "server/data"
	"server/utils"
	"strconv"
	"sync"
	"time"
)

// db都采用了debug  因为直接打了日志 数据不多，为了避免无日志可寻

// Ranks 获取坦克大战车轮战排名数据
func Ranks(c *gin.Context) {
	c.JSON(http.StatusOK, db.QueryRanks())
}

// CreateBattleInfo 被用于车轮战，客户端上传上来的每场比赛数据 直接入库， 不做更改 组委会意图
func CreateBattleInfo(c *gin.Context) {
	var params db.Tank
	err := c.ShouldBindJSON(&params)
	c.Set(config.RequestLogParamKey, utils.Struct2str(params))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r := db.DB.Debug().Model(db.Tank{}).Where("battle_id=? AND invalid=0", params.BattleId).Find(&struct{}{})
	if r.RowsAffected != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "当前已经有一场id为" + strconv.Itoa(*params.BattleId) + "且有效的比赛了，请重试！",
		})
		return
	}
	db.DB.Debug().Model(db.Tank{}).Create(&params)
	boardCast(WheelWarMode)
	c.JSON(http.StatusOK, gin.H{
		"msg": "it's ok!",
	})
}

// InvalidBattleById 用于裁判宣布某场比赛无效 分为车轮战和决赛
func InvalidBattleById(c *gin.Context) {
	var params struct {
		Id      *int   `json:"id" binding:"required"`
		Invalid bool   `json:"invalid"`
		Type    string `json:"type" binding:"oneof='wheel_war' 'final'"`
	}
	params.Invalid = true
	err := c.ShouldBindJSON(&params)
	c.Set(config.RequestLogParamKey, utils.Struct2str(params))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var d = db.DB.Debug()
	if params.Type == FinalMode {
		d = d.Model(&db.TankFinal{})
	} else {
		d = d.Model(&db.Tank{})
	}
	d.Where("battle_id = ?", params.Id).Update("invalid", params.Invalid)
	boardCast(params.Type)
	c.JSON(http.StatusOK, gin.H{
		"msg": "it's ok!",
	})
}

// CreateFinalBattleInfo 创建用于决赛的比赛数据  直接入库， 不做更改 组委会意图
func CreateFinalBattleInfo(c *gin.Context) {
	var params db.TankFinal
	err := c.ShouldBindJSON(&params)
	c.Set(config.RequestLogParamKey, utils.Struct2str(params))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r := db.DB.Debug().Model(db.TankFinal{}).Where("battle_id=? AND invalid=0", params.BattleId).Find(&struct{}{})
	if r.RowsAffected != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "当前已经有一场id为" + strconv.Itoa(*params.BattleId) + "且有效的比赛了，请重试！",
		})
		return
	}
	db.DB.Debug().Model(&db.TankFinal{}).Create(&params)
	boardCast(FinalMode)
	c.JSON(http.StatusOK, gin.H{
		"msg": "it's ok!",
	})
}

// FinaleRanks 获取决赛的排名数据
func FinalRanks(c *gin.Context) {
	var result = []struct {
		Win   string `json:"player"`
		Count int    `json:"count"`
	}{}

	subQuery := db.DB.Table("(? UNION ALL ?) AS a1",
		db.DB.Model(&db.TankFinal{}).Select([]string{
			"player_a AS player",
			"CASE WHEN win = player_a THEN 1 ELSE 0 END as num",
		}).Where("invalid = 0 AND battle_id < 1000"),
		db.DB.Model(&db.TankFinal{}).Select([]string{
			"player_b AS player",
			"CASE WHEN win = player_b THEN 1 ELSE 0 END as num",
		}).Where("invalid = 0 AND battle_id < 1000"),
	)
	db.DB.Debug().Model(&db.TankFinal{}).Table("(?) as o", subQuery).Select([]string{
		"player AS win",
		"SUM(o.num) AS count",
	}).Group("player").Order("count DESC").Scan(&result)
	c.JSON(http.StatusOK, result)
}

// HavingBattleId 查询battleId 用于龙神的游戏上传数据前先查询是否有这个battle_id的存在
func HavingBattleId(c *gin.Context) {
	mode := c.Query("type")
	id, err := strconv.Atoi(c.Query("battle_id"))
	c.Set(config.RequestLogParamKey, "mode:"+mode+",id:"+c.Query("battle_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "龙神 battle_id 呢？是不是数字啊 ☹️",
		})
		return
	}
	var result *gorm.DB
	if mode == WheelWarMode {
		result = db.DB.Debug().Model(db.Tank{}).Where("battle_id=? AND invalid=0", id).Find(&struct{}{})
	} else if mode == FinalMode {
		result = db.DB.Debug().Model(db.TankFinal{}).Where("battle_id=? AND invalid=0", id).First(&struct{}{})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "龙神，你是不是又想做坏事儿？不传比赛类型，就想骗我。",
		})
		return
	}
	fmt.Println(result.RowsAffected)
	c.JSON(http.StatusOK, gin.H{
		"isExist": result.RowsAffected != 0,
	})
}

// 以下是实时推送逻辑： 用于比赛数据上报后，用ws通过前端数据发生了变化。 前端可行再请求获取数据，这边不直接返回数据了
var connections = []map[net.Conn]bool{
	make(map[net.Conn]bool), // 车轮战连接
	make(map[net.Conn]bool), // 决赛连接
}
var HttpUpgrader = &ws.HTTPUpgrader{Timeout: time.Second}
var WheelWarMode = "wheel_war"
var FinalMode = "final"
var WheelWarModeValue = 0
var FinalModeValue = 1
var modeM = map[string]int{
	WheelWarMode: WheelWarModeValue,
	FinalMode:    FinalModeValue,
}
var lock sync.Mutex

// WebSocketRanks ws处理 使用了ws库，自行查阅文档吧，用的最基本的小功能
func WebSocketRanks(c *gin.Context) {
	mode := c.Query("mode") //wheel_war final
	if mode != WheelWarMode && mode != FinalMode {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "传个type呗，用于标记是车轮战还是决赛哦, ==> " + WheelWarMode + " | " + FinalMode,
		})
		return
	}
	conn, _, _, err := HttpUpgrader.Upgrade(c.Request, c.Writer)
	if err != nil {
		utils.Error.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	lock.Lock()
	connections[modeM[mode]][conn] = true
	lock.Unlock()

	go func() {
		defer func() {
			lock.Lock()
			delete(connections[modeM[mode]], conn)
			lock.Unlock()
			_ = conn.Close()
		}()

		for {
			err := wsutil.WriteServerMessage(conn, ws.OpPing, []byte("ping"))
			if err != nil {
				break
			}
			time.Sleep(time.Second * 3)
		}
	}()
}

// boardCast 数据发生变更，广播出去
func boardCast(mode string) {
	if mode != WheelWarMode && mode != FinalMode {
		return
	}
	index := modeM[mode]
	for conn, _ := range connections[index] {
		_ = wsutil.WriteServerText(conn, []byte("排名改变啦啦啦啦")) // 固定成这个字符串了，前端收到就重取数据
	}
	utils.Info.Println(mode + "广播成功")
}
