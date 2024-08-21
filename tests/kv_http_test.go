package keyvaluestore

import (
	"net/http"
	"testing"
)

func ExampleTest(t *testing.T) {
	res, err := http.Get("http://localhost:8080/items")

	if err != nil {
		t.Fatal()
	}

	println(res)
}
