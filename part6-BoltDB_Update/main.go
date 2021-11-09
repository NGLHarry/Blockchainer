package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

//database name
const dbfile = "blockchain.db"

//bucket
const blocksBucket = "blocks"

func main() {
	//---------create the database---------
	// If the database exists, open it, if not, create a database
	db, err := bolt.Open(dbfile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()

	//insert/update database
	db.Update(func(tx *bolt.Tx) error {
		//Determine if this table exists in the database
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			fmt.Println("No existing blockchain found.Creating a new one ....")
			//create Bucket
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}
			//save data
			// key []byte, value []byte
			err = b.Put([]byte("MaoBuyi"), []byte("https://music.163.com/#/artist?id=12138269"))
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
