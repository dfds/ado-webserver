package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const ORG = "dfds"
const ADO_APIVERSION = "6.1-preview.6"

func main() {
	token := os.Getenv("ADO_TOKEN")
	if token == "" {
		panic("A valid Azure DevOps access token needs to be set in the environment variable: 'ADO_TOKEN'")
	}

	r := mux.NewRouter()
	r.HandleFunc("/builds", GetBuilds)

	println("HTTP server listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

func GetBuilds(w http.ResponseWriter, r *http.Request) {
	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Unable to read request body, sending 500 response")
		w.WriteHeader(500)
	}

	println(string(rawBody))

	reqPayload := getBuildRequest{}
	err = json.Unmarshal(rawBody, &reqPayload)
	if err != nil {
		log.Println("Unable to parse JSON")
		log.Println(err)
		w.WriteHeader(400)
	}


	url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/build/builds?%s", ORG, reqPayload.Project, ADO_APIVERSION)

	httpCli := http.DefaultClient
	adoReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
	}

	token := os.Getenv("ADO_TOKEN")
	adoReq.Header.Set("Authorization", fmt.Sprintf("Basic %s", encodeToBase64(token)))
	adoReq.Header.Set("Content-Type", "application/json")
	adoReq.URL.Query().Set("queryOrder", "startTimeDescending")
	adoReq.URL.Query().Set("api-version", ADO_APIVERSION)

	resp, err := httpCli.Do(adoReq)
	if err != nil {
		log.Println("Request to ADO failed")
		log.Println(err)
		w.WriteHeader(500)
	}

	respRawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read response body, sending 500 response")
		w.WriteHeader(500)
	}

	println(string(respRawBody))
	w.Write(respRawBody)
}

func encodeToBase64(val string) string {
	return base64.StdEncoding.EncodeToString([]byte(val))
}

type getBuildRequest struct {
	Project string `json:"project"`
}