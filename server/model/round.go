package model

type Round struct {
	ID        uint       `json:"id"`
	LevelID   uint       `json:"level_id"`
	Name      string     `json:"name"`
	Image     string     `json:"image"`
	Author    string     `json:"author"`
	Year      string     `json:"year"`
	Language  string     `json:"language"`
	Sentences []Sentence `json:"sentences"`
}

type RoundResp struct {
	ID        uint       `json:"id"`
	LevelID   uint       `json:"level_id"`
	Name      string     `json:"name"`
	Image     string     `json:"image"`
	Author    string     `json:"author"`
	Year      string     `json:"year"`
	Language  string     `json:"language"`
	Completed bool       `json:"completed"`
	Sentences []Sentence `json:"sentences"`
}
