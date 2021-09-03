package main

import (
	"encoding/json"
	"fmt"
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
type Templist struct {
	list []float64
}
	var tl Templist
func main() {
	http.HandleFunc("/", hello)
	go func() {
		http.ListenAndServe(":8090", nil)
	}()
	key :=
	temp, err := getTemp(key)
	for i := 0; i < 100; i++ {
		if err != nil {
			log.Fatalln(err)
		}
		tl.list = append(tl.list, temp)
		fmt.Println(temp, tl)
		time.Sleep(time.Second * 2)
	}

}
func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintln(w, tl.list)
}
func getTemp(key string)(float64, error) {
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







