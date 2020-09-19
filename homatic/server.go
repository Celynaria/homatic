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
	r.HandleFunc("/pair-device", pairDeviceHandler).Methods(http.MethodPost)

	addr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	fmt.Println("addr:", addr)

	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	log.Println("starting...")
	log.Fatal(server.ListenAndServe())
}

func pairDeviceHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("connect to database error", err)
	}
	defer db.Close()

	pair := new(Pair)
	err = json.NewDecoder(r.Body).Decode(pair)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	defer r.Body.Close()

	insetQuery := `INSERT INTO pairs VALUES ($1,$2);`
	_, err = db.Exec(insetQuery, pair.DeviceID, pair.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.Write([]byte(`{"status":"active"}`))
}
