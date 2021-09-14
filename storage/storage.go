package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(filepath string) (storage *Storage, err error) {
	_, err = os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(filepath)
			if err != nil {
				return nil, err
			}
			file.Close()
		} else {
			return nil, err
		}
	}
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}
	var result Storage
	result.db = &db
	err = result.createTable()
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *Storage) Shutdown() {
	s.db.Close()
}

func (s *Storage) createTable() error {
	createWeatherTableSQL := `CREATE TABLE IF NOT EXISTS weather (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,    
    "temperature" FLOAT);`

	statement, err := s.db.Prepare(createWeatherTableSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) InsertWeather(temperature float64) error {
	insertWeatherSQL := `INSERT INTO weather(temperature) VALUES (?)`
	statement, err := s.db.Prepare(insertWeatherSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(temperature)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DisplayWeather() ([]float64, error) {
	row, err := s.db.Query("SELECT * FROM weather")
	if err != nil {
		return nil, err
	}
	var result []float64
	defer row.Close()
	for row.Next() {
		var id int
		var temperature float64
		row.Scan(&id, &temperature)
		result = append(result, temperature)
	}
	return result, nil

}
