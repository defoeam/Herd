package keyvaluestore_test

import (
	"net/http"
	"testing"

	kvs "github.com/defoeam/kvs/internal"
)

func TestSet1(t *testing.T) {
	// Start the kvs server
	go kvs.StartServer()

	tests := []kvs.HTTPTest{
		{
			Name: "Populate the store, pt. 1",
			Args: kvs.HTTPArgs{
				Method:   http.MethodPost,
				Endpoint: "/items",
				Key:      "0",
				Value:    []byte(`"Hello World!"`)},
			Want: `{"key":"0","value":"Hello World!"}`,
		},
		{
			Name: "Populate the store, pt. 2",
			Args: kvs.HTTPArgs{
				Method:   http.MethodPost,
				Endpoint: "/items",
				Key:      "1",
				Value:    []byte("[53,58,203]")},
			Want: `{"key":"1","value":[53,58,203]}`,
		},
		{
			Name: "Get all items, pt. 1",
			Args: kvs.HTTPArgs{
				Method:   http.MethodGet,
				Endpoint: "/items",
				Key:      "",
				Value:    []byte{}},
			Want: "",
		},
		{
			Name: "Get item at key 0",
			Args: kvs.HTTPArgs{
				Method:   http.MethodGet,
				Endpoint: "/items/0",
				Key:      "",
				Value:    []byte{}},
			Want: `{"key":"0","value":"Hello World!"}`,
		},
		{
			Name: "Delete item at key 0",
			Args: kvs.HTTPArgs{
				Method:   http.MethodDelete,
				Endpoint: "/items/0",
				Key:      "",
				Value:    []byte{}},
			Want: `{"key":"0","value":"Hello World!"}`,
		},
		{
			Name: "Get all items, pt. 2",
			Args: kvs.HTTPArgs{
				Method:   http.MethodGet,
				Endpoint: "/items",
				Key:      "",
				Value:    []byte{}},
			Want: "",
		},
	}

	kvs.HandleHTTPTests(t, tests)
}
