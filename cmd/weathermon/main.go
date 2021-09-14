package main

import (
	"encoding/json"
	"fmt"
	"github.com/Lytreksq/weathermon/storage"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Response struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

var sqliteDatabase *storage.Storage

func main() {
	var err error
	sqliteDatabase, err = storage.NewStorage("/home/dima/mybase.db")
	if err != nil {
		log.Fatalln(err)
	}
	defer sqliteDatabase.Shutdown()

	http.HandleFunc("/", hello)
	go func() {
		http.ListenAndServe(":8090", nil)
	}()
	key := ""

	for i := 0; i < 100; i++ {
		temp, err := getTemp(key)
		if err != nil {
			log.Fatalln(err)
		}
		err = sqliteDatabase.InsertWeather(temp)
		if err != nil {
			log.Fatalln(err)
		}
		time.Sleep(time.Second * 2)
	}

}
func hello(w http.ResponseWriter, req *http.Request) {
	temperatures, err := sqliteDatabase.DisplayWeather()
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	fmt.Fprintln(w, temperatures)

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
