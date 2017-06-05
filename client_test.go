package elasticsearch_test

import (
	"log"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/b3ntly/elasticsearch"
	"github.com/b3ntly/elasticsearch/mock"
	"encoding/json"
	"os"
)

type Example struct {
	Message string
}

var (
	testIndex = "test"
	testType = "test"
	testMessage = "hello"
	sampleDocument = &Example{ testMessage }
)

func setupMockServer(){
	log.Fatal(mock.New().ListenAndServe())
}

// returns a client with the configured URL and a single document
func getClient(URL string) (*elasticsearch.Client, error){
	client, err := elasticsearch.New(&elasticsearch.Options{ URL: URL })

	if err != nil {
		return nil, err
	}

	// silent error is acceptable here as elasticsearch will throw a 404 when trying to delete nonexistent index
	_ = client.I(testIndex).Drop()
	return client, nil
}

func TestMain(m *testing.M){
	go setupMockServer()
	retCode := m.Run()
	os.Exit(retCode)
}

// Integration tests for standard usage
func TestClient(t *testing.T){
	// client which will test against a real elasticsearch service
	realClient, err := getClient("http://127.0.0.1:9200")
	require.Nil(t, err)

	// client which will test against elasticsearch/mock
	mockClient, err := getClient("http://127.0.0.1:9201")
	require.Nil(t, err)

	clients := []*elasticsearch.Client { realClient, mockClient }

	for _, client := range clients {
		t.Run("Insert Document", func(t *testing.T){
			t.Run("Returns a newly created ID", func(t *testing.T){
				collection := client.I(testIndex).T(testType)
				body, err := json.Marshal(sampleDocument)
				require.Nil(t, err)
				doc, err := collection.Insert(body)
				require.Nil(t, err)
				assert.NotEqual(t, "", doc.ID)
			})
		})

		t.Run("Find document by property", func(t *testing.T){
			t.Run("Returns at least one document", func(t *testing.T){
				// ensure at least one document
				collection := client.I(testIndex).T(testType)
				body, err := json.Marshal(sampleDocument)
				require.Nil(t, err)
				doc, err := collection.Insert(body)
				require.Nil(t, err)
				assert.NotEqual(t, "", doc.ID)

				// find the document
				docs, err := collection.Search("Message:" + testMessage)
				require.Nil(t, err)
				assert.NotEqual(t, 0, len(docs))
			})
		})

		t.Run("Find document by ID", func(t *testing.T){
			t.Run("Returns at least one document", func(t *testing.T){
				// ensure at least one document
				collection := client.I(testIndex).T(testType)
				body, err := json.Marshal(sampleDocument)
				require.Nil(t, err)
				doc, err := collection.Insert(body)
				require.Nil(t, err)
				assert.NotEqual(t, "", doc.ID)

				// find that document by ID
				result, err := collection.FindById(doc.ID)
				require.Nil(t, err)
				assert.Equal(t, doc.ID, result.ID)
			})
		})

		t.Run("Update Document by ID", func(t *testing.T){
			t.Run("will update a document by ID", func(t *testing.T){
				// ensure at least one document
				collection := client.I(testIndex).T(testType)
				body, err := json.Marshal(sampleDocument)
				require.Nil(t, err)
				doc, err := collection.Insert(body)
				require.Nil(t, err)
				assert.NotEqual(t, "", doc.ID)

				// update the document
				update, err := json.Marshal(&Example{ Message: "World" })
				require.Nil(t, err)
				err = collection.UpdateById(doc.ID, update)
				require.Nil(t, err)

				result, err := collection.FindById(doc.ID)
				require.Nil(t, err)

				// verify the document
				finalDoc := &Example{}
				err = json.Unmarshal(result.Body, finalDoc)
				require.Nil(t, err)
				assert.Equal(t, "World", finalDoc.Message)
			})
		})

		t.Run("Delete Document by ID", func(t *testing.T){
			t.Run("will delete a document by id", func(t *testing.T) {
				// ensure at least one document
				collection := client.I(testIndex).T(testType)
				body, err := json.Marshal(sampleDocument)
				require.Nil(t, err)
				doc, err := collection.Insert(body)
				require.Nil(t, err)
				assert.NotEqual(t, "", doc.ID)

				// delete the document
				err = collection.DeleteById(doc.ID)
				require.Nil(t, err)

				// verify it was deleted
				_, err = collection.FindById(doc.ID)

				// FindById will return an error because the document was deleted
				require.Error(t, err)
			})
		})

		t.Run("Search Index", func(t *testing.T){
			t.Run("will return at least one document", func(t *testing.T){
				// ensure at least one document
				collection := client.I(testIndex).T(testType)
				body, err := json.Marshal(sampleDocument)
				require.Nil(t, err)
				doc, err := collection.Insert(body)
				require.Nil(t, err)
				assert.NotEqual(t, "", doc.ID)


				docs, err := collection.Search("*:*")
				require.Nil(t, err)
				assert.NotEqual(t, 0, len(docs))
			})

			t.Run("searching on a non existent index returns an error", func(t *testing.T){
				collection := mockClient.I("hunter2")
				_, err := collection.Search("*:*")
				require.NotNil(t, err)
			})
		})

		t.Run("Search Type", func(t *testing.T){
			t.Run("returns at least one document", func(t *testing.T){
				// ensure at least one document
				collection := client.I(testIndex).T(testType)
				body, err := json.Marshal(sampleDocument)
				require.Nil(t, err)
				doc, err := collection.Insert(body)
				require.Nil(t, err)
				assert.NotEqual(t, "", doc.ID)

				docs, err := collection.Search("*:*")
				require.Nil(t, err)
				assert.NotEqual(t, 0, len(docs))
			})
		})

		t.Run("Drop Index", func(t *testing.T){
			t.Run("will drop an index without error", func(t *testing.T){
				// ensure the index exists to begin with
				collection := client.I(testIndex).T(testType)
				body, err := json.Marshal(sampleDocument)
				require.Nil(t, err)
				doc, err := collection.Insert(body)
				require.Nil(t, err)
				assert.NotEqual(t, "", doc.ID)

				assert.Nil(t, client.I(testIndex).Drop())
			})
		})
	}
}