package main

import (
	"log"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
	"github.com/codegangsta/negroni"
)

type Response struct {
	From string
	Body string
}

var (
	bname = []byte("responses")
	elog  = log.New(os.Stderr, "[error] ", 0)
	DB    *bolt.DB
)

func k(s string) []byte {
	return []byte(s)
}

func main() {
	var err error
	DB, err = bolt.Open("rappr.boltdb", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Initial bucket creation
	err = DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bname)
		return err
	})

	if err != nil {
		log.Fatal("Error initializing bucket: ", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandler)
	mux.HandleFunc("/sms", smsHandler)
	mux.HandleFunc("/votes", votesHandler)
	mux.HandleFunc("/reset", resetHandler)

	server := negroni.Classic()
	server.UseHandler(mux)
	server.Run(":3001")
}
