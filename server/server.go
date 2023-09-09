package main 

import (
	"fmt"
	"net"
	// "strings"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"server/handler"
)

func checkTable(db *sql.DB, tableName string) (bool, error) {
	query := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s'", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return false, err 
	
	}
	defer rows.Close()

	return rows.Next(), nil
}

func main(){
	
	db, err := sql.Open("sqlite3", "db.db")
	if err != nil {
		fmt.Println("error opening")
		return
	}
	defer db.Close()

	tableName := "users"

	exists, err := checkTable(db, tableName)
	if err != nil {
		fmt.Println("error checking", err)
		return
	}

	if exists {
		fmt.Println("the table exists")

	} else {
		fmt.Println("table does not exists.. creating one")
		createTableSql := `
			CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ip TEXT NOT NULL,
			other TEXT NOT NULL
		)
		`
		
		_, err := db.Exec(createTableSql)
		if err != nil {
			fmt.Println("unable to create table", err)
			return
		}
	}

	rows, err := db.Query("SELECT id FROM users")
	if err != nil {
		fmt.Println("error select")
		return
	}
	defer rows.Close()


	ln, err := net.Listen("tcp", ":12000")
	if err != nil {
		fmt.Println("error starting")
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("accepting error")
			continue
		}


		go handleConn(conn)

	}
}

func index(){
	
	fmt.Println("calling on index")
}


func handleConn(conn net.Conn){

//this will pretty much be the main function
//the "n" value will be put into the function handler and redirect it properly 
// handleConn() -> funcHandle() -> function based off of buffer[:n] value 
	defer conn.Close()

	
	buffer := make([]byte, 1024)

	for {

		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("error reading", conn.RemoteAddr())
			return 
		}

		data := string(buffer[:n])

		fmt.Println(data)

		handler.Exec(data, n, conn)


	}
}