package main

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

const (
	API_URL      = "https://gitlab.com/api/v4/projects/"
	PROJECT_ROOT = "/repository/files/"
)

var (
	httpPort  string
	httpsPort string
	debugMode bool
)

func init() {
	flag.StringVar(&httpPort, "http-port", "80", "Custom port to listen for HTTP requests")
	flag.StringVar(&httpsPort, "https-port", "443", "Custom port to listen for HTTPS requests")
	flag.BoolVar(&debugMode, "debug", false, "Verbose logs for debugging")
}

func main() {
	fmt.Println("ChartLab is starting")
	flag.Parse()

	c := make(chan error)
	http.HandleFunc("/", handleHelmRequest)

	go func() {
		fmt.Println(fmt.Sprintf("Listening for HTTP on port: %v", httpPort))
		c <- http.ListenAndServe(fmt.Sprintf(":%v", httpPort), nil)
	}()

	_, crtErr := os.Open("tls/tls.crt")
	_, keyErr := os.Open("tls/tls.key")

	// Listen and Serve TLS only if keys exist
	if !errors.Is(crtErr, os.ErrNotExist) && !errors.Is(keyErr, os.ErrNotExist) {
		go func() {
			fmt.Println(fmt.Sprintf("Listening for HTTPS on port: %v", httpsPort))
			c <- http.ListenAndServeTLS(fmt.Sprintf(":%v", httpsPort), "tls/tls.crt", "tls/tls.key", nil)
		}()
	}

	// Block forever, or until one of our listeners exits with an error
	err := <-c

	fmt.Println(err)
}

func handleHelmRequest(res http.ResponseWriter, req *http.Request) {
	if debugMode {
		fmt.Println("Incoming request:")
		printRequest(req)
	}

	// Convert the Helm request into a GitLab API request
	apiReq, err := convertRequest(req)

	if debugMode {
		fmt.Println("API request: ")
		printRequest(apiReq)
	}

	if err != nil {
		errMsg := "Failed to convert request to API call: " + err.Error()
		fmt.Println(errMsg)
		sendErrorResponse(res, http.StatusBadRequest, errMsg)
		return
	}

	// Use the converted request to make a GitLab API call
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

	// Forward the API response body which contains the requested file to Helm
	fmt.Fprint(res, string(responseBody))
}

func convertRequest(req *http.Request) (*http.Request, error) {
	// Example path: /<project-id>/path/to/file.yaml
	// Ignore the first element which is and empty string
	splitPath := strings.Split(req.URL.Path, "/")[1:]

	// Path must have at least a project-id and a file name
	if len(splitPath) < 2 {
		return nil, fmt.Errorf("Invalid URL. Use 'http://<node-ip>:<node-port>/<gitlab-project-id>'")
	}

	projectID := splitPath[0]
	filePath := strings.Join(splitPath[1:], "%2F")

	// Create a URL for GitLab API
	// Example: https://gitlab.com/api/v4/projects/<project-id>/repository/files/<file-path>/raw
	url := API_URL + projectID + PROJECT_ROOT + filePath + "/raw"
	apiReq, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("Failed to create HTTP request: " + err.Error())
	}

	// Get the base64 encoded Authorization header which contains "Basic encodedtoken"
	encodedToken := req.Header.Get("Authorization")

	if encodedToken == "" || len(strings.Split(encodedToken, " ")) < 2 {
		return nil, fmt.Errorf("Failed to convert auth token. You must provide a username and password.")
	}

	// The base64 encoded string needs to be decoded for the GitLab API request
	// Decoded string will be "username:token"
	token, err := base64.StdEncoding.DecodeString(strings.Split(encodedToken, " ")[1])

	if err != nil {
		return nil, fmt.Errorf("Failed to convert auth token: " + err.Error())
	}

	// GitLab expects an unencoded string containing our PAT
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
	reqDump, err := httputil.DumpRequest(req, true)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(reqDump))
}
