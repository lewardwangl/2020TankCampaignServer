# 1024服务说明
1. 测试地址： `10.11.159.156:7834`
2. 使用框架： go-gin、gorm、ws
3. 联系方式： 大数据中心前端切图仔[ewardwang@sohu-inc.com]

# 目录说明
```textmate
.
├── README.md
├── config 项目配置文件夹
│   └── config.go 项目配置
├── data 数据库相关
│   ├── db.go 数据库的初始化工作
│   ├── leetcode.go 力扣大赛的数据库model 设置
│   └── tank.go  坦克大战的数据库model 设置
├── go.mod
├── go.sum
├── handlers  请求处理函数
│   ├── leetcode
│   │   └── leetcode.go 力扣大赛
│   └── tank
│       └── tank.go  坦克大战的
├── log  本地测试环境下：测试日志
│   ├── errors.log 发生错误 
│   ├── info.log 日志
│   ├── req.log 请求日志，在中间件中打印的
│   ├── sql.log 数据库操作日志，db打印的
│   └── warning.log 警告日志
├── main.go
├── middleware 中间件
│   ├── auth.go 认证
│   ├── cors.go 跨域
│   └── log.go  日志
├── public  测试用的。。。
│   └── index.html
├── release.sh  部署脚本
└── utils  帮助方法
    ├── log.go 日志方法
    ├── openFile.go  打开文件，无就创建
    └── structToString.go 转换结构
```

# 部署【~~仅仅为了方便我~~ 】
1. 部署地址： `10.11.159.156`  
2. 部署方式： 执行`release.sh`脚本，按提示输入所需的所有的环境变量

<span style="color:red;font-weight:bold">注意：本地启动时，请注意要配置好所有的环境变量，具体所需环境变量请查看config文件</span>

# 接口列表

1. 获取力扣大赛的排名数据: `GET /api/leetcode/ranks`

    ```js
    // response exp:
    [{
      "rank": 1,
      "name": "yipengzhang217537",
      "email": "yipengzhang217537@sohu-inc.com",
      "score": 1356,
      "bout": "B场",
      "bonus": 300
    }]
    ```

2. 获取坦克大战排名数据：`GET /api/tank/ranks`
    ```js
    // response exp: 
    [
        {
            "name": "22",
            "score": 690
        },
        {
            "name": "测试A",
            "score": 0
        }
    ]
    ```
    
3. 获取坦克大战决赛排名数据: `GET /api/tank/ranks/final`
    ```js
   // response exp:
   [
       {
           "player": "测试B",
           "count": 4
       },
       {
           "player": "测试A",
           "count": 2
       },
       {
           "player": "test2",
           "count": 2
       },
       {
           "player": "test1",
           "count": 1
       }
   ]
   ```

4. 上传当前场次的比赛信息 `POST /api/tank/rank?_tc={...}`  Request: 

   |   参数    |  类型  |       说明        | 默认  |
   | :-------: | :----: | :---------------: | :------: |
   | battle_id |  int   |   比赛场次号码    |  无   |
   |  playerA  | string |     A选手名称     |  无   |
   |  scoreA   |  int   | A选手当前场次分数 |  无   |
   |  playerB  | string |     B选手名称     |  无   |
   |  scoreB   |  int   | B选手当前场次分数 |  无   |
   |  invalid  |  bool  | 当前场次是否无效  | false |

5. 上传当前场次比赛信息，决赛使用 `POST /api/tank/rank/final?_tc={...}`  Request:

   |   参数    |  类型  |       说明       | 默认  |
   | :-------: | :---------------: | :--------------: | :---------------: |
   | battle_id |  int   |   比赛场次号码   |  无   |
   |  playerA  | string |    A选手名称     |  无   |
   |  playerB  | string |    B选手名称     |  无   |
   |  invalid  |  bool  | 当前场次是否无效 | false |
   |    win    | string |  当前场次谁赢了  |  无   |

6. 宣布某场次比赛无效 `PATCH /api/tank/rank/invalid?_tc={...}`  Request:

   |  参数   | 类型 |   说明   | 默认 |
   | :-----: | :---------------: | :------: | :---------------: |
   |   id    | int  | 数据库id |  无  |
   | invalid | bool | 是否无效 | true |
   | type | bool | 什么类型比赛，"wheel_war"，"final" | 无 |

7. 通过battle_id查询是否有当前场次比赛`GET /api/tank/battle/isexists?_tc={...}`  Request:

   |   参数    |  类型  |                 说明                 | 默认 |
   | :-------: | :---------------: | :----------------------------------: | :---------------: |
   | battle_id |  int   |              比赛场次id              |  无  |
   |   type    | string | 车轮战（wheel_war）？决赛？（final） |  无  |

8. 实时推送比赛分数统计 `ws /api/tank/ws?mode={xxx}`

   mode = wheel_war || final 

   一旦有数据库更新，会推送一个消息，消息内容为`排名改变啦啦啦啦`
