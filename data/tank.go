package db

import "time"

// 坦克大战model

// 车轮战比赛数据
type Tank struct {
	BattleId  *int    `json:"battle_id" binding:"required"`
	PlayerA   *string `json:"playerA" binding:"required"`
	PlayerB   *string `json:"playerB" binding:"required"`
	ScoreA    *int    `json:"scoreA" binding:"required,gte=0"`
	ScoreB    *int    `json:"scoreB" binding:"required,gte=0"`
	Invalid   bool    `json:"invalid"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// 决赛数据表
type TankFinal struct {
	BattleId  *int    `json:"battle_id" binding:"required"`
	PlayerA   *string `json:"playerA" binding:"required"`
	PlayerB   *string `json:"playerB" binding:"required"`
	Win       *string `json:"win" binding:"required"`
	Invalid   bool    `json:"invalid"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func init() {
	_ = DB.AutoMigrate(&Tank{})
	_ = DB.AutoMigrate(&TankFinal{})
}

type TankRank struct {
	Rank   int    `json:"rank"`
	Player string `json:"player"`
	Score  int    `json:"score"`
	Num    int    `json:"num"`
}

// 获取车轮战排名数据
func QueryRanks() (result []TankRank) {
	result = []TankRank{}
	subQuery := DB.Model(Tank{}).Table("(? UNION ALL ?) as other",
		DB.Model(&Tank{}).Select([]string{
			"player_a as player",
			"score_a as score",
			"case when score_a > score_b then 1 else 0 end as num",
		}).Where("invalid = 0 AND battle_id < 1000"),
		DB.Model(&Tank{}).Select([]string{
			"player_b as player",
			"score_b as score",
			"case when score_b > score_a then 1 else 0 end as num",
		}).Where("invalid = 0 AND battle_id < 1000")).
		Select([]string{
			"player",
			"SUM(score) as score",
			"SUM(num) as num",
		}).
		Group("player").Order("score DESC").Order("num DESC")

	DB.Debug().Model(Tank{}).Table("(?) as other1, (select @n:=0) as other2", subQuery).Select([]string{
		"@n:=@n+1 rank",
		"player",
		"score",
		"num",
	}).Scan(&result)
	return result
}
