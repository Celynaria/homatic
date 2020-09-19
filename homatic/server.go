package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type Pair struct {
	DeviceID int
	UserID   int
}

func main() {
	fmt.Println("hello hometic : I'm Gopher!!")
	r := mux.NewRouter()
	r.Handle("/pair-device", PairDevice(createPairDevice{})).Methods(http.MethodPost)

	addr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	fmt.Println("addr:", addr)

	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	log.Println("starting...")
	log.Fatal(server.ListenAndServe())
}

func PairDevice(device Device) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pair Pair
		err := json.NewDecoder(r.Body).Decode(pair)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		defer r.Body.Close()
		err = device.Pair(pair)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		w.Write([]byte(`{"status":"active"}`))
	}
}

type Device interface {
	Pair(p Pair) error
}

type createPairDevice struct{}

func (createPairDevice) Pair(pair Pair) error {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("connect to database error", err)
	}
	defer db.Close()

	insetQuery := `INSERT INTO pairs VALUES ($1,$2);`
	_, err = db.Exec(insetQuery, pair.DeviceID, pair.UserID)

	return err
}
