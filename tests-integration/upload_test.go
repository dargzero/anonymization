//+build integration

package tests_integration

import (
	"encoding/json"
	"github.com/dargzero/anonymization/anonmodel"
	"github.com/dargzero/anonymization/swagger"
	"testing"
)

var dataset = "/datasets/upload-test"
var upload = "/upload"

func TestApi_UploadSession(t *testing.T) {

	t.Run("create upload session", func(t *testing.T) {
		setup()
		defer teardown()
		status, _ := send("POST", upload, "u_session.json")
		if status != 200 {
			t.Errorf("unexpected status: %v", status)
		}
	})

	t.Run("existing upload session", func(t *testing.T) {
		setup()
		defer teardown()
		send("POST", upload, "u_session.json")
		status, _ := send("POST", upload, "u_session.json")
		if status != 400 {
			t.Errorf("unexpected status: %v", status)
		}
	})

	t.Run("delete upload session", func(t *testing.T) {
		setup()
		defer teardown()
		_, sessRes := send("POST", upload, "u_session.json")
		var session anonmodel.UploadSessionData
		json.Unmarshal([]byte(sessRes), &session)
		sessionPath := upload + "/" + session.SessionID
		status, _ := call("DELETE", sessionPath)
		if status != 204 {
			t.Errorf("unexpected status: %v", status)
		}
	})

	t.Run("delete non-existing upload session", func(t *testing.T) {
		setup()
		defer teardown()
		sessionPath := upload + "/non-existing-session"
		status, _ := call("DELETE", sessionPath)
		if status != 404 {
			t.Errorf("unexpected status: %v", status)
		}
	})
}

func TestApi_UploadData(t *testing.T) {

	setup()
	defer teardown()

	t.Run("upload data to invalid session", func(t *testing.T) {
		status, _ := send("POST", upload+"/invalid", "u_data_machine.json")
		if status != 400 {
			t.Errorf("unexpected status: %v", status)
		}
	})

	t.Run("upload medical data", func(t *testing.T) {
		_, sessRes := send("POST", upload, "u_session.json")
		var session anonmodel.UploadSessionData
		json.Unmarshal([]byte(sessRes), &session)
		sessionPath := upload + "/" + session.SessionID + "?last=true"
		status, body := send("POST", sessionPath, "u_data_medical.json")
		var result swagger.UploadResponse
		json.Unmarshal([]byte(body), &result)
		if status != 200 || !result.InsertSuccessful || !result.FinalizeSuccessful {
			t.Errorf("failed to upload: status=%v, result=%v", status, result)
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
