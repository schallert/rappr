package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/boltdb/bolt"
)

type TwiML struct {
	XMLName xml.Name `xml:"Response"`

	Message string `xml:",omitempty"`
}

type Votes struct {
	Derek uint
	Jay   uint
}

func mainHandler(w http.ResponseWriter, req *http.Request) {
	// Catch-all for root route
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	tally := Votes{0, 0}

	err := DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bname)
		if bucket == nil {
			return errors.New("Could not find bucket")
		}

		err := bucket.ForEach(func(key, data []byte) error {
			var r Response
			err := json.Unmarshal(data, &r)
			if err != nil {
				return err
			}

			if r.Body == "derek" {
				tally.Derek++
			}

			if r.Body == "jay" {
				tally.Jay++
			}

			return nil
		})

		if err != nil {
			elog.Println("Error tallying votes: ", err)
			return err
		}

		return nil
	})

	tmpl, err := template.New("index").Parse(INDEX_HTML)
	if err != nil {
		elog.Println("Error parsing index template: ", err)
		return
	}

	err = tmpl.Execute(w, tally)
	if err != nil {
		elog.Println("Error executing template: ", err)
		return
	}
}

func smsHandler(w http.ResponseWriter, req *http.Request) {
	parsedBody := strings.ToLower(strings.Trim(req.PostFormValue("Body"), " "))

	if parsedBody != "derek" && parsedBody != "jay" {
		twiml := TwiML{Message: `Valid votes are "Jay" or "Derek"`}

		w.Header().Set("Content-Type", "application/xml")
		if err := xml.NewEncoder(w).Encode(twiml); err != nil {
			elog.Println("Error encoding xml response", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	res := Response{
		From: req.PostFormValue("From"),
		Body: parsedBody,
	}

	marshalled, err := json.Marshal(res)
	if err != nil {
		elog.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bname)
		if err != nil {
			return err
		}

		err = bucket.Put(k(res.From), marshalled)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		elog.Println("Error storing response: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	twiml := TwiML{Message: "Your vote has been recorded"}
	w.Header().Set("Content-Type", "application/xml")
	if err := xml.NewEncoder(w).Encode(twiml); err != nil {
		elog.Println("Error encoding xml response", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func votesHandler(w http.ResponseWriter, req *http.Request) {
	responses := []Response{}

	err := DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bname)
		if bucket == nil {
			return errors.New("Error retrieving bucket")
		}

		bucket.ForEach(func(k, data []byte) error {
			var r Response
			err := json.Unmarshal(data, &r)
			if err == nil {
				responses = append(responses, r)
			} else {
				elog.Println("Error unmarshaling: ", err)
			}

			return nil
		})

		return nil
	})

	if err != nil {
		elog.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)

}

func resetHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Query().Get("secret") != "14758f1afd44c09b7992073ccf00b43d" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := DB.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket(bname)
		tx.CreateBucketIfNotExists(bname)

		return nil
	})

	if err != nil {
		http.Error(w, "Error clearing DB", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "DB cleared")
}
