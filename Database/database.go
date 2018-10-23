package Database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type MySqlDB struct {
	conn *sql.DB
}

type ConfigSql struct {
	Username, Password, Host, DatabaseName string
	Port                                   int
}

func NewSQLDatabase(c ConfigSql) (*MySqlDB, error) {
	cred := c.Username + ":" + c.Password + "@"

	dataStoreName := fmt.Sprintf("%stcp([%s]:%d)/%s", cred, c.Host, c.Port, c.DatabaseName)

	conn, err := sql.Open("mysql", dataStoreName)
	if err != nil {
		return nil, fmt.Errorf("could not get a connection: %v", err)
	}
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("could not establish connection: %v", err)
	}

	return &MySqlDB{
		conn: conn,
	}, nil
}

func (db *MySqlDB) AddSong(artist, title string) error {
	query := `INSERT INTO songs (song, url) VALUES (?, ?)`

	r, err := db.conn.Exec(query, artist, title)
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

func (db *MySqlDB) Close() {
	db.conn.Close()
}
