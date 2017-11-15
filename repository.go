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

func insert(queueName string, clientId string, objectType string, objectId string) {
	stmtIns, err := connection.Prepare("INSERT INTO store (queue_name, client_id, object_type, object_id) VALUES(?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer stmtIns.Close()
	_, err = stmtIns.Exec(queueName, clientId, objectType, objectId)
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

type ClientQueuePair struct {
	clientId string
	queueName string
}

func getClients(objectType string, objectId string) []ClientQueuePair {
	rows, err := connection.Query("SELECT queue_name, client_id FROM store WHERE object_type=? AND object_id=?", objectType, objectId)
	if err != nil {
		fmt.Println(err)
		return []ClientQueuePair{}
	}

	var clientQueuePairs []ClientQueuePair

	for rows.Next() {
		var clientQueuePair ClientQueuePair
		rows.Scan(&clientQueuePair.queueName, &clientQueuePair.clientId)

		clientQueuePairs = append(clientQueuePairs, clientQueuePair)
	}

	return clientQueuePairs
}
