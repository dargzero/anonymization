//+build integration

package tests_integration

import (
	"fmt"
	"github.com/dargzero/anonymization/anonmodel"
	"github.com/dargzero/anonymization/swagger"
	"testing"
	"time"
)

func TestApi_Data(t *testing.T) {

	t.Run("download dataset", func(t *testing.T) {
		dataset := "test-download-dataset"
		uploadTestData(dataset)
		defer teardown(dataset)
		path := "/data/" + dataset
		status, body := waitUntilDataAppears(path)
		if status != 200 {
			t.Errorf("unexpected status: %v, %v", status, body)
		}
	})

	t.Run("download anonymized dataset", func(t *testing.T) {
		dataset := "test-anon-dataset"
		uploadTestData(dataset)
		defer teardown(dataset)
		status, body := waitUntilDataAppears("/anon/" + dataset)
		if status != 200 {
			t.Errorf("unexpected status: %v, %v", status, body)
		}
		t.Logf("%v", body)
	})

	t.Run("download document", func(t *testing.T) {

		tests := []struct {
			dataset string
			path    string
		}{
			{"test-download-single-document", "/data/"},
			{"test-download-anonymized-document", "/anon/"},
		}

		for _, test := range tests {
			t.Run(test.dataset, func(t *testing.T) {
				uploadTestData(test.dataset)
				defer teardown(test.dataset)
				path := test.path + test.dataset
				status, body := waitUntilDataAppears(path)
				if status != 200 {
					t.Errorf("unexpected status: %v, %v", status, body)
				}
				var result swagger.ListDataResponse
				mustUnmarshal([]byte(body), &result)
				id := result.Result[0]["_id"]
				status, body = call("GET", fmt.Sprintf("%s/%s", path, id))
				if status != 200 {
					t.Errorf("unexpected status: %v, %v", status, body)
				}
				var doc anonmodel.Document
				mustUnmarshal([]byte(body), &doc)
				if doc["_id"] != id {
					t.Errorf("expected: %v, got %v", id, doc["_id"])
				}
			})
		}
	})
}

func uploadTestData(dataset string) {
	sessionPath := setupSession(dataset, "ds_machine.json")
	status, body := sendResource("POST", sessionPath, "u_data_machine.json")
	if status != 200 {
		panic(fmt.Sprintf("upload failed: %v, %v", status, body))
	}
}

func waitUntilDataAppears(path string) (int, string) {
	var status int
	var body string
	wait(func() bool {
		status, body = call("GET", path)
		return status != 404
	}, 8000)
	return status, body
}

func wait(predicate func() bool, timeout int) {
	for !predicate() && timeout > 0 {
		time.Sleep(100 * time.Millisecond)
		timeout = timeout - 100
	}
}
