package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/comail/colog"
	"github.com/gorilla/mux"
	"github.com/undeadops/konveyer"
	"github.com/undeadops/konveyer/handler"
)

const appName = "konveyer-server"

func init() {
	colog.Register()
	colog.ParseFields(true)
	colog.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	runtime := konveyer.Initialize()
	defer runtime.Stop()

	log.Println("info: Server starting up on port: :", runtime.Port)
	router := mux.NewRouter()
	router.Handle("/project", handler.Handler{runtime, handler.GetProjects}).Methods("GET")
	router.Handle("/project", handler.Handler{runtime, handler.CreateProjects}).Methods("PUT")
	router.Handle("/project/{projectName}", handler.Handler{runtime, handler.DescribeProject}).Methods("GET")
	router.HandleFunc("/", Index).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+runtime.Port, router))

}

// Index - list index
func Index(w http.ResponseWriter, r *http.Request) {
	payload := "{ status: OK }"
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
