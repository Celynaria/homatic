package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
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
	pair := new(Pair)
	_ = json.NewDecoder(r.Body).Decode(pair)
	marshal, _ := json.Marshal(pair)
	w.Write(marshal)
}
