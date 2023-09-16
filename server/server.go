package main 

import (
	"os/user"
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
	"log"
	"encoding/json"
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
	user, err := user.Current()
	if err != nil {
		fmt.Println("error with finding user id")
		return
	}
	if user.Uid != "0" {
		fmt.Println("this script requires root privs")
		os.Exit(0)
		return
	}
	db, err := sql.Open("sqlite3", "db.db")
	if err != nil {
		fmt.Println("error opening")
		return
	}
	defer db.Close()

	tableName := "general"

	peertableExist, err := checkTable(db, "peertable")
	if err != nil {
		fmt.Println("error checking", err)
		return
	}

	if peertableExist{
		fmt.Println("peertable does exist")

		
	} else {
		fmt.Println("peertable does not exists")
		fmt.Println("creating one now...")
		createPeerTableSQL := `
		CREATE TABLE IF NOT EXISTS peertable (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			peer_uuid TEXT NOT NULL,
			peer_ip TEXT NOT NULL,
			peer_public_key TEXT NOT NULL
		);
		`
		_, err := db.Exec(createPeerTableSQL)
		if err != nil {
			fmt.Println("peertable err")
			return
		}
		main()
	}


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
			newUuid, err := uuid.NewUUID()
			if err != nil {
				fmt.Println("error making uuid")
				return
			}
		
			fmt.Println(newUuid, " new uuid here")	

			fmt.Println("file:")
			fmt.Println(string(fileContents))

			//need a way to return the value of public ip to all incoming requests and then ask for 
			// your own server ip
			//outside user (ip included)=> server => outside user

			stmt, err := db.Prepare("INSERT INTO general VALUES(NULL, ?, ?)")
			if err != nil {
				fmt.Println("error inserting data into genral table")
				return
			}
			defer stmt.Close()

			_, err = stmt.Exec(newUuid, string(fileContents))
			if err != nil {
				fmt.Println("error inserting")
				return
			}
			fmt.Println("data is inserted")

		}
		
		


	} else {
		fmt.Println("table does not exists.. creating one")
		createTableSql := `
			CREATE TABLE IF NOT EXISTS general (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			public_uuid TEXT NOT NULL,
			public_key TEXT NOT NULL
		);
		`
		
		_, err := db.Exec(createTableSql)
		if err != nil {
			fmt.Println("unable to create table", err)
			return
		}
		main()
		
	}

	// rows, err := db.Query("SELECT id FROM users")
	// if err != nil {
	// 	fmt.Println("error select")
	// 	return
	// }
	// defer rows.Close()

	CONFIG, err := ioutil.ReadFile("./../config/config.json")
	if err != nil {
		fmt.Println("error reading file")
		return
	}


	var configData map[string]interface{}

	if err := json.Unmarshal(CONFIG, &configData); err != nil {
		log.Fatal("failed to unmarshal data", err)
		return
	}
	listenPort := configData["port"].(string)
	

	ln, err := net.Listen("tcp", ":" + listenPort)
	if err != nil {
		fmt.Println("error starting")
		return
	}
	fmt.Println("listening on", ":" + listenPort)
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("accepting error")
			continue
		}


		go handleConn(conn, db)

	}
}

func index(){
	
	fmt.Println("calling on index")
}


func handleConn(conn net.Conn, db *sql.DB){
//this will pretty much be the main function
//the "n" value will be put into the function handler and redirect it properly 
// handleConn() -> funcHandle() -> function based off of buffer[:n] value 
	defer conn.Close()
	// newUuid, err := uuid.NewUUID()
	// if err != nil {
	// 	fmt.Println("error making uuid")
	// 	return
	// }

	// fmt.Println(newUuid, " new uuid hrer")
	
	buffer := make([]byte, 1024)

	for {
		

		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("error reading", conn.RemoteAddr())
			return 
		}

		data := string(buffer[:n])
		
		handler.Exec(data, n, conn, db)


	}
}
