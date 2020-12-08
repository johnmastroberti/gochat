// This package manages all database access.
package db

import (
	"errors"
	"time"

	"github.com/boltdb/bolt"
	"golang.org/x/crypto/bcrypt"
)

// The users database, only used by the server
var users *bolt.DB
var userEmails = []byte("emails")
var userPassHashes = []byte("hashes")
var userNotFoundError = errors.New("User not in database")
var userExistsError = errors.New("User already exists in database")

// The messages database, used by both client and server
var Messages *bolt.DB

// Open the users database and initialize if necessary
func UserDBInit(path string) error {
	// Open the database
	var err error
	users, err = bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	// Create the email and hash buckets if they do not exist
	return users.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(userEmails)
		_, err := tx.CreateBucketIfNotExists(userPassHashes)
		return err
	})
}

// Add a new user to the database, along with their email and
// password hash
func AddNewUser(username string, email string, password string) error {
	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// Add it to the database
	err = users.Update(func(tx *bolt.Tx) error {
		emailBucket := tx.Bucket(userEmails)
		hashBucket := tx.Bucket(userPassHashes)
		// Make sure the user does not exist yet
		if hashBucket.Get([]byte(username)) != nil {
			return userExistsError
		}
		emailBucket.Put([]byte(username), []byte(email))
		return hashBucket.Put([]byte(username), hash)
	})
	return err
}

// Check whether or not a username and password is correct
// Returns true if the password is correct
func AuthenticateUser(username string, password string) bool {
	err := users.View(func(tx *bolt.Tx) error {
		// Retreive hash from database
		b := tx.Bucket(userPassHashes)
		hash := b.Get([]byte(username))
		if hash == nil {
			return userNotFoundError
		}
		// Check if the password matches
		e := bcrypt.CompareHashAndPassword(hash, []byte(password))
		return e
	})
	if err != nil {
		return false
	}
	return true
}
