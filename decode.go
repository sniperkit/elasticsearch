package elasticsearch

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Reference to an elasticsearch document. To simplify deserialization between
// the elasticsearch REST response and the struct type response of this library,
// we do not decode the document but instead leave it serialized as a *json.RawMessage.
//
// Elasticsearch also separates document IDs from document bodies, hence the separate struct fields.
type Document struct {
	ID string
	Body json.RawMessage
}

func IndexResponseToDocument(HTTPResponseBody []byte) (*Document, error){
	response := &IndexResponse{}
	err := json.Unmarshal(HTTPResponseBody, response)

	if err != nil {
		return nil, err
	}

	if response.Created != true {
		return nil, errors.New("Failed to create document.")
	}

	return &Document{ ID: response.ID, Body: HTTPResponseBody }, err
}

func DeleteIndexResponseToDocument(HTTPResponseBody []byte) error {
	response := &DeleteIndexResponse{}

	err := json.Unmarshal(HTTPResponseBody, response)
	if err != nil {
		return err
	}

	if response.Acknowledged != true {
		return errors.New("Failed to drop index.")
	}

	return nil
}

func GetDocumentResponseToDocument(HTTPResponseBody []byte) (*Document, error){
	response := &GetDocumentResponse{}
	err := json.Unmarshal(HTTPResponseBody, response)

	if err != nil {
		return nil, err
	}

	if response.Found == false {
		return nil, errors.New(fmt.Sprintf("Failed to get document with id: %v", response.ID))
	}

	return &Document{ ID: response.ID, Body: response.Document }, err
}

func SearchResponseToDocument(HTTPResponseBody []byte) ([]*Document, error){
	//fmt.Println(string(HTTPResponseBody))
	response := &SearchResponse{}
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

func DeleteDocumentResponseToDocument(HTTPResponseBody []byte) error {
	response := &DeleteIndexResponse{}
	err := json.Unmarshal(HTTPResponseBody, response)

	if err != nil {
		return err
	}

	if response.Acknowledged != true {
		return errors.New("Delete Index was not acknowledged.")
	}

	return nil
}

func UpdateDocumentResponseToDocument(HTTPResponseBody []byte) error {
	response := &IndexResponse{}
	err := json.Unmarshal(HTTPResponseBody, response)

	if err != nil {
		return errors.New("Failed to unmarshal response")
	}

	if response.Created == true {
		return errors.New("Accidentally upserted document...")
	}

	return nil
}

