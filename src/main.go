package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const ORG = "dfds"
const ADO_APIVERSION = "6.1-preview.6"

func main() {
	// Acquire ADO token from environment variable "ADO_TOKEN"
	token := os.Getenv("ADO_TOKEN")
	if token == "" {
		panic("A valid Azure DevOps access token needs to be set in the environment variable: 'ADO_TOKEN'")
	}

	// Routing
	r := mux.NewRouter()
	r.HandleFunc("/builds", GetBuilds)

	r.Use(corsHandler(r))

	// Launch HTTP server
	println("HTTP server listening on :8080")
	if err := http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, handlers.CompressHandler(r))); err != nil {
		log.Fatal(err)
	}
}

// HTTP route handler : /builds
func GetBuilds(w http.ResponseWriter, r *http.Request) {
	// Read HTTP request body
	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Unable to read request body, sending 500 response")
		w.WriteHeader(500)
	}

	// Try to parse JSON into 'getBuildRequest' struct
	reqPayload := getBuildRequest{}
	err = json.Unmarshal(rawBody, &reqPayload)
	if err != nil {
		log.Println("Unable to parse JSON")
		log.Println(err)
		w.WriteHeader(400)
	}

	// Put together the API endpoint
	url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/build/builds?%s", ORG, reqPayload.Project, ADO_APIVERSION)

	// Create API request
	httpCli := http.DefaultClient
	adoReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
	}

	// Get ADO token
	token := os.Getenv("ADO_TOKEN")

	// Set necessary API request headers
	adoReq.Header.Set("Authorization", fmt.Sprintf("Basic %s", encodeToBase64(token)))
	adoReq.Header.Set("Content-Type", "application/json")
	adoReq.URL.Query().Set("queryOrder", "startTimeDescending")
	adoReq.URL.Query().Set("api-version", ADO_APIVERSION)

	// Send API request
	resp, err := httpCli.Do(adoReq)
	if err != nil {
		log.Println("Request to ADO failed")
		log.Println(err)
		w.WriteHeader(500)
	}

	// Read API response body
	respRawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read response body, sending 500 response")
		w.WriteHeader(500)
	}

	// Return API response
	w.Write(respRawBody)
}

func encodeToBase64(val string) string {
	return base64.StdEncoding.EncodeToString([]byte(val))
}

// Ensure the correct Access-Control-Allow-Origin header is set.
// Currently rather primitive and hardcoded, but it'll suffice for now.
func corsHandler(r *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			switch req.Host {
			case "backstage.dfds.cloud":
				w.Header().Set("Access-Control-Allow-Origin", "backstage.dfds.cloud")
			case "localhost:8080":
				w.Header().Set("Access-Control-Allow-Origin", "localhost:8080")
			case "localhost:3000":
				w.Header().Set("Access-Control-Allow-Origin", "localhost:8080")
			}

			next.ServeHTTP(w, req)
		})
	}
}

// API request payload
type getBuildRequest struct {
	Project string `json:"project"`
}