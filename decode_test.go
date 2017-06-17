package elasticsearch

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_decodingError(t *testing.T) {
	malformedJSON := []byte("l23kej230")
	errorJSON := []byte(`{"found": "false", "created": "false", "acknowledged":"false"}`)
	bulkErrorJSON := []byte(`
		"items": [
			"index": {"created": "false", "acknowledged": "false"},
			"delete": {"created": "false", "acknowledged": "false"}
		]
	`)

	t.Run("Will return error if it cannot unmarshal to mock.Generic", func(t *testing.T) {
		// run it against all decoding functions
		_, err := indexResponseToDocument(malformedJSON)
		require.Error(t, err)

		require.Error(t, deleteDocumentResponseToDocument(malformedJSON))

		_, err = getDocumentResponseToDocument(malformedJSON)
		require.Error(t, err)

		_, err = searchResponseToDocument(malformedJSON)
		require.Error(t, err)

		require.Error(t, deleteDocumentResponseToDocument(malformedJSON))

		require.Error(t, updateDocumentResponseToDocument(malformedJSON))

		_, err = bulkInsertResponseToIDs(malformedJSON)
		require.Error(t, err)

		_, err = bulkUpdateResponseToIDs(malformedJSON)
		require.Error(t, err)

		_, err = bulkDeleteResponseToIDs(malformedJSON)
		require.Error(t, err)
	})

	t.Run("Will return error if state requirement is not met", func(t *testing.T) {
		_, err := indexResponseToDocument(errorJSON)
		require.Error(t, err)

		require.Error(t, deleteDocumentResponseToDocument(errorJSON))

		_, err = getDocumentResponseToDocument(errorJSON)
		require.Error(t, err)

		require.Error(t, deleteIndexResponseToDocument(errorJSON))

		require.Error(t, updateDocumentResponseToDocument(errorJSON))

		_, err = bulkInsertResponseToIDs(bulkErrorJSON)
		require.Error(t, err)

		_, err = bulkDeleteResponseToIDs(bulkErrorJSON)
		require.Error(t, err)
	})
}
