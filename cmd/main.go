package main

import (
	"MyMusic/Database"
	"MyMusic/server"
	"fmt"
	"net/http"
	"os"
)

func main() {
	const port = "8080"
	db, err := Database.NewSQLDatabase(Database.ConfigSql{
		Username: os.Getenv("USER_SQLCLOUD"),
		Password: os.Getenv("PASS_SQLCLOUD"),
		DatabaseName: "mymusic",
		Host: "localhost",
		Port: 3306,
	})
	if err != nil {
		fmt.Printf("mysql: %s", err.Error())
	}

	s := server.NewServer(db)
	http.ListenAndServe(":"+port, s.Mux)
}
