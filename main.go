package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TaskInfo struct {
	Organization       string `json:"organization"`
	ID                 string `json:"id"`
	TaskType           string `json:"taskType"`
	ComponentID        string `json:"componentId"`
	ComponentKey       string `json:"componentKey"`
	ComponentName      string `json:"componentName"`
	ComponentQualifier string `json:"componentQualifier"`
	Status             string `json:"status"`
	SubmittedAt        string `json:"submittedAt"`
	StartedAt          string `json:"startedAt"`
	ExecutedAt         string `json:"executedAt"`
	ExecutionTimeMs    int64  `json:"executionTimeMs"`
	Logs               bool   `json:"logs"`
	ErrorMessage       string `json:"errorMessage"`
	HasErrorStacktrace string `json:"hasErrorStacktrace"`
	HasScannerContext  bool   `json:"hasScannerContext"`
}

type TaskStatusResponse struct {
	Tasks []TaskInfo `json:"tasks"`
}

func sonarAPIRequest(
	sonarHostURL string, apiMethod string, httpMethod string, username string, password string) (
	[]byte, error) {

	var responseData []byte
	client := &http.Client{}
	url := sonarHostURL + apiMethod
	componentRequest, err := http.NewRequest(
		httpMethod,
		url,
		nil)
	if err != nil {
		return responseData, err
	}

	componentRequest.SetBasicAuth(username, password)
	resp, err := client.Do(componentRequest)
	if err != nil {
		return responseData, err
	}
	defer resp.Body.Close()

	responseData, err = ioutil.ReadAll(resp.Body)
	return responseData, err
}

func apiCeActivityStatus(
	sonarHostURL string, projectKey string, username string, password string) (TaskStatusResponse, error) {

	var response TaskStatusResponse

	body, err := sonarAPIRequest(
		sonarHostURL,
		"api/ce/activity?status=PENDING,IN_PROGRESS&component="+projectKey,
		"GET",
		username,
		password)

	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err

	}
	return response, nil
}

func waitForPendingTasks(sonarHostURL string, projectKey string, username string, password string) bool {
	tasksCount := 1
	fmt.Println("Waiting for pending tasks to finish...")
	for tasksCount > 0 {
		response, err := apiCeActivityStatus(sonarHostURL, projectKey, username, password)
		if err != nil {
			fmt.Println("Failed to perform SonarQube API request for activities!")
			fmt.Println("Error: ", err)
			return false
		}
		tasksCount = len(response.Tasks)
		fmt.Printf("\r%d pending tasks remaining for %s component...", tasksCount, projectKey)
	}
	fmt.Println("\n\nNo more pending tasks left for component!\n")
	return true
}

func main() {
	if waitForPendingTasks("http://localhost:9000/", "externals", "12ec798e45e39d289780709b537cfc92a6f8bfa3", "") {
		fmt.Println("No pending tasks found!")

	}
}
