/*
 * Data Anonymization Server
 *
 * This is a data anonymization server. You can set the anonymization requirements for the different datasets individually, and upload data to them. The uploaded data is anonymized on the server and can be then downloaded.
 *
 * API version: 0.1-alpha
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

import (
	"github.com/dargzero/anonymization/anonmodel"
)

// CreateDatasetRequest is the JSON object the client sends to create a dataset
type CreateDatasetRequest struct {
	Settings anonmodel.AnonymizationSettings    `json:"settings"`
	Fields   []anonmodel.FieldAnonymizationInfo `json:"fields"`
}

// DatasetResponse is the JSON object the server sends to describe a dataset
type DatasetResponse struct {
	Name     string                             `json:"name"`
	Settings anonmodel.AnonymizationSettings    `json:"settings"`
	Fields   []anonmodel.FieldAnonymizationInfo `json:"fields"`
}

func createDataset(datasetName string, request *CreateDatasetRequest) anonmodel.Dataset {
	return anonmodel.Dataset{
		Name:     datasetName,
		Settings: request.Settings,
		Fields:   request.Fields,
	}
}

func createDatasetResponse(dataset *anonmodel.Dataset) DatasetResponse {
	return DatasetResponse{
		Name:     dataset.Name,
		Settings: dataset.Settings,
		Fields:   dataset.Fields,
	}
}

// DatasetsResponse is the JSON object the server sends to list datasets
type DatasetsResponse []DatasetResponse

func createDatasetsResponse(array []anonmodel.Dataset) DatasetsResponse {
	datasets := make([]DatasetResponse, len(array))
	for ix, item := range array {
		datasets[ix] = createDatasetResponse(&item)
	}

	return DatasetsResponse(datasets)
}
