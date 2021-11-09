package main

import (
	"log"

	"github.com/boltdb/bolt"
)

//database name
const dbfile = "blockchain.db"

//bucket
const blocksBucket = "blocks"

func main() {
	db, err := bolt.Open(dbfile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
}
