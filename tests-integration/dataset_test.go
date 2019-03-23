//+build integration

package tests_integration

import (
	"github.com/dargzero/anonymization/anonmodel"
	"testing"
)

func TestApi_InvalidDataSet(t *testing.T) {

	path := "/datasets/crate-bad"
	payload := "ds_invalid.json"
	status, _ := sendResource("PUT", path, payload)
	if status != 400 {
		t.Errorf("unexpected status: %v", status)
	}
}

func TestApi_DataSets(t *testing.T) {

	payload := "ds_medical.json"

	t.Run("delete existing dataset", func(t *testing.T) {
		path := "/datasets/delete-existing-test"
		sendResource("PUT", path, payload)
		status, _ := call("DELETE", path)
		if status != 204 {
			t.Errorf("delete: unexpected status: %v", status)
		}
	})

	t.Run("delete non existing dataset", func(t *testing.T) {
		path := "/datasets/delete-non-existing-test"
		call("DELETE", path)
		status, _ := call("DELETE", path)
		if status != 404 {
			t.Errorf("delete: unexpected status: %v", status)
		}
	})

	t.Run("create dataset", func(t *testing.T) {
		path := "/datasets/create-dataset"
		call("DELETE", path)
		status, _ := sendResource("PUT", path, payload)
		if status != 201 {
			t.Errorf("create: unexpected status: %v", status)
		}
	})

	t.Run("get dataset metadata", func(t *testing.T) {
		path := "/datasets/get-dataset-metadata"
		call("DELETE", path)
		sendResource("PUT", path, payload)
		_, body := call("GET", path)
		var actual anonmodel.Dataset
		mustUnmarshal([]byte(body), &actual)
		if actual.Name != "get-dataset-metadata" {
			t.Errorf("invalid dataset metadata: %v", body)
		}
	})

	t.Run("list datasets", func(t *testing.T) {
		path1 := "/datasets/new-dataset1"
		path2 := "/datasets/new-dataset2"
		call("DELETE", path1)
		call("DELETE", path2)
		sendResource("PUT", path1, payload)
		sendResource("PUT", path2, payload)
		_, body := call("GET", "/datasets")
		var actual []anonmodel.Dataset
		mustUnmarshal([]byte(body), &actual)

		assertContains := func(coll []anonmodel.Dataset, name string) {
			for _, ds := range coll {
				if ds.Name == name {
					return
				}
			}
			t.Errorf("dataset not found: %v", name)
		}

		assertContains(actual, "new-dataset1")
		assertContains(actual, "new-dataset2")
	})

}
