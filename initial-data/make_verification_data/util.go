package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func MkdirIfNotExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0777); err != nil {
			panic(err)
		}
	}
}

func getSnapshotFromRequest(serverName string, request Request) Snapshot {
	req, err := http.NewRequest(request.Method, fmt.Sprintf("%s%s", serverName, request.Resource), bytes.NewBuffer([]byte(request.Body)))
	if err != nil {
		fmt.Printf("Error: on making httpRequest [%v %v]", request.Method, request.Resource)
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.URL.RawQuery = request.Query
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	return Snapshot{
		Request: request,
		Response: Response{
			StatusCode: res.StatusCode,
			Body:       string(bytes),
		},
	}
}
