package db

// 力扣大赛model

type LeetCode struct {
	Rank  uint16 `json:"rank"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Score uint16 `json:"score"`
	Bout  string `json:"bout"`
	Bonus uint16 `json:"bonus"`
}

func init() {
	_ = DB.AutoMigrate(&LeetCode{})
}
