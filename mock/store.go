package mock

import (
	"encoding/json"
	"errors"
	"sync"
)

// store provides the facilities for replicating in-memory operations
// via a flat map structure
type store struct {
	sync.Mutex

	// index:type:ids:document
	Indexes map[string]map[string]map[string]*Document
}

func newStore() *store {
	return &store{
		Indexes: make(map[string]map[string]map[string]*Document),
	}
}

// helpers that should be called only in a safe (locked) context
func (s *store) getOrCreateIndex(name string) map[string]map[string]*Document {
	index, exists := s.Indexes[name]

	if !exists {
		index = make(map[string]map[string]*Document)
		s.Indexes[name] = index
	}
	return index
}

func (s *store) getOrCreateType(index string, name string) map[string]*Document {
	indice := s.getOrCreateIndex(index)
	_type, exists := indice[name]

	if !exists {
		_type = make(map[string]*Document)
		indice[name] = _type
	}
	return _type
}

// endhelpers

func (s *store) insert(index string, _type string, payload []byte) (*Document, error) {
	s.Lock()
	defer s.Unlock()

	document := &Document{ID: ULID()}
	err := json.Unmarshal(payload, &document.Body)

	if err != nil {
		return nil, err
	}

	collection := s.getOrCreateType(index, _type)
	collection[document.ID] = document
	return document, nil
}

// search index currently returns the entire store
func (s *store) searchIndex(index string) ([]*SearchHit, error) {
	s.Lock()
	defer s.Unlock()
	hits := []*SearchHit{}

	// todo: Build out a data model for elasticsearch/mock o(n^3) is too much
	if indice, exists := s.Indexes[index]; exists {
		for typeName, _type := range indice {
			for _, doc := range _type {
				body, _ := json.Marshal(doc.Body)
				hit := &SearchHit{
					ID:     doc.ID,
					Index:  index,
					Type:   typeName,
					Score:  0.0,
					Source: body,
				}

				hits = append(hits, hit)
			}
		}
	} else {
		return nil, errors.New("Index does not exist.")
	}

	return hits, nil
}

// search type currently returns the entire store
func (s *store) searchType(index string, _type string) ([]*SearchHit, error) {
	s.Lock()
	defer s.Unlock()
	hits := []*SearchHit{}

	if collection, exists := s.Indexes[index][_type]; exists {
		for _, doc := range collection {
			body, _ := json.Marshal(doc.Body)
			hit := &SearchHit{
				ID:     doc.ID,
				Index:  index,
				Type:   _type,
				Score:  0.0,
				Source: body,
			}

			hits = append(hits, hit)
		}
	} else {
		return nil, errors.New("Type does not exist")
	}

	return hits, nil
}

func (s *store) deleteIndex(name string) {
	s.Lock()
	defer s.Unlock()
	delete(s.Indexes, name)
}

func (s *store) getDocument(index string, _type string, ID string) *Document {
	s.Lock()
	defer s.Unlock()

	if doc, exists := s.Indexes[index][_type][ID]; exists {
		return doc
	}

	return nil
}

func (s *store) upsertDocument(index string, _type string, ID string, body []byte) (bool, error) {
	s.Lock()
	defer s.Unlock()
	var err error
	var updated bool

	if document, exists := s.Indexes[index][_type][ID]; !exists {
		_, err = s.insert(index, _type, body)
		updated = false
	} else {
		update := &map[string]json.RawMessage{}
		err = json.Unmarshal(body, update)

		if err != nil {
			return false, err
		}

		for k, v := range *update {
			document.Body[k] = v
		}

		s.Indexes[index][_type][ID] = document
		updated = true
	}

	return updated, err
}

func (s *store) deleteDocument(index string, _type string, ID string) bool {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.Indexes[index][_type][ID]; !exists {
		return false
	}

	delete(s.Indexes[index][_type], ID)
	return true
}
