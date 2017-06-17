package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/b3ntly/elasticsearch/mock"
	"github.com/cch123/elasticsql"
	"io/ioutil"
	"net/http"
)

// rest interface with elasticsearch
type rest struct {
	HTTPClient *http.Client
	BaseURI    string
}

// Call the elasticsearch Search API for  given index
func (r *rest) searchIndex(index string, queryString string) ([]*Document, error) {
	URL, err := buildURI(r.BaseURI, map[string]string{"index": index, "suffix": "_search"}, map[string]string{"q": queryString})

	if err != nil {
		return nil, err
	}

	body, err := r.request("GET", URL, nil)

	if err != nil {
		return nil, err
	}

	return searchResponseToDocument(body)
}

// Call the elasticsearch Index API
func (r *rest) deleteIndex(index string) error {
	URL, err := buildURI(r.BaseURI, map[string]string{"index": index}, nil)

	if err != nil {
		return err
	}

	body, err := r.request("DELETE", URL, nil)

	if err != nil {
		return err
	}

	return deleteIndexResponseToDocument(body)
}

func (r *rest) searchSQL(index string, _type string, sql string) ([]*Document, error) {
	URL, err := buildURI(r.BaseURI, map[string]string{"index": index, "type": index, "suffix": "_search"}, nil)

	if err != nil {
		return nil, err
	}

	query, _, err := elasticsql.Convert(sql)

	if err != nil {
		return nil, err
	}

	body, err := r.request("GET", URL, []byte(query))

	if err != nil {
		return nil, err
	}

	return searchResponseToDocument(body)
}

