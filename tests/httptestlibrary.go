package keyvaluestore

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
)

var url = "http://localhost:8080"

type HTTPTest struct {
	Name string
	Args HTTPArgs
	Want string
}

type HTTPArgs struct {
	Method   string // GET, POST, DELETE
	Endpoint string // /items, items/:key, /items/keys, /items/values
	Key      string
	Value    []byte
}

// Executes a specific request defined by a singular http test.
func (test *HTTPTest) ExecuteRequest() (string, error) {
	switch test.Args.Method {
	case "GET":
		return getMessage(&test.Args)
	case "POST":
		return postMessage(&test.Args)
	case "DELETE":
		return deleteMessage(&test.Args)
	default:
		return "", errors.New("invalid http method provided")
	}
}

// Builds the json formatted post request body.
func (args *HTTPArgs) GetJSONString() string {
	return fmt.Sprintf(`{"key":"%s","value":%s}`, args.Key, args.Value)
}

// Public method that handles a series of http tests.
func HandleHTTPTests(t *testing.T, tests []HTTPTest) {
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			res, err := test.ExecuteRequest()
			if err != nil {
				t.Fatalf("Error %s encountered while executing test %s", err, test.Name)
			}

			// If want is not empty, and response is not equal to want
			if test.Want != "" && res != test.Want {
				t.Errorf("Expected %s, got %s", test.Want, res)
			}
		})
	}
}

// Method to interface the POST endpoint.
func postMessage(args *HTTPArgs) (string, error) {
	jsonData := []byte(args.GetJSONString())

	// Create a new HTTP request with context
	req, err := http.NewRequestWithContext(context.Background(), "POST", url+args.Endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %s", err)
		return "", err
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	return sendReqAndGetResp(req)
}

// Method to interface the GET endpoints.
func getMessage(args *HTTPArgs) (string, error) {
	// Create a new HTTP request with context
	req, err := http.NewRequestWithContext(context.Background(), "GET", url+args.Endpoint, nil)
	if err != nil {
		log.Printf("Error creating request: %s", err)
		return "", err
	}

	return sendReqAndGetResp(req)
}

// Method to interface the DELETE endpoints.
func deleteMessage(args *HTTPArgs) (string, error) {
	// Build the request with context
	req, err := http.NewRequestWithContext(context.Background(), "DELETE", url+args.Endpoint, nil)
	if err != nil {
		log.Printf("Error building request: %s", err)
		return "", err
	}

	return sendReqAndGetResp(req)
}

// Method to send http requests and read response bodies.
func sendReqAndGetResp(req *http.Request) (string, error) {
	// Create an HTTP client and send the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %s", err)
		return "", err
	}
	defer res.Body.Close()

	// Read response
	temp, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response body: %s", err)
		return "", err
	}

	return string(temp), nil
}
