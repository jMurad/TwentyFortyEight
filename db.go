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

func (db *CoreDB) GetPlayer(name string) (pl Players, err error) {
	conn, err := db.Conn.Acquire(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error acquiring connection: %v\n", err)
		return Players{}, err
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		`SELECT score, level, sizebox
        FROM players
        WHERE name==$1`,
		name,
	)
	err = row.Scan(&pl.score, &pl.level, &pl.sizebox)
	if err != nil {
		return Players{}, err
	}
	return
}

func (db *CoreDB) AddPlayer(pl *Players) error {
	conn, err := db.Conn.Acquire(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error acquiring connection: %v\n", err)
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(),
		`INSERT INTO players (name, score, level, sizebox) 
		VALUES ($1, $2, $3, $4::timestamptz)`,
		pl.name, pl.score, pl.level, pl.sizebox,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return err
	}
	return nil
}