// Call the elasticsearch Search API for  given index
func (r *rest) searchType(index string, _type string, queryString string) ([]*Document, error) {
	URL, err := buildURI(r.BaseURI, map[string]string{"index": index, "type": _type, "suffix": "_search"}, map[string]string{"q": queryString})

	if err != nil {
		return nil, err
	}

	body, err := r.request("GET", URL, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return searchResponseToDocument(body)
}

// Call the elasticsearch Index API
func (r *rest) insertDocument(index string, _type string, doc []byte) (string, error) {
	URL, err := buildURI(r.BaseURI, map[string]string{"index": index, "type": _type}, map[string]string{"refresh": "true"})

	if err != nil {
		return "", err
	}

	body, err := r.request("POST", URL, doc)

	if err != nil {
		return "", err
	}

	return indexResponseToDocument(body)
}

// Call the elasticsearch Bulk API with insert operations
func (r *rest) bulkInsertDocuments(index string, _type string, docs [][]byte) ([]string, error) {
	// construct an NDJSON payload that satisfies the Elasticsearch API
	payload := make([][]byte, len(docs)*2)

	// insert a bulk operation prefix before each document in the docs slice
	for i := 0; i < len(docs); i++ {
		operation, err := json.Marshal(&mock.BulkIndex{Index: &mock.Base{Index: index, Type: _type}})

		if err != nil {
			return nil, err
		}

		payload[i*2] = operation
		payload[(i*2)+1] = docs[i]
	}

	URL, err := buildURI(r.BaseURI, map[string]string{"suffix": "_bulk"}, map[string]string{"refresh": "true"})

	if err != nil {
		return nil, err
	}

	body, err := r.bulkRequest("POST", URL, payload)

	if err != nil {
		return nil, err
	}

	return bulkInsertResponseToIDs(body)
}

// Call the elasticsearch Document API
func (r *rest) getDocument(index string, _type string, ID string) (*Document, error) {
	URL, err := buildURI(r.BaseURI, map[string]string{"index": index, "type": _type, "suffix": ID}, nil)

	if err != nil {
		return nil, err
	}

	body, err := r.request("GET", URL, nil)

	if err != nil {
		return nil, err
	}

	return getDocumentResponseToDocument(body)
}

// Call the elasticsearch Document API
func (r *rest) updateDocument(index string, _type string, ID string, doc []byte) error {
	URL, err := buildURI(r.BaseURI, map[string]string{"index": index, "type": _type, "suffix": ID}, map[string]string{"refresh": "true"})

	if err != nil {
		return err
	}

	body, err := r.request("PUT", URL, doc)

	if err != nil {
		return err
	}

	return updateDocumentResponseToDocument(body)
}

// Call the elasticsearch Bulk API with update operations
func (r *rest) bulkUpdateDocuments(index string, _type string, docs []*Document) ([]string, error) {
	// construct an NDJSON payload that satisfies the Elasticsearch API
	payload := make([][]byte, len(docs)*2)

	// insert a bulk operation prefix before each document in the docs slice
	for i := 0; i < len(docs); i++ {
		operation, err := json.Marshal(&mock.BulkUpdate{Update: &mock.Resource{Index: index, Type: _type, ID: docs[i].ID}})

		if err != nil {
			return nil, err
		}

		payload[i*2] = operation

		doc, err := json.Marshal(&mock.BulkUpdatePayload{Doc: docs[i].Body})

		if err != nil {
			return nil, err
		}

		payload[(i*2)+1] = doc
	}

	URL, err := buildURI(r.BaseURI, map[string]string{"suffix": "_bulk"}, nil)

	if err != nil {
		return nil, err
	}

	body, err := r.bulkRequest("POST", URL, payload)

	if err != nil {
		return nil, err
	}

	return bulkUpdateResponseToIDs(body)
}

// Call the elasticsearch Document API
func (r *rest) deleteDocument(index string, _type string, ID string) error {
	URL, err := buildURI(r.BaseURI, map[string]string{"index": index, "type": _type, "suffix": ID}, map[string]string{"refresh": "true"})

	if err != nil {
		return err
	}

	body, err := r.request("DELETE", URL, nil)

	if err != nil {
		return err
	}

	return deleteDocumentResponseToDocument(body)
}

// Call the elasticsearch Bulk API with delete operations
func (r *rest) bulkDeleteDocuments(index string, _type string, IDs []string) ([]string, error) {
	// construct an NDJSON payload that satisfies the Elasticsearch bulk API delete operation
	payload := make([][]byte, len(IDs))

	// insert a bulk operation prefix before each document in the docs slice
	for idx, ID := range IDs {
		operation, err := json.Marshal(&mock.BulkDelete{Delete: &mock.Resource{Index: index, Type: _type, ID: ID}})

		if err != nil {
			return nil, err
		}

		payload[idx] = operation
	}

	URL, err := buildURI(r.BaseURI, map[string]string{"suffix": "_bulk"}, map[string]string{"refresh": "true"})

	if err != nil {
		return nil, err
	}

	body, err := r.bulkRequest("POST", URL, payload)

	if err != nil {
		return nil, err
	}

	return bulkDeleteResponseToIDs(body)
}

func (r *rest) buildRequest(method string, url string, body []byte) (*http.Request, error) {
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

func (r *rest) buildBulkRequest(method string, url string, bodies [][]byte) (*http.Request, error) {
	buffer := new(bytes.Buffer)

	for _, body := range bodies {
		buffer.Write(body)
		buffer.Write([]byte("\n"))
	}

	req, err := http.NewRequest(method, url, buffer)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func (r *rest) sendRequest(req *http.Request) ([]byte, error) {
	response, err := r.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 299 {
		return nil, errorResponseToError(contents)
	}

	return contents, err
}

// Generic method to make a JSON request against a configured endpoint.
func (r *rest) request(method string, url string, body []byte) ([]byte, error) {
	req, err := r.buildRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	return r.sendRequest(req)
}

// Generic method to make an NDJSON request against a configured endpoint
func (r *rest) bulkRequest(method string, url string, bodies [][]byte) ([]byte, error) {
	req, err := r.buildBulkRequest(method, url, bodies)

	if err != nil {
		return nil, err
	}

	return r.sendRequest(req)
}
