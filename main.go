package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

// App struct
type App struct {
	URLTemplate string
	Accept      string
	Params      *Params
	Client      http.Client
}

// Params struct
type Params struct {
	Query  string
	Source string
	Size   string
	Offset string
}

// Body struct
type Body struct {
	Result Result `json:"result"`
}

// Result struct
type Result struct {
	Documents []Document `json:"document"`
}

// Document struct
type Document struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

var app *App

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	params := InitParams()
	app = NewApp(params)

	searchRes, searchErr := app.Search()
	if searchErr != nil {
		panic(searchErr)
	}
	for _, d := range searchRes.Result.Documents {
		log.Println(d.Title)
	}
}

// NewApp constructor
func NewApp(params *Params) *App {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &App{
		URLTemplate: os.Getenv("URL_TEMPLATE"),
		Accept:      os.Getenv("ACCEPT_HEADER"),
		Params:      params,
		Client:      http.Client{},
	}
}

// InitParams initializes params
func InitParams() *Params {
	return &Params{
		Query:  url.QueryEscape(os.Getenv("Q_PARAM")),
		Source: url.QueryEscape(os.Getenv("SRC_PARAM")),
		Size:   os.Getenv("SIZE_PARAM"),
		Offset: os.Getenv("OFFSET_PARAM"),
	}
}

// Search method
func (app *App) Search() (*Body, error) {
	url := fmt.Sprintf(
		app.URLTemplate,
		app.Params.Query,
		app.Params.Source,
		app.Params.Size,
		app.Params.Offset)

	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	req.Header.Set("Accept", app.Accept)

	response, responseErr := app.Client.Do(req)
	if responseErr != nil {
		return nil, responseErr
	}

	var body *Body
	decodeErr := json.NewDecoder(response.Body).Decode(&body)
	if decodeErr != nil {
		return nil, decodeErr
	}
	return body, nil
}
