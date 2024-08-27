package keyvaluestore

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
)

var url string = "http://localhost:8080"

type HttpTest struct {
	name string
	args HttpArgs
	want string
}

type HttpArgs struct {
	method   string // GET, POST, DELETE
	endpoint string // /items, items/:key, /items/keys, /items/values
	key      string
	value    []byte
}

// Executes a specific request defined by a singular http test
func (test *HttpTest) ExecuteRequest() (string, error) {
	switch test.args.method {
	case "GET":
		return getMessage(&test.args)
	case "POST":
		return postMessage(&test.args)
	case "DELETE":
		return deleteMessage(&test.args)
	default:
		return "", errors.New("invalid http method provided")
	}
}

// Builds the json formatted post request body
func (args *HttpArgs) GetJSONString() string {
	return fmt.Sprintf(`{"key":"%s","value":%s}`, args.key, args.value)
}

// Public method that handles a series of http tests.
func HandleHTTPTests(t *testing.T, tests []HttpTest) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.ExecuteRequest()
			if err != nil {
				t.Fatalf("Error %s encountered while executing test %s", err, test.name)
			}

			// If want is not empty, and response is not equal to want
			if test.want != "" && res != test.want {
				t.Errorf("Expected %s, got %s", test.want, res)
			}
		})
	}
}

// Method to interface the POST endpoint.
func postMessage(args *HttpArgs) (string, error) {
	jsonData := []byte(args.GetJSONString())

	// Build and send the request
	res, err := http.Post(url+args.endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending request: %s", err)
		return "", err
	}

	// Read the response body
	body, err := readResponseBody(res)
	if err != nil {
		log.Printf("Error reading response body: %s", err)
		return "", err
	}

	return string(body), nil
}

// Method to interface the GET endpoints.
func getMessage(args *HttpArgs) (string, error) {
	// Make request
	res, err := http.Get(url + args.endpoint)
	if err != nil {
		log.Printf("Error sending request: %s", err)
		return "", err
	}

	// Handle response
	body, err := readResponseBody(res)
	if err != nil {
		log.Printf("Error reading response: %s", err)
		return "", err
	}

	return string(body), nil
}

// Method to interface the DELETE endpoints.
func deleteMessage(args *HttpArgs) (string, error) {

	// Build the request
	req, err := http.NewRequest("DELETE", url+args.endpoint, nil)
	if err != nil {
		log.Printf("Error building request: %s", err)
		return "", err
	}

	// Send the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %s", err)
		return "", err
	}

	// Read response
	body, err := readResponseBody(res)
	if err != nil {
		log.Printf("Error reading response body: %s", err)
	}

	return body, nil
}

// Method to read response bodies
func readResponseBody(resp *http.Response) (body string, err error) {
	defer resp.Body.Close()
	temp, err := io.ReadAll(resp.Body)
	return string(temp), err
}
