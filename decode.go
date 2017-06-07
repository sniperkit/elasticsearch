package elasticsearch

import (
	"github.com/b3ntly/elasticsearch/mock"
	"encoding/json"
	"errors"
	"fmt"
)

// Reference to an elasticsearch document. To simplify deserialization between
// the elasticsearch rest response and the struct type response of this library,
// we do not decode the document but instead leave it serialized as a *json.RawMessage.
//
// Elasticsearch also separates document IDs from document bodies, hence the separate struct fields.
type Document struct {
	ID string
	Body json.RawMessage
}

func errorResponseToError(HTTPResponseBody []byte) error {
	response := &mock.Generic{}
	err := json.Unmarshal(HTTPResponseBody, response)

	if err != nil {
		return err
	}

	reason := ""
	for _, r := range response.Error.RootCause {
		reason += "," + r.Reason
	}

	return errors.New(reason)
}

func indexResponseToDocument(HTTPResponseBody []byte) (string, error){
	response := &mock.Generic{}
	err := json.Unmarshal(HTTPResponseBody, response)

	if err != nil {
		return "", err
	}

	if response.Created != true {
		return "", errors.New("Failed to create document.")
	}

	return response.ID, err
}



func deleteIndexResponseToDocument(HTTPResponseBody []byte) error {
	response := &mock.Generic{}

	err := json.Unmarshal(HTTPResponseBody, response)
	if err != nil {
		return err
	}

	if response.Acknowledged != true {
		return errors.New("Failed to drop index.")
	}

	return nil
}

func getDocumentResponseToDocument(HTTPResponseBody []byte) (*Document, error){
	//fmt.Println(string(HTTPResponseBody))
	response := &mock.Generic{}
	err := json.Unmarshal(HTTPResponseBody, response)

	if err != nil {
		return nil, err
	}

	if response.Found == false {
		return nil, errors.New(fmt.Sprintf("Failed to get document with id: %v", response.ID))
	}

	return &Document{ ID: response.ID, Body: response.Source }, err
}

func searchResponseToDocument(HTTPResponseBody []byte) ([]*Document, error){
	response := &mock.Generic{}
	err := json.Unmarshal(HTTPResponseBody, response)

	if err != nil {
		return nil, err
	}

	documents := make([]*Document, len(response.Hits.Hits))
	for i, val := range response.Hits.Hits {
		documents[i] = &Document{
			ID: val.ID,
			Body: val.Source,
		}
	}

	return documents, err
}

func deleteDocumentResponseToDocument(HTTPResponseBody []byte) error {
	response := &mock.Generic{}
	err := json.Unmarshal(HTTPResponseBody, response)

	if err != nil {
		return err
	}

	if response.Found != true {
		return errors.New("Document was not found.")
	}

	return nil
}

func updateDocumentResponseToDocument(HTTPResponseBody []byte) error {
	response := &mock.Generic{}
	err := json.Unmarshal(HTTPResponseBody, response)

	if err != nil {
		return errors.New("Failed to unmarshal response")
	}

	if response.Created == true {
		return errors.New("Accidentally upserted document...")
	}

	return nil
}

func bulkInsertResponseToIDs(HTTPResponseBody []byte) ([]string, error) {
	response := &mock.Generic{}
	err := json.Unmarshal(HTTPResponseBody, response)

	if err != nil {
		return nil, err
	}

	inserted := make([]string, len(response.Items))
	for idx, item := range response.Items {
		if item.Index.Created == false {
			err = errors.New("Some documents were not inserted.")
		}
		inserted[idx] = item.Index.ID
	}

	return inserted, err
}

func bulkUpdateResponseToIDs(HTTPResponseBody []byte)([]string, error){
	response := &mock.Generic{}
	err := json.Unmarshal(HTTPResponseBody, response)

	if err != nil {
		return nil, err
	}

	updated := make([]string, len(response.Items))

	for idx, item := range response.Items {
		updated[idx] = item.Update.ID
	}

	return updated, err
}

func bulkDeleteResponseToIDs(HTTPResponseBody []byte)([]string, error){
	//fmt.Println(string(HTTPResponseBody))
	response := &mock.Generic{}
	err := json.Unmarshal(HTTPResponseBody, response)

	if err != nil {
		return nil, err
	}

	deleted := make([]string, len(response.Items))

	for idx, item := range response.Items {
		if item.Delete.Found == false {
			err = errors.New("Some documents were not found and thus not deleted.")
		}
		deleted[idx] = item.Delete.ID
	}

	return deleted, err
}



