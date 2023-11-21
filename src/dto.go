package src

type Player struct {
	ID           string   `json:"id"`
	PlayerNumber int      `json:"playerNumber"`
	Username     string   `json:"username"`
	Collected    []string `json:"collected"`
}
