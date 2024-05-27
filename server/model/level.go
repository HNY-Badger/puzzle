package model

type Level struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Rounds []Round `json:"rounds"`
}

type Request struct {
	Cookie string `json:"cookie"`
	Id     uint   `json:"id"`
}

type LevelResp struct {
	ID     uint        `json:"id"`
	Name   string      `json:"name"`
	Rounds []RoundResp `json:"rounds"`
}
