package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
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
	url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/build/builds", ORG, reqPayload.Project)

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

	transformResp, err := transformBuildsApiResponse(respRawBody)
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

func encodeToBase64(val string) string {
	return base64.StdEncoding.EncodeToString([]byte(val))
}

func transformBuildsApiResponse(rawResp []byte) ([]sortedBuildResponse, error) {
	parsed := getBuildResponse{}
	err := json.Unmarshal(rawResp, &parsed)
	if err != nil {
		return nil, errors.New("unable to parse JSON response")
	}

	// Get-LatestBuilds replica
	// Sort by name
	sort.SliceStable(parsed.Build, func(i ,j int) bool {
		var si string = parsed.Build[i].Definition.Name
		var sj string = parsed.Build[j].Definition.Name
		var si_lower = strings.ToLower(si)
		var sj_lower = strings.ToLower(sj)
		if si_lower == sj_lower {
			return si < sj
		}
		return si_lower < sj_lower
	})

	// Ensure uniqueness
	uniqueness := make(map[string]string)
	var values []Build

	for _, k := range parsed.Build {
		_, exists := uniqueness[k.Definition.Name]
		if !exists {
			uniqueness[k.Definition.Name] = ""
			values = append(values, k)
		}
	}

	// Sort by queuetime descending
	sort.SliceStable(values, func(i, j int) bool {
		return values[i].QueueTime.After(values[j].QueueTime)
	})

	// Put in DTO
	var payload []sortedBuildResponse
	for _, k := range values {
		dto := sortedBuildResponse{
			Status:       k.Status,
			Result:       k.Result,
			BuildNumber:  k.BuildNumber,
			QueueTime:    k.QueueTime,
			StartTime:    k.StartTime,
			FinishTime:   k.FinishTime,
			SourceBranch: k.SourceBranch,
			PipelineName: k.Definition.Name,
			ProjectId:    k.Project.ID,
			BuildPageLink: fmt.Sprintf("https://dev.azure.com/dfds/%s/_build/results?buildId=%d&view=results", k.Project.ID, k.ID),
		}
		payload = append(payload, dto)
	}

	return payload, nil
}

// Ensure the correct Access-Control-Allow-Origin header is set.
// Currently rather primitive and hardcoded, but it'll suffice for now.
func corsHandler(r *mux.Router) mux.MiddlewareFunc {
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

// API request payload
type getBuildRequest struct {
	Project string `json:"project"`
}

type getBuildResponse struct {
	Count int     `json:"count"`
	Build []Build `json:"value"`
}

type sortedBuildResponse struct {
	Status string `json:"status"`
	Result string `json:"result"`
	BuildNumber string `json:"buildNumber"`
	QueueTime time.Time `json:"queueTime"`
	StartTime time.Time `json:"startTime"`
	FinishTime time.Time `json:"finishTime"`
	SourceBranch string `json:"sourceBranch"`
	PipelineName string `json:"pipelineName"` // Definition.Name
	ProjectId string `json:"projectId"` // Project.Id
	BuildPageLink string `json:"buildPageLink"`
}

type Build struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Web struct {
			Href string `json:"href"`
		} `json:"web"`
		SourceVersionDisplayURI struct {
			Href string `json:"href"`
		} `json:"sourceVersionDisplayUri"`
		Timeline struct {
			Href string `json:"href"`
		} `json:"timeline"`
		Badge struct {
			Href string `json:"href"`
		} `json:"badge"`
	} `json:"_links"`
	Properties struct {
	} `json:"properties"`
	Tags              []interface{} `json:"tags"`
	ValidationResults []interface{} `json:"validationResults"`
	Plans             []struct {
		PlanID string `json:"planId"`
	} `json:"plans"`
	TriggerInfo struct {
	} `json:"triggerInfo"`
	ID          int       `json:"id"`
	BuildNumber string    `json:"buildNumber"`
	Status      string    `json:"status"`
	Result      string    `json:"result"`
	QueueTime   time.Time `json:"queueTime"`
	StartTime   time.Time `json:"startTime"`
	FinishTime  time.Time `json:"finishTime"`
	URL         string    `json:"url"`
	Definition  struct {
		Drafts      []interface{} `json:"drafts"`
		ID          int           `json:"id"`
		Name        string        `json:"name"`
		URL         string        `json:"url"`
		URI         string        `json:"uri"`
		Path        string        `json:"path"`
		Type        string        `json:"type"`
		QueueStatus string        `json:"queueStatus"`
		Revision    int           `json:"revision"`
		Project     struct {
			ID             string    `json:"id"`
			Name           string    `json:"name"`
			Description    string    `json:"description"`
			URL            string    `json:"url"`
			State          string    `json:"state"`
			Revision       int       `json:"revision"`
			Visibility     string    `json:"visibility"`
			LastUpdateTime time.Time `json:"lastUpdateTime"`
		} `json:"project"`
	} `json:"definition"`
	Project struct {
		ID             string    `json:"id"`
		Name           string    `json:"name"`
		Description    string    `json:"description"`
		URL            string    `json:"url"`
		State          string    `json:"state"`
		Revision       int       `json:"revision"`
		Visibility     string    `json:"visibility"`
		LastUpdateTime time.Time `json:"lastUpdateTime"`
	} `json:"project"`
	URI           string `json:"uri"`
	SourceBranch  string `json:"sourceBranch"`
	SourceVersion string `json:"sourceVersion"`
	Queue         struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Pool struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			IsHosted bool   `json:"isHosted"`
		} `json:"pool"`
	} `json:"queue"`
	Priority     string `json:"priority"`
	Reason       string `json:"reason"`
	RequestedFor struct {
		DisplayName string `json:"displayName"`
		URL         string `json:"url"`
		Links       struct {
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
		} `json:"_links"`
		ID         string `json:"id"`
		UniqueName string `json:"uniqueName"`
		ImageURL   string `json:"imageUrl"`
		Descriptor string `json:"descriptor"`
	} `json:"requestedFor"`
	RequestedBy struct {
		DisplayName string `json:"displayName"`
		URL         string `json:"url"`
		Links       struct {
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
		} `json:"_links"`
		ID         string `json:"id"`
		UniqueName string `json:"uniqueName"`
		ImageURL   string `json:"imageUrl"`
		Descriptor string `json:"descriptor"`
	} `json:"requestedBy"`
	LastChangedDate time.Time `json:"lastChangedDate"`
	LastChangedBy   struct {
		DisplayName string `json:"displayName"`
		URL         string `json:"url"`
		Links       struct {
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
		} `json:"_links"`
		ID         string `json:"id"`
		UniqueName string `json:"uniqueName"`
		ImageURL   string `json:"imageUrl"`
		Descriptor string `json:"descriptor"`
	} `json:"lastChangedBy"`
	OrchestrationPlan struct {
		PlanID string `json:"planId"`
	} `json:"orchestrationPlan"`
	Logs struct {
		ID   int    `json:"id"`
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"logs"`
	Repository struct {
		ID                 string      `json:"id"`
		Type               string      `json:"type"`
		Name               string      `json:"name"`
		URL                string      `json:"url"`
		Clean              interface{} `json:"clean"`
		CheckoutSubmodules bool        `json:"checkoutSubmodules"`
	} `json:"repository"`
	KeepForever       bool        `json:"keepForever"`
	RetainedByRelease bool        `json:"retainedByRelease"`
	TriggeredByBuild  interface{} `json:"triggeredByBuild"`
}
