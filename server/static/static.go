package static

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"godb/model"
	"io"
	"os"
)

func ConvertStaticToDb(db *sql.DB) {
	// ConvertOriginalToCustom()
	err := createTables(db)
	if err != nil {
		panic(err)
	}
	err = fillTables(db)
	if err != nil {
		panic(err)
	}
}

func createTables(db *sql.DB) error {
	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS` + "`level`" + `(
		id INTEGER PRIMARY KEY,
		name TEXT)`)
	if err != nil {
		return err
	}
	statement.Exec()
	statement, err = db.Prepare(`CREATE TABLE IF NOT EXISTS` + "`round`" + `(
		id INTEGER PRIMARY KEY,
		level_id INTEGER,
		name TEXT,
		image TEXT,
		author TEXT,
		year TEXT,
		language TEXT)`)
	if err != nil {
		return err
	}
	statement.Exec()
	statement, err = db.Prepare(`CREATE TABLE IF NOT EXISTS` + "`sentence`" + `(
		id INTEGER PRIMARY KEY,
		round_id INTEGER,
		text TEXT,
		translation TEXT)`)
	if err != nil {
		return err
	}
	statement.Exec()
	statement, err = db.Prepare(`CREATE TABLE IF NOT EXISTS` + "`completed`" + `(
		id INTEGER PRIMARY KEY,
		round_id INTEGER,
		user_id INTEGER)`)
	if err != nil {
		return err
	}
	statement.Exec()
	return nil
}

func fillTables(db *sql.DB) error {
	for i := uint(1); i <= 6; i++ {
		jsonFile, err := os.Open(fmt.Sprintf("./static/data/%d.json", i))
		if err != nil {
			return err
		}
		byteValue, err := io.ReadAll(jsonFile)
		jsonFile.Close()
		if err != nil {
			return err
		}

		level := model.Level{}
		json.Unmarshal(byteValue, &level)
		statement, err := db.Prepare("INSERT OR IGNORE INTO `level` (id, name) VALUES (?, ?)")
		if err != nil {
			return err
		}
		_, err = statement.Exec(level.ID, level.Name)
		if err != nil {
			return err
		}
		for j := 0; j < len(level.Rounds); j++ {
			statement, err = db.Prepare("INSERT OR IGNORE INTO `round`" + `(
				id,
				level_id,
				name,
				image,
				author,
				year,
				language) VALUES (?, ?, ?, ?, ?, ?, ?)`)
			if err != nil {
				return err
			}
			_, err = statement.Exec(level.Rounds[j].ID,
				level.Rounds[j].LevelID,
				level.Rounds[j].Name,
				level.Rounds[j].Image,
				level.Rounds[j].Author,
				level.Rounds[j].Year,
				level.Rounds[j].Language)
			if err != nil {
				return err
			}
			for k := 0; k < len(level.Rounds[j].Sentences); k++ {
				statement, err = db.Prepare("INSERT OR IGNORE INTO `sentence`" + `(
					id,
					round_id,
					text,
					translation) VALUES (?, ?, ?, ?)`)
				if err != nil {
					return err
				}
				_, err = statement.Exec(level.Rounds[j].Sentences[k].ID,
					level.Rounds[j].Sentences[k].RoundID,
					level.Rounds[j].Sentences[k].Text,
					level.Rounds[j].Sentences[k].Translation)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
