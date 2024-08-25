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

type HttpTest struct {
	name string
	args HttpArgs
	want string
}

type HttpArgs struct {
	method   string // GET, POST, DELETE
	endpoint string
	key      string
	value    string
}

func (test *HttpTest) ExecuteRequest() (string, error) {
	switch test.args.method {
	case "GET":
		// todo
	case "POST":
		return postMessage(test.args.GetJSONString())
	case "DELETE":
		// todo
	default:
		return "", errors.New("invalid http method provided")
	}

	return "", nil
}

func (args *HttpArgs) GetJSONString() string {
	return fmt.Sprintf(`{"key":"%s","value":"%s"}`, args.key, args.value)
}

func HandleTests(t *testing.T, tests []HttpTest) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.ExecuteRequest()
			if err != nil {
				t.Fatalf("Error %s encountered while executing test %s", err, test.name)
			}

			if res != test.want {
				t.Errorf("Expected %s, got %s", test.want, res)
			}
		})
	}
}

// Method to interface the /items POST endpoint.
func postMessage(jsonString string) (string, error) {
	jsonData := []byte(jsonString)

	url := "http://localhost:8080"

	// Build the request
	req, err := http.NewRequest("POST", url+"/items", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %s", err)
		return "", err
	}

	// Set the content type to application/json
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %s", err)
		return "", err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error sending request: %s", err)
		return "", err
	}

	bodyString := string(body)
	log.Printf("Response body: %s", bodyString)
	return bodyString, nil
}

// Method to interface the /items/ GET endpoint

// Method to interface the /items/:key GET endpoint

// Method to interface the /items/keys GET endpoint

// Method to interface the /items/values GET endpoint

// Method to interface the /items DELETE endpoint

// Method to interface the /items/:key DELETE endpoint
