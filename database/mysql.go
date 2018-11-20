package Database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"cloud.google.com/go/firestore"
)

type Repo struct {
	UserInfo *sql.DB
	Session  *firestore.Client
}

type ConfigSql struct {
	Username, Password, Host, DatabaseName string
	Port                                   int
}

func NewRepository(sqlConfig ConfigSql, fireStoreConfig ConfigFireStore) (*Repo, error) {
	sql, err := newSQLDatabase(sqlConfig)
	if err != nil {
		return nil, fmt.Errorf("SQL: %v", err)
	}
	fireStore, err := newSessionStore(fireStoreConfig)
	if err != nil {
		return nil, fmt.Errorf("firestore: %v", err)
	}

	return &Repo {
		UserInfo: sql,
		Session:  fireStore,
	}, nil
}

func newSQLDatabase(config ConfigSql) (*sql.DB, error) {
	cred := config.Username + ":" + config.Password + "@"

	dataStoreName := fmt.Sprintf("%stcp([%s]:%d)/%s", cred, config.Host, config.Port, config.DatabaseName)

	conn, err := sql.Open("mysql", dataStoreName)
	if err != nil {
		return nil, fmt.Errorf("could not get a connection: %v", err)
	}
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("could not establish connection: %v", err)
	}

	return conn, nil
}

func (db *Repo) AddSong(artist, title string) error {
	query := `INSERT INTO songs (song, url) VALUES (?, ?)`

	r, err := db.UserInfo.Exec(query, artist, title)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}
	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected: %v", err)
	} else if rowsAffected != 1 {
		return fmt.Errorf("expected 1 row affected, got %d", rowsAffected)
	}
	return nil

}

func (db *Repo) CloseSQL() {
	db.UserInfo.Close()
}
