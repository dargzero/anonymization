//+build integration

package tests_integration

import (
	"bytes"
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

func send(method, apiPath, payload string) (int, string) {
	req := request(method, apiPath, resource(payload))
	res := do(req)
	return read(res)
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
