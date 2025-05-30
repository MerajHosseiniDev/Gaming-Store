package main

import (
	"database/sql"
	"errors"
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
	query := "SELECT name, price FROM videogames where id=?"
	rows := db.QueryRow(query, v.Id)
	err := rows.Scan(&v.Name, &v.Price)
	if err != nil {
		return err
	}
	return nil
}

func (v *Videogame) createVideoGame(db *sql.DB) error {
	query := "insert into videogames(name,price) values(?, ?)"
	result, err := db.Exec(query, v.Name, v.Price)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil{
		return err
	}
	v.Id = int(id)
	return nil
}

func (v *Videogame) updateVideoGame(db *sql.DB) error {
	query := "UPDATE videogames SET name = ?, price = ? WHERE id = ?"
	result, err := db.Exec(query, v.Name, v.Price, v.Id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected ==0 {
		return errors.New("not such row exists")
	}
	return err
}