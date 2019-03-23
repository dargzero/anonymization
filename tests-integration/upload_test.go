//+build integration

package tests_integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dargzero/anonymization/swagger"
	"io"
	"testing"
)

func TestApi_UploadSession(t *testing.T) {

	upload := "/upload"

	t.Run("create upload session", func(t *testing.T) {
		dataset := "create-upload-session"
		setup(dataset)
		defer teardown(dataset)
		status, _ := send("POST", upload, sessionCfg(dataset))
		if status != 200 {
			t.Errorf("unexpected status: %v", status)
		}
	})

	t.Run("existing upload session", func(t *testing.T) {
		dataset := "existing-session-test"
		setup(dataset)
		defer teardown(dataset)
		send("POST", upload, sessionCfg(dataset))
		status, _ := send("POST", upload, sessionCfg(dataset))
		if status != 400 {
			t.Errorf("unexpected status: %v", status)
		}
	})
}

func TestApi_UploadData(t *testing.T) {

	upload := "/upload"

	t.Run("upload data to invalid session", func(t *testing.T) {
		dataset := "upload-invalid-session"
		setup(dataset)
		defer teardown(dataset)
		status, _ := sendResource("POST", upload+"/invalid", "u_data_medical.json")
		if status != 400 {
			t.Errorf("unexpected status: %v", status)
		}
	})

	t.Run("upload data", func(t *testing.T) {
		dataset := "upload-data"
		sessionPath := setupSession(dataset)
		defer teardown(dataset)
		status, body := sendResource("POST", sessionPath, "u_data_medical.json")
		var result swagger.UploadResponse
		json.Unmarshal([]byte(body), &result)
		if status != 200 || !result.InsertSuccessful || !result.FinalizeSuccessful {
			t.Errorf("failed to upload: status=%v, result=%v", status, result)
		}
	})

}

func setup(dataset string) {
	path := "/datasets/" + dataset
	call("DELETE", path)
	sendResource("PUT", path, "ds_medical.json")
}

func setupSession(dataset string) string {
	setup(dataset)
	_, res := send("POST", "/upload", sessionCfg(dataset))
	var session swagger.CreateUploadSessionResponse
	json.Unmarshal([]byte(res), &session)
	return fmt.Sprintf("/upload/%s?last=true", session.SessionID)
}

func teardown(dataset string) {
	path := "/datasets/" + dataset
	call("DELETE", path)
}

func sessionCfg(dataset string) io.Reader {
	req := swagger.CreateUploadSessionRequest{
		DatasetName: dataset,
	}
	b, _ := json.Marshal(req)
	return bytes.NewReader(b)
}
