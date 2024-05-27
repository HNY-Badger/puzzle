package model

type Sentence struct {
	ID          uint   `json:"id"`
	RoundID     uint   `json:"round_id"`
	Text        string `json:"text"`
	Translation string `json:"translation"`
}
