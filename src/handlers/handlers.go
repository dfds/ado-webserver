package handlers

import (	
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	helpers "ado-pipeline/helpers"
	types "ado-pipeline/types"
	identity "ado-pipeline/identity"
)

const ORG = "dfds"
const ADO_APIVERSION = "6.1-preview.6"

// HTTP route handler : /builds
func GetBuildsHandler(w http.ResponseWriter, r *http.Request) {
	// Read HTTP request body
	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Unable to read request body, sending 500 response")
		w.WriteHeader(500)
	}

	// Try to parse JSON into 'getBuildRequest' struct
	reqPayload := types.GetBuildRequest{}
	err = json.Unmarshal(rawBody, &reqPayload)
	if err != nil {
		log.Println("Unable to parse JSON")
		log.Println(err)
		w.WriteHeader(400)
	}

	// Put together the API endpoint
	url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/build/builds", ORG, reqPayload.Project)

	// Create API request
	httpCli := http.DefaultClient
	adoReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
	}

	// Get an access token
	//TODO: Our app registration probably needs some delegated permissions configured to call ADO.
	var token = identity.AcquireTokenClientSecret()

	// Set necessary API request headers
	adoReq.Header.Set("Authorization", fmt.Sprintf("Basic %s", helpers.EncodeToBase64(token)))
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

	transformResp, err := helpers.TransformBuildsApiResponse(respRawBody)
	if err != nil {
		log.Println("Unable to transform ADO response, sending 500 response")
		log.Println(err)
		w.WriteHeader(500)
	}

	payload, err := json.Marshal(transformResp)
	if err != nil {
		log.Println("Unable to serialize ADO transformed response to JSON")
		log.Println(err)
		w.WriteHeader(500)
	}

	w.Header().Set("Access-Control-Allow-Credentials", "true")
	// Return API response
	w.Write(payload)
}

// Ensure the correct Access-Control-Allow-Origin header is set.
func CorsHandler(r *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			switch req.Header.Get("origin") {
			case "http://backstage.dfds.cloud":
				w.Header().Set("Access-Control-Allow-Origin", "https://backstage.dfds.cloud")
			case "http://localhost:7000":
				w.Header().Set("Access-Control-Allow-Origin", "http://localhost:7000")
			case "http://localhost:3000":
				w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			default:
				w.Header().Set("Access-Control-Allow-Origin", "https://backstage.dfds.cloud")
			}

			next.ServeHTTP(w, req)
		})
	}
}