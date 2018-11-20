package main

import (
	"GoMusic/database"
	"GoMusic/server"
	"fmt"
	"net/http"
	"os"
	"context"
)

func main() {
	const port = "8080"

	configSql := Database.ConfigSql{
		Username: os.Getenv("USER_SQLCLOUD"),
		Password: os.Getenv("PASS_SQLCLOUD"),
		DatabaseName: "mymusic",
		Host: "localhost",
		Port: 3306,
	}

	configFireStore := 	Database.ConfigFireStore {
		ProjectID: "mymusic-220213",
	}

	repo, err := Database.NewRepository(configSql , configFireStore)
	if err != nil {
		fmt.Printf("Repo %s", err.Error())
	}
	defer repo.CloseSQL()
	//
	_ = repo.AddSession(context.Background(), "aaa", "kecske@kicsi.com")
	fmt.Println("exists: ", repo.IsValidSession(context.Background(), "aaa"))

	asd, err := repo.GetSessionEmail(context.Background(), "aaa")
	if err != nil {
		fmt.Println("get: ", err)
	}
	fmt.Println("session email: ", asd)

	s := server.NewServer(repo)
	http.ListenAndServe(":"+port, s.Mux)
}
