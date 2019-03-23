package tests_integration

import (
	"testing"
)

var dataset = "/datasets/upload-test"

func TestApi_UploadSession(t *testing.T) {

	uploadPath := "/upload"

	t.Run("create upload session", func(t *testing.T) {
		setup()
		defer teardown()
		status, _ := send("POST", uploadPath, "u_session.json")
		if status != 200 {
			t.Errorf("unexpected status: %v", status)
		}
	})

	t.Run("existing upload session", func(t *testing.T) {
		setup()
		defer teardown()
		send("POST", uploadPath, "u_session.json")
		status, _ := send("POST", uploadPath, "u_session.json")
		if status != 400 {
			t.Errorf("unexpected status: %v", status)
		}
	})
}

func setup() {
	call("DELETE", dataset)
	send("PUT", dataset, "ds_medical.json")
}

func teardown() {
	call("DELETE", dataset)
}
