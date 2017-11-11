package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"fmt"
)

var connection *sql.DB

func connect(host string, port int, user string, password string, database string) *sql.DB {
	var connectString string = user + ":" + password + "@tcp(" + host + ":" + strconv.Itoa(port) + ")/" + database

	fmt.Println("Try connect to the database: " + connectString)

	db, err := sql.Open("mysql", connectString)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Connected to the database: " + connectString)

	return db
}

func insert(clientId string, objectType string, objectId string) {
	stmtIns, err := connection.Prepare("INSERT INTO store (client_id, object_type, object_id) VALUES(?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer stmtIns.Close()
	_, err = stmtIns.Exec(clientId, objectType, objectId)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func delete(clientId string, objectType string, objectId string) {
	stmtOut, err := connection.Prepare("DELETE FROM store WHERE client_id=? AND object_type=? AND object_id=?")
	_, err = stmtOut.Exec(clientId, objectType, objectId)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func deleteClient(clientId string) {
	stmtOut, err := connection.Prepare("DELETE FROM store WHERE client_id=?")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = stmtOut.Exec(clientId)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getClients(objectType string, objectId string) []string {
	rows, err := connection.Query("SELECT client_id FROM store WHERE object_type=? AND object_id=?", objectType, objectId)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	var clientsIds []string = []string{}

	for rows.Next() {
		var clientId string
		rows.Scan(&clientId)

		clientsIds = append(clientsIds, clientId)
	}

	return clientsIds
}
