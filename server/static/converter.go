package static

import (
	"encoding/json"
	"fmt"
	"godb/model"
	"io"
	"os"
	"strings"
)

func ConvertOriginalToCustom() error {
	type Sentence struct {
		Text        string `json:"textExample"`
		Translation string `json:"textExampleTranslate"`
	}
	type RoundData struct {
		Name     string `json:"name"`
		Image    string `json:"imageSrc"`
		Author   string `json:"author"`
		Year     string `json:"year"`
		Language string `json:"language"`
	}
	type Round struct {
		Data      RoundData  `json:"levelData"`
		Sentences []Sentence `json:"words"`
	}
	type Level struct {
		Rounds []Round `json:"rounds"`
	}
	roundID := 1
	sentenceID := 1
	for i := uint(1); i <= 6; i++ {
		jsonFile, err := os.Open(fmt.Sprintf("./static/data/original%d.json", i))
		if err != nil {
			return err
		}
		byteValue, err := io.ReadAll(jsonFile)
		jsonFile.Close()
		if err != nil {
			return err
		}

		level := Level{}
		json.Unmarshal(byteValue, &level)

		var rounds []model.Round
		for k := 0; k < len(level.Rounds); k++ {
			current := level.Rounds[k]
			var sentences []model.Sentence
			for j := 0; j < len(current.Sentences); j++ {
				sentence := model.Sentence{
					ID:          uint(sentenceID),
					RoundID:     uint(roundID),
					Text:        current.Sentences[j].Text,
					Translation: current.Sentences[j].Translation,
				}
				sentences = append(sentences, sentence)
				sentenceID++
			}
			round := model.Round{
				ID:        uint(roundID),
				LevelID:   i,
				Name:      current.Data.Name,
				Image:     fmt.Sprintf("assets/%d/%s", i, strings.Split(current.Data.Image, "/")[1]),
				Author:    current.Data.Author,
				Year:      current.Data.Year,
				Language:  "en-US",
				Sentences: sentences,
			}
			rounds = append(rounds, round)
			roundID++
		}
		newLevel := model.Level{
			ID:     i,
			Name:   fmt.Sprintf("Уровень %d", i),
			Rounds: rounds,
		}
		levelByte, err := json.Marshal(newLevel)
		if err != nil {
			panic(err)
		}
		os.WriteFile(fmt.Sprintf("./static/data/%d.json", i), levelByte, 0644)
	}

	return nil
}
