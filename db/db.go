// This package manages all database access.
package db

import (
	"log"
	"time"

	"github.com/boltdb/bolt"
)

// The users database, only used by the server
var Users *bolt.DB
var userEmails = []byte("emails")
var userSalts = []byte("salts")
var userPassHashes = []byte("hashes")

// The messages database, used by both client and server
var Messages *bolt.DB

// Open the users database and initialize if necessary
func UserDBInit(path string) error {
	var err error
	Users, err = bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Panic(err)
	}
	//return Users.Update(func(tx *bolt.Tx) error {
	//	tx.CreateBucketIfNotExists(
}
