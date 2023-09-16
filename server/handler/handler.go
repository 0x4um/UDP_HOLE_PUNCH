package handler

import (
	"fmt"
	"strings"
	"net"
	"database/sql"
	"encoding/json"
	"io/ioutil"
)

func dialCheck(conn net.Conn){
	fmt.Println("checking for dial")
}

func function1(conn net.Conn) {
	conn.Write([]byte("return hello"))
	fmt.Println("function1 hello from function1")
}

func newPeer(db *sql.DB, conn net.Conn, id string){
	fmt.Println("new peer func")
	fmt.Println("id id ", id)
	altID := strings.TrimSpace(id)
	getPeerTableSQL := `SELECT peer_uuid FROM peertable WHERE peer_uuid = ?`
	stmt, err := db.Prepare(getPeerTableSQL)
	if err != nil {
		fmt.Println("error getting data from peertable", err)
		return 
	}
	var done string
	err = stmt.QueryRow(altID).Scan(&done)
	if err != nil {
		fmt.Println("error scanning value", err)
		conn.Write([]byte("fuck you"))
		return
	}
	fmt.Println(done)
}

func findPeers(db *sql.DB, conn net.Conn){
	
	checkForLocalPeersSQL := `SELECT peer_uuid, peer_ip, peer_public_key FROM peertable`
	rows, err := db.Query(checkForLocalPeersSQL)
	if err != nil {
		fmt.Println("error checking for sql", err)
		return
	}
	defer rows.Close()
	rowCount := 0

	

	for rows.Next() {
		rowCount++
		var peer_uuid, peer_ip, peer_public_key string
		if err := rows.Scan(&peer_uuid, &peer_ip, &peer_public_key); err != nil {
			fmt.Println(err)
			return

		}
		fmt.Println(peer_uuid)
		fmt.Println(peer_ip)
		fmt.Println(peer_public_key)
		fmt.Println(rowCount)
	}
	if rowCount >= 1 {
		fmt.Println("there are rows here")
		fmt.Println("writing to peer nodes")
		getPublicUUIDSQL := `SELECT public_uuid FROM general`
		getPublicUUID, err := db.Query(getPublicUUIDSQL)
		if err != nil {
			fmt.Println("error getting uuid", err)
			return
		}
		defer getPublicUUID.Close()
		var peer_uuid string
		for getPublicUUID.Next() {
			//var public_uuid string
			if err := getPublicUUID.Scan(&peer_uuid); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(peer_uuid)

			
		}
		fmt.Println("writing using", peer_uuid)
		conn.Write([]byte("findpeer:" + peer_uuid))
		
	} else {
		fmt.Println("check the boot strap node")
		CONFIG, err := ioutil.ReadFile("./../config/config.json")
		
		if err != nil {
			fmt.Println("error reading file", err)
			return
		
		}
		var configData map[string]interface{}

		if err := json.Unmarshal(CONFIG, &configData); err != nil {
			fmt.Println("error herer", err)
			return
		}

		bootNode := configData["bootstrap"].(string)
		fmt.Println(bootNode)
		if bootNode == "" {
			fmt.Println("boot node is null")
		}
	}
}

func Exec(data string, n int, conn net.Conn, db *sql.DB) {

	fmt.Println(data)
	fmt.Println(n)
	// functions := map[string]func(){
	// 	"function1": function1,
	// 	"function2": function2,
	// }
	// buffer := make([]byte, 1024)



	other := strings.TrimSpace(data)
	var start int
	for i, char := range other {
		if string(char) == ":" {
			fmt.Println(i)
			start = i
		}
	}

	fmt.Println(start, "found here")
	fmt.Println(other, " other other")
	fmt.Println(other[:start], "start parse")
	fmt.Println(other[start:], "id here")
	switch other[:start] {
	case "findpeer":
		findPeers(db, conn)
	case "newpeer":
		newPeer(db, conn, other[(start + 1):])
	default:
		fmt.Println("def")
	}
}