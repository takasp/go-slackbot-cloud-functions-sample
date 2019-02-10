package function

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/nlopes/slack"
)

var verificationToken string
var appId string

func init() {
	verificationToken = os.Getenv("VERIFICATION_TOKEN")
	appId = os.Getenv("APP_ID")
}

func HelloCommand(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !s.ValidateToken(verificationToken) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/hello":
		params := &slack.Msg{ResponseType: "in_channel", Text: "こんにちは、<@" + s.UserID + ">さん"}
		b, err := json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func WeatherCommand(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !s.ValidateToken(verificationToken) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/weather":
		text := strings.Fields(s.Text)
		if len(text) == 0 {
			params := &slack.Msg{Text: "引数に都市名を入力してください。"}
			b, err := json.Marshal(params)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
			return
		}

		var data weatherData
		var city string

		switch text[0] {
		case "tokyo":
			city = "東京"
			data, err = query("1850147")
		case "osaka":
			city = "大阪"
			data, err = query("1853908")
		case "nagoya":
			city = "名古屋"
			data, err = query("1856057")
		default:
			params := &slack.Msg{Text: text[0] + "には対応していません。"}
			b, err := json.Marshal(params)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := weatherSlackMessageBuilder(city, data)
		b, err := json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func weatherSlackMessageBuilder(city string, data weatherData) *slack.Msg {
	text := fmt.Sprintf("%s(%s)", data.Weather[0].Main, data.Weather[0].Description)
	temperature := fmt.Sprintf("%.0f ℃", math.Floor(data.Main.Temp+.5))
	windSpeed := fmt.Sprintf("%.1f m/s", math.Floor(data.Wind.Speed+.5))

	attachment := slack.Attachment{
		Fallback:  "Current Weather",
		Title:     "現在の" + city + "のお天気",
		TitleLink: "https://openweathermap.org/city/" + strconv.Itoa(data.Id),
		Text:      text,
		Fields: []slack.AttachmentField{
			{
				Title: "気温",
				Value: temperature,
				Short: true,
			},
			{
				Title: "風速",
				Value: windSpeed,
				Short: true,
			},
		},
		ThumbURL:   "https://openweathermap.org/img/w/" + data.Weather[0].Icon + ".png",
		Footer:     "OpenWeatherMap",
		FooterIcon: "https://openweathermap.org/themes/openweathermap/assets/vendor/owm/img/icons/logo_60x60.png",
		Ts:         data.Dt,
	}

	var attachments []slack.Attachment

	params := &slack.Msg{ResponseType: "in_channel", Attachments: append(attachments, attachment)}
	return params
}

func query(cityId string) (weatherData, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?units=metric&APPID=" + appId + "&id=" + cityId)
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var d weatherData

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}

	return d, nil
}

type weatherData struct {
	Id      int         `json:"id"`
	Name    string      `json:"name"`
	Weather []Weather   `json:"weather"`
	Main    Main        `json:"main"`
	Wind    Wind        `json:"wind"`
	Dt      json.Number `json:"dt"`
}

type Weather struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Main struct {
	Temp float64 `json:"temp"`
}

type Wind struct {
	Speed float64 `json:"speed"`
}
