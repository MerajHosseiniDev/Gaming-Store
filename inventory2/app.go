package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	DB		*sql.DB
	Router	*mux.Router
}

func (app * App) Initialise() error {
	connectionString := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", DbUser, DbPassword, DbName)
	var err error
	app.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}

	err = app.DB.Ping()
	if err != nil {
		app.DB.Close()
		return fmt.Errorf("failed to ping database: %v", err)
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
	app.Router.HandleFunc("/videogames/{id}", app.getVideoGame).Methods("GET")
	app.Router.HandleFunc("/videogames", app.createVideoGame).Methods("POST")
	app.Router.HandleFunc("/videogames/{id}", app.updateVideoGame).Methods("PUT")
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

func (app *App) getVideoGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, err := strconv.Atoi(vars["id"])
	if err !=nil {
		sendError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	v := Videogame{Id: key}
	err = v.getVideoGame(app.DB)
	if err != nil {
		switch err {
		case sql. ErrNoRows:
			sendError(w, http.StatusNotFound, "product not found")
		default:
			sendError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	sendResponse(w, http.StatusOK, v)
}

func (app *App) createVideoGame(w http.ResponseWriter, r *http.Request) {
	var v Videogame
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	err = v.createVideoGame(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusCreated, v)
}

func (app *App) updateVideoGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, err := strconv.Atoi(vars["id"])
	if err != nil{
		sendError(w, http.StatusBadRequest, "invalid product id")
		return
	}
	var v Videogame
	err = json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid payload Request")
		return
	}
	v.Id = key
	err = v.updateVideoGame(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
	}
	sendResponse(w, http.StatusOK, v)
}
