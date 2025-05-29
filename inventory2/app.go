package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	DB		*sql.DB
	Router	*mux.Router
}

func (app * App) Initialise() error {
	connetionString := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", DbUser, DbPassword, DbName)
	var err error
	app.DB, err = sql.Open("mysql", connetionString)
	if err != nil {
		return err
	}

	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRouters()
	return nil
}

func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router))
}

func (app *App) handleRouters() {
	app.Router.HandleFunc("/home", app.homePage)
	app.Router.HandleFunc("/videogames", app.getVideoGames).Methods("GET")
}

func sendError(w http.ResponseWriter, statusCode int, err string) {
	error_message := map[string]string{"error":err}
	sendResponse(w, statusCode, error_message)

}
func sendResponse(w http.ResponseWriter, statusCode int, payLoad interface{}) {
	response, err := json.Marshal(payLoad)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func (app *App) homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Homepage!")
	log.Println("Endpoint hit: homepage")
}

func (app *App) getVideoGames(w http.ResponseWriter, r *http.Request) {
	videoGames, err := getVideoGames(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
	}
	sendResponse(w, http.StatusOK, videoGames)
}