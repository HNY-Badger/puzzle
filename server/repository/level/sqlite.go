package level

import (
	"database/sql"
	"fmt"
	"godb/model"
)

type SQLRepo struct {
	DB *sql.DB
}

func (r *SQLRepo) GetLevel(levelId uint, userId uint) (model.LevelResp, error) {
	level := model.LevelResp{}
	levelRows, err := r.DB.Query(`SELECT * FROM `+"`level`"+` WHERE id = ?`, levelId)
	if err != nil {
		return level, fmt.Errorf("failed to perfom select: %w", err)
	}
	for levelRows.Next() {
		levelRows.Scan(&level.ID, &level.Name)
	}
	roundRows, err := r.DB.Query(`SELECT * FROM `+"`round`"+` WHERE level_id = ?`, levelId)
	if err != nil {
		return level, fmt.Errorf("failed to perfom select: %w", err)
	}
	for roundRows.Next() {
		round := model.RoundResp{}
		roundRows.Scan(&round.ID, &round.LevelID, &round.Name, &round.Image, &round.Author, &round.Year, &round.Language)
		completedRows, err := r.DB.Query(`SELECT * FROM `+"`completed`"+` WHERE round_id = ? AND user_id = ?`, round.ID, userId)
		if err != nil {
			return level, fmt.Errorf("failed to perfom select: %w", err)
		}
		completed := false
		for completedRows.Next() {
			completed = true
		}
		round.Completed = completed
		sentenceRows, err := r.DB.Query(`SELECT * FROM `+"`sentence`"+` WHERE round_id = ?`, round.ID)
		if err != nil {
			return level, fmt.Errorf("failed to perfom select: %w", err)
		}
		for sentenceRows.Next() {
			sentence := model.Sentence{}
			sentenceRows.Scan(&sentence.ID, &sentence.RoundID, &sentence.Text, &sentence.Translation)
			round.Sentences = append(round.Sentences, sentence)
		}
		level.Rounds = append(level.Rounds, round)
	}
	return level, nil
}

func (r *SQLRepo) CompleteRound(roundId uint, userId uint) (error) {
	check, err := r.DB.Query(`SELECT * FROM `+"`completed`"+` WHERE round_id = ? AND user_id = ?`, roundId, userId)
	if err != nil {
		return err
	}
	exist := false
	for check.Next() {
		exist = true
	}
	if exist {
		return nil
	}
	statement, err := r.DB.Prepare("INSERT INTO `completed` (round_id, user_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(roundId, userId)
	if err != nil {
		return err
	}
	return nil
}
