package main

import (
	"net/http"

	"github.com/takasp/go-slackbot-cloud-functions-sample"
)

func main() {
	http.HandleFunc("/hello", function.HelloCommand)
	http.HandleFunc("/weather", function.WeatherCommand)
	http.ListenAndServe(":8080", nil)
}
