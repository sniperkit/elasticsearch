package elasticsearch

import (
	"bytes"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/pkg/errors"
	"strings"
)

type (
	// Contains methods for querying the elasticsearch REST API.
	REST struct {
		HTTPClient *http.Client
		BaseURL string
	}

	// Response structure for the elasticsearch Index API
	IndexResponse struct {
		Shards struct {
			Total int `json:"total"`
			Successful int `json:"successful"`
			Failed int `json:"failed"`
		} `json:"_shards"`
		Index string `json:"_index"`
		Type string `json:"_type"`
		ID string `json:"_id"`
		Version int64 `json:"_version"`
		Created bool `json:"created"`
		Result string `json:"result"`
	}

	// Response structure for the elasticsearch Document-get API
	GetDocumentResponse struct {
		Shards struct {
			Total int `json:"total"`
			Successful int `json:"successful"`
			Failed int `json:"failed"`
		} `json:"_shards"`
		Index string `json:"_index"`
		Type string `json:"_type"`
		ID string `json:"_id"`
		Version int64 `json:"_version"`
		Found bool `json:"found"`
		Document json.RawMessage `json:"_source"`
	}

	// Response structure for the elasticsearch Search API
	SearchResponse struct {
		TimedOut bool `json:"timed_out"`
		Took int `json:"took"`
		Shards struct {
			Total int `json:"total"`
			Successful int `json:"successful"`
			Failed int `json:"failed"`
		} `json:"_shards"`
		Hits struct {
			Total int `json:"total"`
			MaxScore float64 `json:"max_score"`
			Hits []struct {
				Index string `json:"_index"`
				Type string `json:"_type"`
				ID string `json:"_id"`
				Score float64 `json:"_score"`
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	// Response structure for the elasticsearch Index-delete API
	DeleteIndexResponse struct {
		Shards struct {
			Total int `json:"total"`
			Successful int `json:"successful"`
			Failed int `json:"failed"`
		} `json:"_shards"`
		Acknowledged bool `json:"acknowledged"`
	}

	// Response structure for the elasticsearch Document-delete
	DeleteDocumentResponse struct {
		Shards struct {
			Total int `json:"total"`
			Failed int `json:"failed"`
			Successful int `json:"successful"`
		} `json:"_shards"`
		Found bool `json:"found"`
		Index string `json:"_index"`
		Type string `json:"_type"`
		ID string `json:"_id"`
		Version int `json:"_version"`
		Result string `json:"result"`
	}

	// Response structure for the elasticsearch Document-update API
	UpdateDocumentResponse struct {
		Shards struct {
			Total int `json:"total"`
			Successful int `json:"successful"`
			Failed int `json:"failed"`
		} `json:"_shards"`
		Index string `json:"_index"`
		Type string `json:"_type"`
		ID string `json:"_id"`
		Version int `json:"_version"`
		Result string `json:"result"`
		Created bool `json:"created"`
	}
)

// Call the elasticsearch Search API for  given index
func (r *REST) SearchIndex(index string, queryString string) ([]*Document, error){
	qs := fmt.Sprintf("_search?q=%v", queryString)
	URL := r.BuildURL(index, qs)
	body, err := r.Request("GET", URL, nil)

	if err != nil {
		return nil, err
	}

	return SearchResponseToDocument(body)
}

// Call the elasticsearch Index API
func (r *REST) DeleteIndex(index string) error {
	URL := r.BuildURL(index)
	body, err := r.Request("DELETE", URL, nil)

	if err != nil {
		return err
	}

	return DeleteIndexResponseToDocument(body)
}

// Call the elasticsearch Search API for  given index
func (r *REST) SearchType(index string, _type string, queryString string) ([]*Document, error){
	qs := fmt.Sprintf("_search?q=%v", queryString)
	URL := r.BuildURL(index, _type, qs)
	body, err := r.Request("GET", URL, nil)

	if err != nil {
		return nil, err
	}

	return SearchResponseToDocument(body)
}

// Call the elasticsearch Index API
func (r *REST) InsertDocument(index string, _type string, doc []byte) (*Document, error){
	URL := r.BuildURL(index, _type, "?refresh")
	body, err := r.Request("POST", URL, doc)

	if err != nil {
		return nil, err
	}

	return IndexResponseToDocument(body)
}

// Call the elasticsearch Document API
func (r *REST) GetDocument(index string, _type string, ID string) (*Document, error){
	URL := r.BuildURL(index, _type, ID)
	body, err := r.Request("GET", URL, nil)

	if err != nil {
		return nil, err
	}

	return GetDocumentResponseToDocument(body)
}

// Call the elasticsearch Document API
func (r *REST) UpdateDocument(index string, _type string, ID string, doc []byte) error {
	URL := r.BuildURL(index, _type, ID, "?refresh")
	body, err := r.Request("PUT", URL, doc)

	if err != nil {
		return err
	}

	return UpdateDocumentResponseToDocument(body)
}

// Call the elasticsearch Document API
func (r *REST) DeleteDocument(index string, _type string, ID string) error {
	URL := r.BuildURL(index, _type, ID, "?refresh")
	body, err := r.Request("GET", URL, nil)

	if err != nil {
		return err
	}

	return DeleteDocumentResponseToDocument(body)
}

// Concatenate a URL from an array of strings.
func (r *REST) BuildURL(parts ...string) string {
	parts = append([]string{r.BaseURL },  parts...)
	return strings.Join(parts, "/")
}

func (r *REST) BuildRequest(method string, url string, body []byte) (*http.Request, error){
	var req *http.Request
	var err error

	if body == nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer([]byte{}))
	} else {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	}

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func (r *REST) SendRequest(req *http.Request) ([]byte, error){
	response, err := r.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	// fmt.Println(response.StatusCode)

	if response.StatusCode >= 299 {
		return nil, errors.New(fmt.Sprintf("Invalid status code during InsertDocument: %v", response.StatusCode))
	}

	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

// Generic method to make a JSON request against a configured endpoint.
func (r *REST) Request(method string, url string, body []byte) ([]byte, error){
	req, err := r.BuildRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	return r.SendRequest(req)
}