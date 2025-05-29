package main

import "database/sql"

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