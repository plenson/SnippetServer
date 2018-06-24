// Copyright Peter Lenson.
// All Rights Reserved

package routes

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dbpkg "github.com/textioHQ/interview-peter-lenson/database"
	"log"
	"net/http"
)

// Maximum snippet length
const MAXIMUM_SNIPPET_LEN = 300

// Snippet type definition
type Snippet struct {
	ID     string `json:"id,omitempty"`
	Text   string `json:"text,omitempty"`
	Shared bool   `json:"shared,omitempty"`
}

// ALERT
// is a package level global

//#IF BUILD BOLT
var sniptDB dbpkg.SBolt

// Initialize route connection to database
//#IF BUILD BOLT
func InitBolt(snipsdb dbpkg.SBolt) {
	sniptDB = snipsdb
}

//#IF BUILD BOW
//var sniptDB dbpkg.SBow

//#IF BUILD BOW
//func InitBow(snipsdb dbpkg.SBow) {
//	sniptDB = snipsdb
//}

// Gets a snippet from the server given its id.
//
// It is invoked when the user calls (GET) the /snippets/{id} route"
// and will return it in a JSON format:
//
// {
//  "id":"0e5693ae-6449-4f67-88d3-9d18ae344300",
//  "text":"Mary had a little lamb.",
//  "shared":true
// }
func GetSnippetEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	get_requests_cnt.Inc()

	item, err := sniptDB.Get(params["id"], "Text")
	if err != nil {
		log.Printf("Couldn't find snippet %v", err)
	} else {
		json.NewEncoder(w).Encode(Snippet{params["id"], item, true})
		return
	}

	json.NewEncoder(w).Encode(&Snippet{})
}

type NumItems struct {
	Items int `json:"items,omitempty"`
}

// Gets all snippets from the server.
//
// It is invoked when the user calls (GET) the /snippets route"
// and will return all snippets as a list of separate JSON formatted statements:
//
// {"items":3}
// {"id":"0d725943-f06a-4747-bf78-3bdf70572087","text":"Little lamb","shared":true}
// {"id":"0e5693ae-6449-4f67-88d3-9d18ae344300","text":"Mary had a little lamb.","shared":true}
// {"id":"107f0644-f6be-4762-a7a6-5e6bcf869b7e","text":"Little lamb","shared":true}
func GetSnippetsEndpoint(w http.ResponseWriter, r *http.Request) {
	getall_requests_cnt.Inc()
	items, err := sniptDB.GetAll()
	if err != nil {
		log.Printf("Couldn't find snippets %v", err)
	} else {
		json.NewEncoder(w).Encode(NumItems{len(items)})
		for _, item := range items {

			snip, err := sniptDB.Get(item, "Text")
			if err != nil {
				log.Printf("Problem getting items  %v", err)
			} else {

				json.NewEncoder(w).Encode(Snippet{item, snip, true})
			}
		}

		return

		//json.NewEncoder(w).Encode(items)
		//return
	}
	//json.NewEncoder(w).Encode(snippets)
}

// Creates a snippet on the server.
//
// It is invoked when the user calls (POST) the /snippet/ route
// with body containing the snippet in the example format below.
//
// {
//  "text":"Mary had a little lamb.",
// }
//
// A successful post will return the id of the snippet, shared status of true.
//
//{"id":"b7c04ec2-1307-478b-9468-55e62c77245d","shared":true}
//
func CreateSnippetEndpoint(w http.ResponseWriter, r *http.Request) {
	set_requests_cnt.Inc()
	var snippet Snippet
	_ = json.NewDecoder(r.Body).Decode(&snippet)

	if len(snippet.Text) <= MAXIMUM_SNIPPET_LEN {
		id, err := sniptDB.Set("Text", snippet.Text)
		if err != nil {
			http.Error(w, "Failed to write snippet!", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(Snippet{id, "", true})
	} else {

		log.Printf("Snipe length %d", len(snippet.Text))
		http.Error(w, "Snippets must be less than 300 characters!", http.StatusRequestEntityTooLarge)
		return
	}

}

// Deletes a snippet from the server given its id.
//
// It is invoked when the user calls (DELETE) the /snippets/{id} route"
// and will return the id of deleted snippet in a JSON format:
//
// { "id":"0e5693ae-6449-4f67-88d3-9d18ae344300", "shared":true}
//
func DeleteSnippetEndpoint(w http.ResponseWriter, r *http.Request) {
	del_requests_cnt.Inc()
	params := mux.Vars(r)

	err := sniptDB.Del(params["id"])
	if err != nil {
		log.Printf("Couldn't find snippet %v", err)
	} else {
		json.NewEncoder(w).Encode(Snippet{params["id"], "", true})
		return
	}

}

// Status handler used for monitoring health of server.
//
// It is invoked when the user calls the /status route
// and will return a string with the message "API is up and running".
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	status_requests_cnt.Inc()

	w.Write([]byte("API is up and running"))
}

var (
	status_requests_cnt = promauto.NewCounter(
		prometheus.CounterOpts{Name: "status_handler_total", Help: "Status Handler requested."})

	get_requests_cnt = promauto.NewCounter(
		prometheus.CounterOpts{Name: "get_handler_total", Help: "Get Handler requested."})

	getall_requests_cnt = promauto.NewCounter(
		prometheus.CounterOpts{Name: "getall_handler_total", Help: "GetAll Handler requested."})

	set_requests_cnt = promauto.NewCounter(
		prometheus.CounterOpts{Name: "set_handler_total", Help: "Set Handler requested."})

	del_requests_cnt = promauto.NewCounter(
		prometheus.CounterOpts{Name: "del_handler_total", Help: "Del Handler requested."})
)

// Sets up route handlers.
func SetUpRouteHandlers() *mux.Router {
	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/snippets", GetSnippetsEndpoint).Methods("GET")
	router.HandleFunc("/snippet/{id}", GetSnippetEndpoint).Methods("GET")
	router.HandleFunc("/snippet/", CreateSnippetEndpoint).Methods("POST")
	router.HandleFunc("/snippet/{id}", DeleteSnippetEndpoint).Methods("DELETE")
	router.HandleFunc("/status", StatusHandler).Methods("GET")
	return router
}
