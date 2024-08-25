package keyvaluestore

import (
	"testing"

	kvs "github.com/defoeam/kvs/kv"
)

// Before any tests can be ran on the http server,
// we must actually start it.
func init() {
	// Start the server on a seperate thread.
	go kvs.StartServer()
}

func TestSet1(t *testing.T) {
	tests := []HttpTest{
		{
			name: "Populate the store, pt. 1",
			args: HttpArgs{method: "POST", endpoint: "items", key: "0", value: "Hello World!"},
			want: `{"key":"0","value":"Hello World!"}`,
		},
		{
			name: "Populate the store, pt. 2",
			args: HttpArgs{method: "POST", endpoint: "items", key: "1", value: "[53,58,203]"},
			want: `{"key":"1","value":"[53,58,203]"}`,
		},
	}

	HandleTests(t, tests)
}
