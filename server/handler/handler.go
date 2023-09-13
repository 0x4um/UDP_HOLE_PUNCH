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
	if rowCount > 1 {
		fmt.Println("there are rows here")
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
	switch other[:start] {
	case "findpeer":
		findPeers(db, conn)
	default:
		fmt.Println("def")
	}
}