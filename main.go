package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type VideoGame struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Publisher string `json:"publisher"`
	Genre     string `json:"genre"`
}

func main() {
	db, err := initTCPConnectionPool()
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, GCloud WebApp!  123<3")
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.GET("/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Message string }{Message: "Hello, World!"})
	})

	e.GET("/addvg",func(c echo.Context) error {
		err := db.QueryRow("INSERT INTO videogames(name,publisher,genre) VALUES ($1,$2,$3)","GTA","ROCKSTAR","SHOOTER").Err()
		if err != nil {
			return c.JSON(http.StatusInternalServerError,struct{Error string}{Error: err.Error()})
		}
		return c.JSON(http.StatusOK, struct{ Message string }{Message: "VideoGame added!"})
	})

	e.GET("/vg",func(c echo.Context) error {
		rows, err := db.Query("SELECT * FROM videogames")
		if err != nil {
			return c.JSON(http.StatusBadRequest,struct{Error string}{Error: err.Error()})
		}
		defer rows.Close()

		var videogames []VideoGame

		// Loop through rows, using Scan to assign column data to struct fields.
		for rows.Next() {
			var vg VideoGame
			if err := rows.Scan(&vg.ID, &vg.Name, &vg.Publisher,
				&vg.Genre); err != nil {
				return c.JSON(http.StatusInternalServerError,struct{Error string}{Error: err.Error()})
			}
			videogames = append(videogames, vg)
		}
		return c.JSON(http.StatusOK, videogames)
	})

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Warning: %s environment variable not set.\n", k)
	}
	return v
}

func initTCPConnectionPool() (*sql.DB, error) {
	var (
		dbUser = mustGetenv("DB_USER")
		dbPass = mustGetenv("DB_PASS")
		dbName = mustGetenv("DB_NAME")
	)

	dbURI := fmt.Sprintf("user=%s password=%s database=%s host=127.0.0.1 port=5432", dbUser, dbPass, dbName)
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection to db: %w", err)
	}
	dbPool.SetMaxIdleConns(5)
	dbPool.SetMaxOpenConns(7)
	dbPool.SetConnMaxLifetime(1800)
	return dbPool, nil
}
