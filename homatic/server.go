package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Celynaria/Homatic/homatic/logger"
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

type CustomResponseWriter interface {
	JSON(statusCode int, data interface{})
}

type JSONResponseWriter struct {
	http.ResponseWriter
}

type CustomHandlerFunc func(CustomResponseWriter, *http.Request)

func (handler CustomHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler(&JSONResponseWriter{w}, r)
}

func (w *JSONResponseWriter) JSON(statusCode int, data interface{}) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func main() {
	err := run()
	if err != nil {
		log.Fatal("can't start application", err)
	}
}

func run() error {
	fmt.Println("hello hometic : I'm Gopher!!")

	r := mux.NewRouter()
	r.Use(logger.MiddleWare)
	r.Handle("/pair-device", PairDevice(CreatePairDeviceFunc(createPairDevice))).Methods(http.MethodPost)

	addr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	fmt.Println("addr:", addr)

	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	log.Println("starting...")
	return server.ListenAndServe()
}

func PairDevice(device Device) CustomHandlerFunc {
	return func(w CustomResponseWriter, r *http.Request) {
		l := logger.GetLog(r.Context())
		l.Info("pair-device")

		var pair Pair
		err := json.NewDecoder(r.Body).Decode(&pair)
		if err != nil {
			l.Error(err.Error())
			w.JSON(http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()
		err = device.Pair(pair)

		if err != nil {
			l.Error(err.Error())
			w.JSON(http.StatusBadRequest, err.Error())
			return
		}

		w.JSON(http.StatusOK, map[string]interface{}{"status": "active"})
	}
}

type Device interface {
	Pair(p Pair) error
}

type CreatePairDeviceFunc func(p Pair) error

func (fn CreatePairDeviceFunc) Pair(p Pair) error {
	return fn(p)
}

func createPairDevice(pair Pair) error {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("connect to database error", err)
	}
	defer db.Close()

	insetQuery := `INSERT INTO pairs VALUES ($1,$2);`
	_, err = db.Exec(insetQuery, pair.DeviceID, pair.UserID)

	return err
}
