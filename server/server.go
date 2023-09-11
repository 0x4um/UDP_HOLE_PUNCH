package main 

import (
	"fmt"
	"net"
	// "strings"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"server/handler"
	"github.com/google/uuid"
	"os"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
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

func checkForPublicSQL(db *sql.DB) (int, error) {
	fmt.Println("the table exists... checking for public data")
	checkForPublicSQL := `select * from general`

	rows, err := db.Query(checkForPublicSQL)
	if err != nil {
		fmt.Println("error checking for sql", err)
	}
	defer rows.Close()

	rowCount := 0

	for rows.Next() {
		rowCount++
	}

	if err := rows.Err(); err != nil {
		fmt.Println("error iterating", err)
	}
	return rowCount, err
}

func main(){
	
	db, err := sql.Open("sqlite3", "db.db")
	if err != nil {
		fmt.Println("error opening")
		return
	}
	defer db.Close()

	tableName := "general"

	exists, err := checkTable(db, tableName)
	if err != nil {
		fmt.Println("error checking", err)
		return
	}

	if exists {

		rowCount, err := checkForPublicSQL(db)
		if err != nil {
			fmt.Println("error checking for public sql", err)
			return
		}
		fmt.Println(rowCount, "here")
		if rowCount == 1{
			fmt.Println("there is a public profile here")
		} else {
			fmt.Println("there is not a public profile here")
			fmt.Println("creating one now..")
			privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				fmt.Println("error making priv key", err)
			}
			privateKeyPEM := &pem.Block{
				Type: "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
			}

			privateKeyFile, err := os.Create("private_key.pem")
			if err != nil {
				fmt.Println("error creating private key", err)
				return 
			}
			defer privateKeyFile.Close()
			err = pem.Encode(privateKeyFile, privateKeyPEM)
			if err != nil {
				fmt.Println("error encoding private key", err)
				return
			} 	

			fmt.Println("saved")

			publicKey := &privateKey.PublicKey
			publicKeyPEM, err := x509.MarshalPKIXPublicKey(publicKey)
			if err != nil {
				fmt.Println("error encoding public key", err)
				return
			}
			publicKeyFile, err := os.Create("public_key.pem")
			if err != nil {
				fmt.Println("error createing public key", err)
				return
			}
			defer publicKeyFile.Close()
			err = pem.Encode(publicKeyFile, &pem.Block{
				Type: "RSA PUBLIC KEY",
				Bytes: publicKeyPEM,
			})

			if err != nil {
				fmt.Println("error encoding pub key", err)
				return
			}
			fmt.Println("done")

			file, err := os.Open("public_key.pem")
			if err != nil {
				fmt.Println("error opening")
				return
			}
			defer file.Close()

			fileContents, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Println("error", err)
				return
			}

			fmt.Println("file:")
			fmt.Println(string(fileContents))

			//need a way to return the value of public ip to all incoming requests and then ask for 
			// your own server ip
			//outside user (ip included)=> server => outside user
		}
		
		


	} else {
		fmt.Println("table does not exists.. creating one")
		createTableSql := `
			CREATE TABLE IF NOT EXISTS general (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			public_ip TEXT NOT NULL,
			public_key TEXT NOT NULL,
		)
		`
		
		_, err := db.Exec(createTableSql)
		if err != nil {
			fmt.Println("unable to create table", err)
			return
		}
		os.Exit(0)
		
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
		newUuid, err := uuid.NewUUID()
		if err != nil {
			fmt.Println("error making uuid")
			return
		}

		fmt.Println(newUuid)

		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("error reading", conn.RemoteAddr())
			return 
		}

		data := string(buffer[:n])
		handler.Exec(data, n, conn)


	}
}