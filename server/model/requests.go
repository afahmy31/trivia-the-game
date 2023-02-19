package model

type UpdateScoreRequest struct {
	TeamName  string `json:"teamName"`
	Operation string `json:"operation"`
	GameId    string `json:"gameId"`
}
