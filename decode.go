package elasticsearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/b3ntly/elasticsearch/mock"
	"github.com/b3ntly/insertjson"
)

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

func indexResponseToDocument(HTTPResponseBody []byte) (string, error) {
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

func getDocumentResponseToDocument(HTTPResponseBody []byte) ([]byte, error) {
	//fmt.Println(string(HTTPResponseBody))
	response := &mock.Generic{}
	err := json.Unmarshal(HTTPResponseBody, response)

	if err != nil {
		return nil, err
	}

	if response.Found == false {
		return nil, errors.New(fmt.Sprintf("Failed to get document with id: %v", response.ID))
	}

	return insertjson.Property("_id", response.ID, response.Source), err
}

func searchResponseToDocument(HTTPResponseBody []byte) ([][]byte, error) {
	response := &mock.Generic{}
	err := json.Unmarshal(HTTPResponseBody, response)

	if err != nil {
		return nil, err
	}

	documents := make([][]byte, len(response.Hits.Hits))
	for i, val := range response.Hits.Hits {
		documents[i] = insertjson.Property("_id", val.ID, val.Source)
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

func bulkUpdateResponseToIDs(HTTPResponseBody []byte) ([]string, error) {
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

func bulkDeleteResponseToIDs(HTTPResponseBody []byte) ([]string, error) {
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
