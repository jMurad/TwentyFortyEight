package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var PoolConn *pgxpool.Pool

type CoreDB struct {
	Conn *pgxpool.Pool
}

func (db *CoreDB) DBinit() {
	var err error

	urlDb := "postgres://tester:353694@localhost:5432/example"
	db.Conn, err = pgxpool.New(context.Background(), urlDb)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to Database: %v\n", err)
		os.Exit(1)
	}
}

func (db *CoreDB) GetPlayer(plr *Players) error {
	conn, err := db.Conn.Acquire(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetPlayer.Error acquiring connection: %v\n", err)
		return err
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		`SELECT score, level, sizebox
        FROM players
        WHERE name=$1 AND sizebox=$2`,
		plr.name, plr.sizebox,
	)

	err = row.Scan(&plr.score, &plr.level, &plr.sizebox)
	if err != nil {
		// fmt.Fprintf(os.Stderr, "GetPlayer.Error row.Scan: %v\n", err)
		return err
	}
	return nil
}

func (db *CoreDB) GetBestPlayers(sizebox string) ([]Players, error) {
	conn, err := db.Conn.Acquire(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetBestPlayers.Error acquiring connection: %v\n", err)
		return []Players{}, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(),
		`SELECT *
        FROM players
		WHERE sizebox=$1
        ORDER BY sizebox DESC, score DESC
		LIMIT 5`,
		sizebox,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetBestPlayers.QueryRow failed: %v\n", err)
		return []Players{}, err
	}
	defer rows.Close()

	var prows []Players
	for rows.Next() {
		row := Players{}
		err := rows.Scan(&row.name, &row.score, &row.level, &row.sizebox)
		if err != nil {
			return []Players{}, err
		}
		prows = append(prows, row)
	}
	return prows, nil
}

func (db *CoreDB) AddPlayer(plr *Players) error {
	conn, err := db.Conn.Acquire(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "AddPlayer.Error acquiring connection: %v\n", err)
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(),
		`INSERT INTO players (name, score, level, sizebox) 
		VALUES ($1, $2, $3, $4)`,
		plr.name,
		plr.score,
		plr.level,
		plr.sizebox,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "AddPlayer.QueryRow failed: %v\n", err)
		return err
	}
	return nil
}

func (db *CoreDB) UpdatePlayer(plr *Players) error {
	conn, err := db.Conn.Acquire(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "UpdatePlayer.Error acquiring connection: %v\n", err)
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(),
		`UPDATE players 
		SET score = $1, 
			level = $2, 
			sizebox = $3 
		WHERE name = $4`,
		plr.score,
		plr.level,
		plr.sizebox,
		plr.name,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "UpdatePlayer.QueryRow failed: %v\n", err)
		return err
	}
	return nil
}
