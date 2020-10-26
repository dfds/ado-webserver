package helpers

import (	
	"sort"
	"strings"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	types "ado-pipeline/types"
)

func EncodeToBase64(val string) string {
	return base64.StdEncoding.EncodeToString([]byte(val))
}

func TransformBuildsApiResponse(rawResp []byte) ([]types.SortedBuildResponse, error) {	
	parsed := types.GetBuildResponse{}

	println(string(rawResp));

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
	var values []types.Build

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
	var payload []types.SortedBuildResponse
	for _, k := range values {
		dto := types.SortedBuildResponse{
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