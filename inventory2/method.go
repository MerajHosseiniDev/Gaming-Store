package main

import (
	"database/sql"
	"fmt"
)

type Videogame struct {
	Id    int  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func getVideoGames(db *sql.DB) ([]Videogame, error) {
	query := "SELECT id, name, price FROM videogames"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	videoGames := []Videogame{}
	for rows.Next() {
		var v Videogame
		err := rows.Scan(&v.Id, &v.Name, &v.Price)
		if err !=nil {
			return nil, err
		}
		videoGames = append(videoGames, v)
	}
	return videoGames, nil
}

func (v *Videogame) getVideoGame(db *sql.DB) error {
	query := fmt.Sprintf("SELECT name, price FROM videogames where id=%v", v.Id)
	rows := db.QueryRow(query)
	err := rows.Scan(&v.Name, v.Price)
	if err != nil {
		return err
	}
	return nil
}