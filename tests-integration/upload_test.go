//+build integration

package tests_integration

import (
	"github.com/dargzero/anonymization/swagger"
	"testing"
)

func TestApi_UploadSession(t *testing.T) {

	datasetType := "ds_medical.json"
	upload := "/upload"

	t.Run("create upload session", func(t *testing.T) {
		dataset := "create-upload-session"
		setup(dataset, datasetType)
		defer teardown(dataset)
		status, _ := send("POST", upload, sessionCfg(dataset))
		if status != 200 {
			t.Errorf("unexpected status: %v", status)
		}
	})

	t.Run("existing upload session", func(t *testing.T) {
		dataset := "existing-session-test"
		setup(dataset, datasetType)
		defer teardown(dataset)
		send("POST", upload, sessionCfg(dataset))
		status, _ := send("POST", upload, sessionCfg(dataset))
		if status != 400 {
			t.Errorf("unexpected status: %v", status)
		}
	})
}

func TestApi_UploadData(t *testing.T) {

	datasetType := "ds_medical.json"
	upload := "/upload"

	t.Run("upload data to invalid session", func(t *testing.T) {
		dataset := "upload-invalid-session"
		setup(dataset, datasetType)
		defer teardown(dataset)
		status, _ := sendResource("POST", upload+"/invalid", "u_data_medical.json")
		if status != 400 {
			t.Errorf("unexpected status: %v", status)
		}
	})

	t.Run("upload data", func(t *testing.T) {
		dataset := "upload-data"
		sessionPath := setupSession(dataset, datasetType)
		defer teardown(dataset)
		status, body := sendResource("POST", sessionPath, "u_data_medical.json")
		var result swagger.UploadResponse
		mustUnmarshal([]byte(body), &result)
		if status != 200 || !result.InsertSuccessful || !result.FinalizeSuccessful {
			t.Errorf("failed to upload: status=%v, result=%v", status, result)
		}
	})
}
