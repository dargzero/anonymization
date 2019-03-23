//+build integration

package tests_integration

import "testing"

func TestApi_Ping(t *testing.T) {
	status, _ := call("GET", "/ping")
	if status != 200 {
		t.Errorf("ping: unexpected status: %v", status)
	}
}
