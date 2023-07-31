package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Expression struct {
}

type MockRequestResponse struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

type MockRequest struct {
	HttpRequestUrl string              `json:"httpRequestUrl"`
	Method         string              `json:"method"`
	Response       MockRequestResponse `json:"response"`
	Expressions    []Expression        `json:"expressions"`
}

// ReadMapping /*
func ReadMapping() []*MockRequest {
	var requests []*MockRequest
	root := os.Getenv("cloud_mock_root")
	if root == "" {
		root, _ = os.Getwd()
	}
	requestFolder := filepath.Join(root, "requests")
	_ = filepath.WalkDir(requestFolder, func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(d.Name(), ".json") && !d.IsDir() {
			var mockRequests []*MockRequest
			bytes, _ := os.ReadFile(path)
			_ = json.Unmarshal(bytes, &mockRequests)
			for _, mockRequest := range mockRequests {
				requests = append(requests, mockRequest)
			}
		}
		return nil
	})
	return requests
}

func SaveRequest(request *Request) {
	root := os.Getenv("cloud_mock_root")
	if root == "" {
		root, _ = os.Getwd()
	}
	requestFolder := filepath.Join(root, "requests")
	mockRequest := &MockRequest{
		HttpRequestUrl: fmt.Sprintf("%s%s", request.Host, request.Path),
		Method:         request.Method,
		Expressions:    []Expression{},
		Response: MockRequestResponse{
			Type:  "raw",
			Value: []byte("{}"),
		},
	}
	mockRequestBytes, _ := json.Marshal([]*MockRequest{mockRequest})
	dir := filepath.Join(requestFolder, "tmp", url.QueryEscape(request.Host))
	err := os.MkdirAll(dir, os.ModeDir|0777)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(dir, url.QueryEscape(request.Path)+".json"), mockRequestBytes, 0666)
	if err != nil {
		log.Fatal(err)
	}
}

type Request struct {
	Host   string
	Path   string
	Method string
}

func ExpressRule(mockRequests []*MockRequest, r *Request) string {
	hostPath := fmt.Sprintf("%s%s", r.Host, r.Path)
	var match bool
	for _, mockRequest := range mockRequests {
		if strings.HasPrefix(mockRequest.HttpRequestUrl, "~") {
			match, _ = regexp.MatchString(mockRequest.HttpRequestUrl[1:], hostPath)
		} else {
			match = mockRequest.HttpRequestUrl == hostPath
		}
		if match && r.Method == mockRequest.Method {
			log.Println(hostPath)
			switch mockRequest.Response.Type {
			case "raw":
				{
					return string(mockRequest.Response.Value)
				}
			default:
				return "{}"
			}
		}
	}
	// generate the blank requests if not match
	SaveRequest(r)
	return "{}"
}
