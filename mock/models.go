package mock

import (
	"encoding/json"
)

type (
	// Represents an Elasticsearch Document
	Document struct {
		ID string
		Body map[string]json.RawMessage
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

	SearchHit struct {
		Index string `json:"_index"`
		Type string `json:"_type"`
		ID string `json:"_id"`
		Score float64 `json:"_score"`
		Source json.RawMessage `json:"_source"`
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
			Hits []*SearchHit `json:"hits"`
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