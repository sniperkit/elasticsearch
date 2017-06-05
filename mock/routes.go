package mock

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/gorilla/mux"
)

func DeleteIndex(w http.ResponseWriter, req *http.Request){
	vars := mux.Vars(req)
	index := vars["index"]

	database.deleteIndex(index)

	// return proper response
	resp := DeleteIndexResponse{
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

func SearchIndex(w http.ResponseWriter, req *http.Request){
	vars := mux.Vars(req)
	index := vars["index"]

	hits, err := database.searchIndex(index)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := SearchResponse{}
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

func SearchType(w http.ResponseWriter, req *http.Request){
	// todo: fix nomenclature
	vars := mux.Vars(req)
	index := vars["index"]
	_type := vars["_type"]

	hits, err := database.searchType(index, _type)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := SearchResponse{}
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

func InsertDocument(w http.ResponseWriter, req *http.Request){
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
	resp := IndexResponse{
		Index: index,
		Type: _type,
		ID: doc.ID,
		Created: true,
		Result: "created",
	}

	js, err := json.Marshal(resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func GetDocumentByID(w http.ResponseWriter, req *http.Request){
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
		resp := GetDocumentResponse{
			Index: index,
			Type: _type,
			ID: ID,
			Found: true,
			Document: body,
		}

		js, err := json.Marshal(resp)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	} else {
		resp := GetDocumentResponse{
			Index: index,
			Type: _type,
			ID: ID,
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

func UpdateDocumentByID(w http.ResponseWriter, req *http.Request){
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
	if !updated { result = "created" }

	resp := &UpdateDocumentResponse{
		Index: index,
		Type: _type,
		ID: ID,
		Created: !updated,
		Result: result,
	}

	js, err := json.Marshal(resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func DeleteDocumentByID(w http.ResponseWriter, req *http.Request){
	vars := mux.Vars(req)
	index := vars["index"]
	_type := vars["_type"]
	ID := vars["id"]

	deleted := database.deleteDocument(index, _type, ID)

	resp := &DeleteDocumentResponse{
		Found: deleted,
		ID: ID,
		Index: index,
		Type: _type,
	}

	js, err := json.Marshal(resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}