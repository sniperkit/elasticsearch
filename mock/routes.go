package mock

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

// accept a raw request body and return two slices containing
// individual bulk request operations and their payloads
// note the request body is NDJSON not regular JSON
func parseBulkRequest(body []byte) ([]*Operation, [][]byte, error) {
	// parse the first object of the body and determine if
	// this is a bulk delete operation (which must be handled separately)
	objects := bytes.Split(body, []byte("\n"))

	firstOperation := &Operation{}
	err := json.Unmarshal(objects[0], firstOperation)

	if err != nil {
		return nil, nil, err
	}

	// objects are formatted as a series of JSON objects in an alternating
	// pattern of operation, payload, operation, payload
	// except for delete operations which do not contain payloads
	// thus we must only process operations for pure delete requests
	var operations []*Operation
	var payloads [][]byte
	if firstOperation.Delete != nil {
		// the last object is an empty line, conforming with elasticsearch's implementation of NDJSON...
		operations = make([]*Operation, len(objects)-1)
		payloads = make([][]byte, 0)

		// process each object into an operation
		for idx, object := range objects {
			// handle the edge case of a blank line
			if len(object) == 0 {
				continue
			}

			op := &Operation{}
			err := json.Unmarshal(object, op)

			if err != nil {
				return nil, nil, err
			}

			operations[idx] = op
		}

		// for all other cases we split objects into a slice
		// operations containing all odd-indexed elements of objects
		// and payloads containing all even-indexed elements
	} else {
		operations = make([]*Operation, len(objects)/2)
		payloads = make([][]byte, len(objects)/2)

		for i := 0; i < len(objects)/2; i++ {
			op := new(*Operation)
			err := json.Unmarshal(objects[i*2], op)

			if err != nil {
				return nil, nil, err
			}

			operations[i] = *op
			payloads[i] = objects[i*2+1]
		}
	}

	return operations, payloads, err
}

func BulkAPI(w http.ResponseWriter, req *http.Request) {
	contents, err := ioutil.ReadAll(req.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	operations, payloads, err := parseBulkRequest(contents)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	for idx, operation := range operations {
		// handle bulk insert operations
		if operation.Index != nil {
			payload := payloads[idx]
			doc, err := database.insert(operation.Index.Index, operation.Index.Type, payload)

			if err != nil {
				// return without executing any operations
				// this does not match elasticsearch's process
				http.Error(w, "testing", http.StatusBadRequest)
				return
			}

			operation.Index.Created = true
			operation.Index.ID = doc.ID

			// handle bulk update operations
		} else if operation.Update != nil {
			payload := payloads[idx]
			updated, err := database.upsertDocument(operation.Update.Index, operation.Update.Type, operation.Update.ID, payload)

			if err != nil {
				// return without executing any operations
				// this does not match elasticsearch's process
				http.Error(w, "testing", http.StatusBadRequest)
				return
			}

			operation.Update.Created = !updated

		} else if operation.Delete != nil {
			deleted := database.deleteDocument(operation.Delete.Index, operation.Delete.Type, operation.Delete.ID)
			operation.Delete.Found = deleted
		}
	}

	response := &Generic{Items: operations}
	js, err := json.Marshal(response)

	if err != nil {
		// return without executing any operations
		// this does not match elasticsearch's process
		http.Error(w, "testing", http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func DeleteIndex(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	index := vars["index"]

	database.deleteIndex(index)

	// return proper response
	resp := Generic{
		Acknowledged: true,
	}

	js, err := json.Marshal(resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func SearchIndex(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	index := vars["index"]

	hits, err := database.searchIndex(index)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := Generic{}
	resp.Hits.Hits = hits
	resp.Hits.Total = len(hits)

	js, err := json.Marshal(resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func SearchType(w http.ResponseWriter, req *http.Request) {
	// todo: fix nomenclature
	vars := mux.Vars(req)
	index := vars["index"]
	_type := vars["_type"]

	hits, err := database.searchType(index, _type)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := Generic{}
	resp.Hits.Hits = hits
	resp.Hits.Total = len(resp.Hits.Hits)

	js, err := json.Marshal(resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func InsertDocument(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	index := vars["index"]
	_type := vars["_type"]

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	doc, err := database.insert(index, _type, body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// return proper response
	resp := Generic{
		Index:   index,
		Type:    _type,
		ID:      doc.ID,
		Created: true,
		Result:  "created",
	}

	js, err := json.Marshal(resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func GetDocumentByID(w http.ResponseWriter, req *http.Request) {
	// handle req
	vars := mux.Vars(req)
	index := vars["index"]
	_type := vars["_type"]
	ID := vars["id"]

	// if document exists return as GetDocumentResponse
	if doc := database.getDocument(index, _type, ID); doc != nil {
		body, err := json.Marshal(doc.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// return proper response
		resp := Generic{
			Index:  index,
			Type:   _type,
			ID:     ID,
			Found:  true,
			Source: body,
		}

		js, err := json.Marshal(resp)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	} else {
		resp := Generic{
			Index: index,
			Type:  _type,
			ID:    ID,
			Found: false,
		}

		js, err := json.Marshal(resp)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func UpdateDocumentByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	index := vars["index"]
	_type := vars["_type"]
	ID := vars["id"]

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated, err := database.upsertDocument(index, _type, ID, body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := "updated"
	if !updated {
		result = "created"
	}

	resp := &Generic{
		Index:   index,
		Type:    _type,
		ID:      ID,
		Created: !updated,
		Result:  result,
	}

	js, err := json.Marshal(resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func DeleteDocumentByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	index := vars["index"]
	_type := vars["_type"]
	ID := vars["id"]

	deleted := database.deleteDocument(index, _type, ID)

	resp := &Generic{
		Found: deleted,
		ID:    ID,
		Index: index,
		Type:  _type,
	}

	js, err := json.Marshal(resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
