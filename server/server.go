package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type Log struct {
	mu      sync.Mutex
	records []Record
}

type Record struct {
	Value  []byte `json:"value"`
	Offset int64  `json:"offset"`
}

func (l *Log) Append(record Record) (int64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	record.Offset = int64(len(l.records))
	l.records = append(l.records, record)

	return l.records[len(l.records)-1].Offset, nil
}

func (l *Log) Get(offset int64) (Record, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.records[offset], nil
}

func main() {
	r := mux.NewRouter()
	var log Log
	log_router := r.PathPrefix("/log").Subrouter()

	log_router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var record Record
		json.NewDecoder(r.Body).Decode(&record)

		if record.Value == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No value in request body"))
			return
		}

		offset, err := log.Append(record)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		json.NewEncoder(w).Encode(offset)
	}).Methods("POST")

	log_router.HandleFunc("/{offset}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		offset := vars["offset"]
		offset_int, err := strconv.ParseInt(offset, 10, 64)
		record, err_get := log.Get(offset_int)

		if err_get != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		json.NewEncoder(w).Encode(record)
	}).Methods("GET")

	http.ListenAndServe(":80", r)
}
