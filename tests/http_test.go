package keyvaluestore_test

import (
	"testing"

	kvs "github.com/defoeam/kvs/kv"
	kvstesting "github.com/defoeam/kvs/tests"
)

var serverRunning = false

// Utility method which starts the kvs server on a seperate thread if it isn't already running.
// This method should be called at the start of each test set.
func ensureServerRunning() {
	if !serverRunning {
		go kvs.StartServer()
		serverRunning = true
	}
}

func TestSet1(t *testing.T) {
	ensureServerRunning()

	tests := []kvstesting.HTTPTest{
		{
			Name: "Populate the store, pt. 1",
			Args: kvstesting.HTTPArgs{Method: "POST", Endpoint: "/items", Key: "0", Value: []byte(`"Hello World!"`)},
			Want: `{"key":"0","value":"Hello World!"}`,
		},
		{
			Name: "Populate the store, pt. 2",
			Args: kvstesting.HTTPArgs{Method: "POST", Endpoint: "/items", Key: "1", Value: []byte("[53,58,203]")},
			Want: `{"key":"1","value":[53,58,203]}`,
		},
		{
			Name: "Get all items, pt. 1",
			Args: kvstesting.HTTPArgs{Method: "GET", Endpoint: "/items", Key: "", Value: []byte{}},
			Want: "",
		},
		{
			Name: "Get item at key 0",
			Args: kvstesting.HTTPArgs{Method: "GET", Endpoint: "/items/0", Key: "", Value: []byte{}},
			Want: `{"key":"0","value":"Hello World!"}`,
		},
		{
			Name: "Delete item at key 0",
			Args: kvstesting.HTTPArgs{Method: "DELETE", Endpoint: "/items/0", Key: "", Value: []byte{}},
			Want: `{"key":"0","value":"Hello World!"}`,
		},
		{
			Name: "Get all items, pt. 2",
			Args: kvstesting.HTTPArgs{Method: "GET", Endpoint: "/items", Key: "", Value: []byte{}},
			Want: "",
		},
	}

	kvstesting.HandleHTTPTests(t, tests)
}
