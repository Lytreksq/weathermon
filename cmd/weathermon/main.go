package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Response struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func createTable(db *sql.DB) {
	createWeatherTableSQL := `CREATE TABLE Weather (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,    
    "temperature" FLOAT);`

	log.Println("Create weather table...")
	statement, err := db.Prepare(createWeatherTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("weather table created")
}

func insertWeather(db *sql.DB, temperature float64) {
	log.Println("Inserting weather record ...")
	insertWeatherSQL := `INSERT INTO weather(temperature) VALUES (?)`
	statement, err := db.Prepare(insertWeatherSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(temperature)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func displayWeather(db *sql.DB) []float64 {
	row, err := db.Query("SELECT * FROM weather")
	if err != nil {
		log.Fatal(err)
	}
	var result []float64
	defer row.Close()
	for row.Next() {
		var id int
		var temperature float64
		row.Scan(&id, &temperature)
		result = append(result, temperature)
	}
	return result
}

var sqliteDatabase *sql.DB

func main() {
	os.Remove("sqlite-database.db")
	log.Println("Creating sql-database.db...")
	file, err := os.Create("sqlite-database.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sql-database.db created")

	sqliteDatabase, _ = sql.Open("sqlite3", "./sqlite-database.db")
	defer sqliteDatabase.Close()
	createTable(sqliteDatabase)

	http.HandleFunc("/", hello)
	go func() {
		http.ListenAndServe(":8090", nil)
	}()
	key :=

	for i := 0; i < 100; i++ {
		temp, err := getTemp(key)
		if err != nil {
			log.Fatalln(err)
		}
		insertWeather(sqliteDatabase, temp)
		time.Sleep(time.Second * 2)
	}

}
func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, displayWeather(sqliteDatabase))
}

func getTemp(key string) (float64, error) {
	url := "http://api.openweathermap.org/data/2.5/weather?q=Moscow&units=metric&appid="
	url = url + key
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	var res Response

	if err := json.Unmarshal(body, &res); err != nil {
		return 0, err
	}

	return res.Main.Temp, nil
}