//+build integration

package tests_integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dargzero/anonymization/swagger"
	"io"
	"io/ioutil"
	"net/http"
)

const baseUrl = "http://localhost:9137/v1"

func call(method, operation string) (int, string) {
	req := request(method, operation, nil)
	res := do(req)
	return read(res)
}

func send(method, apiPath string, payload io.Reader) (int, string) {
	req := request(method, apiPath, payload)
	res := do(req)
	return read(res)
}

func sendResource(method, apiPath, payloadName string) (int, string) {
	return send(method, apiPath, resource(payloadName))
}

func read(res *http.Response) (int, string) {
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return res.StatusCode, string(body)
}

func do(r *http.Request) *http.Response {
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	return res
}

func request(method, operation string, reader io.Reader) *http.Request {
	req, err := http.NewRequest(method, path(operation), reader)
	if err != nil {
		panic("failed to create request: " + operation)
	}
	return req
}

func resource(name string) io.Reader {
	b, err := ioutil.ReadFile("resources/" + name)
	if err != nil {
		panic("test resource not found: " + name)
	}
	return bytes.NewBuffer(b)
}

func path(p string) string {
	return baseUrl + p
}

func setup(dataset, datasetType string) {
	path := "/datasets/" + dataset
	call("DELETE", path)
	sendResource("PUT", path, datasetType)
}

func setupSession(dataset, datasetType string) string {
	setup(dataset, datasetType)
	_, res := send("POST", "/upload", sessionCfg(dataset))
	var session swagger.CreateUploadSessionResponse
	mustUnmarshal([]byte(res), &session)
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

func mustUnmarshal(data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		panic(err)
	}
}
