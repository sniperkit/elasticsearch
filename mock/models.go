package mock

import (
	"encoding/json"
)

type (
	Generic struct {
		TimedOut bool `json:"timed_out"`
		Took int `json:"took"`
		Index string `json:"_index"`
		Type string `json:"_type"`
		ID string `json:"_id"`
		Version int64 `json:"_version"`
		Created bool `json:"created"`
		Result string `json:"result"`
		Score float64 `json:"_score"`
		Source json.RawMessage `json:"_source"`
		Hits SearchResult `json:"hits"`
		Status int
		Error ElasticsearchError `json:"error"`
		Acknowledged bool `json:"acknowledged"`
		Found bool `json:"found"`
		Errors bool `json:"errors"`
		Total int `json:"total"`
		MaxScore float64 `json:"max_score"`

		// field for bulk API
		Items []*Operation `json:"items"`
	}

	// Indicated a Bulk API operation
	Operation struct {
		Index *Generic `json:"index"`
		Update *Generic `json:"update"`
		Delete *Generic `json:"delete"`
	}

	ShardMetadata struct {
		Total int `json:"total"`
		Successful int `json:"successful"`
		Failed int `json:"failed"`
	}

	// Represents an Elasticsearch Document
	Document struct {
		ID string
		Body map[string]json.RawMessage
	}

	SearchResult struct {
		Total int `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits []*SearchHit `json:"hits"`
	}

	SearchHit struct {
		Index string `json:"_index"`
		Type string `json:"_type"`
		ID string `json:"_id"`
		Score float64 `json:"_score"`
		Source json.RawMessage `json:"_source"`
	}

	ElasticsearchError struct {
		RootCause []ErrorDescription `json:"root_cause"`
		Type string `json:"type"`
		Reason string `json:"reason"`
	}

	ErrorDescription struct {
		Type string `json:"type"`
		Reason string `json:"reason"`
	}

	ElasticsearchErrors struct {
		Errors ElasticsearchError `json:"errors"`
		Status int `json:"status"`
	}
)