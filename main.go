package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
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

	app = NewApp()
	http.ListenAndServe(":3000", Routes())
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

// NewApp constructor
func NewApp() *App {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	params := InitParams()
	return &App{
		URLTemplate: os.Getenv("URL_TEMPLATE"),
		Accept:      os.Getenv("ACCEPT_HEADER"),
		Params:      params,
		Client:      http.Client{},
	}
}

// Routes func
func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)
	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api", SearchRoutes())
	})
	return router
}

// SearchRoutes func
func SearchRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", Search)
	return router
}

// Search method
func Search(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf(
		app.URLTemplate,
		app.Params.Query,
		app.Params.Source,
		app.Params.Size,
		app.Params.Offset)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", app.Accept)

	response, _ := app.Client.Do(req)

	var body *Body
	json.NewDecoder(response.Body).Decode(&body)
	render.JSON(w, r, body)
}
