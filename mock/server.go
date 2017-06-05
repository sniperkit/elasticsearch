package mock

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/oklog/ulid"
	"time"
	"math/rand"
)

var (
	// mapping of index:type:documentID:document
	database *store

	// create simple ULID using basic source of entropy
	timestamp = time.Unix(1000000, 0)
	entropy = rand.New(rand.NewSource(timestamp.UnixNano()))
)

func ULID() string {
	return ulid.MustNew(ulid.Timestamp(timestamp), entropy).String()
}

func New() *http.Server {
	database = newStore()
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/{index}", DeleteIndex).Methods("DELETE")
	router.HandleFunc("/{index}/_search", SearchIndex).Methods("GET")
	router.HandleFunc("/{index}/{_type}/_search", SearchType).Methods("GET")
	router.HandleFunc("/{index}/{_type}/", InsertDocument).Methods("POST")
	router.HandleFunc("/{index}/{_type}/{id}", GetDocumentByID).Methods("GET")
	router.HandleFunc("/{index}/{_type}/{id}/", UpdateDocumentByID).Methods("PUT")
	router.HandleFunc("/{index}/{_type}/{id}/", DeleteDocumentByID).Methods("DELETE")

	return &http.Server{
		Handler: router,
		Addr: "127.0.0.1:9201",
	}
}