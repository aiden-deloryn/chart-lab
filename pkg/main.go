package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	API_URL      = "https://gitlab.com/api/v4/projects/"
	PROJECT_ROOT = "/repository/files/"
)

func main() {
	http.HandleFunc("/", handleHelmRequest)

	log.Fatal(http.ListenAndServe(":9000", nil))
}

func handleHelmRequest(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Incoming request:")
	printRequest(req)

	apiReq, err := convertRequest(req)

	fmt.Println("API request: ")
	printRequest(apiReq)

	if err != nil {
		errMsg := "Failed to convert request to API call: " + err.Error()
		fmt.Println(errMsg)
		sendErrorResponse(res, http.StatusBadRequest, errMsg)
		return
	}

	apiRes, err := sendGitLabRequest(apiReq)

	if err != nil {
		errMsg := "Failed to send API request: " + err.Error()
		fmt.Println(errMsg)
		sendErrorResponse(res, http.StatusInternalServerError, errMsg)
		return
	}

	responseBody, err := ioutil.ReadAll(apiRes.Body)

	if err != nil {
		errMsg := "Failed to read response from API: " + err.Error()
		fmt.Println(errMsg)
		sendErrorResponse(res, http.StatusInternalServerError, errMsg)
		return
	}

	fmt.Fprint(res, string(responseBody))
}

func convertRequest(req *http.Request) (*http.Request, error) {
	splitPath := strings.Split(req.URL.Path, "/")[1:]

	if len(splitPath) < 2 {
		return nil, fmt.Errorf("Invalid URL. Use 'http://<host>:<port>/<gitlab-project-id>'")
	}

	projectID := splitPath[0]
	filePath := strings.Join(splitPath[1:], "%2F")

	url := API_URL + projectID + PROJECT_ROOT + filePath + "/raw"
	apiReq, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("Failed to create HTTP request: " + err.Error())
	}

	encodedToken := req.Header.Get("Authorization")

	if encodedToken == "" || len(strings.Split(encodedToken, " ")) < 2 {
		return nil, fmt.Errorf("Failed to convert auth token. You must provide a username and password.")
	}

	token, err := base64.StdEncoding.DecodeString(strings.Split(encodedToken, " ")[1])

	if err != nil {
		return nil, fmt.Errorf("Failed to convert auth token: " + err.Error())
	}

	apiReq.Header.Set("PRIVATE-TOKEN", strings.Split(string(token), ":")[1])

	return apiReq, nil
}

func sendGitLabRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: " + err.Error())
	}

	return res, err
}

func sendErrorResponse(res http.ResponseWriter, httpStatus int, message string) {
	res.WriteHeader(httpStatus)
	fmt.Fprint(res, message)
}

func printRequest(req *http.Request) {
	fmt.Println("=== PROTOCOL ===")
	fmt.Println(req.Proto)
	fmt.Println()
	fmt.Println("=== HOST ===")
	fmt.Println(req.Host)
	fmt.Println()
	fmt.Println("=== METHOD ===")
	fmt.Println(req.Method)
	fmt.Println()
	fmt.Println("=== PATH ===")
	fmt.Println(req.URL.Path)
	fmt.Println()
	// fmt.Println("=== HEADER ===")
	// fmt.Println(req.Header)
	// fmt.Println()
	fmt.Println("=== BODY ===")
	fmt.Println(req.Body)
	fmt.Println()
	fmt.Println("---" + "\n")
}
