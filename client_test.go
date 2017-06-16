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
	"strings"
)

type Example struct {
	Message string `json:"message"`
}

const (
	testIndex = "test"
	testType = "test"
	testMessage = "hello"
	testMessageChange = "world"
	// the number of elements to insert/update/delete in tests of the bulk API
	bulkOperations = 5
)

var (
	sampleDocument = &Example{ testMessage }
)

func setupMockServer(){
	log.Fatal(mock.New().ListenAndServe())
}

// returns a client with the configured URI and a single document
func getClient(URL string) (*elasticsearch.Client, error){
	client, err := elasticsearch.New(&elasticsearch.Options{ URI: URL })

	if err != nil {
		return nil, err
	}

	clean(client)
	return client, nil
}

func clean(client *elasticsearch.Client){
	// silent error is acceptable here as elasticsearch will throw a 404 when trying to delete nonexistent index
	_ = client.I(testIndex).Drop()
}

func TestMain(m *testing.M){
	go setupMockServer()
	retCode := m.Run()
	os.Exit(retCode)
}

// Integration tests for standard usage
func TestClient(t *testing.T){
	// client which will test against a real elasticsearch service
	client, err := getClient("http://127.0.0.1:9200")
	require.Nil(t, err)

	// client which will test against elasticsearch/mock
	mockClient, err := getClient("http://127.0.0.1:9201")
	require.Nil(t, err)

	clients := []*elasticsearch.Client { client, mockClient }

	for _, client := range clients {
		t.Run("Insert Document", func(t *testing.T){
			t.Run("Returns a newly created ID", func(t *testing.T){
				collection := client.I(testIndex).T(testType)
				body, err := json.Marshal(sampleDocument)
				require.Nil(t, err)
				ID, err := collection.Insert(body)
				require.Nil(t, err)
				assert.NotEqual(t, "", ID)
				clean(client)
			})
		})

		t.Run("Bulk Insert Documents", func(t *testing.T){
			t.Run("Returns the proper number of inserted documents with valid IDs", func(t *testing.T){
				collection := client.I(testIndex).T(testType)
				body, err := json.Marshal(sampleDocument)
				require.Nil(t, err)


				inputs := make([][]byte, bulkOperations)
				for i := 0; i < bulkOperations; i++ {
					doc := make([]byte, len(body))
					copy(doc, body)
					inputs[i] = doc
				}

				docs, err := collection.BulkInsert(inputs)
				require.Nil(t, err)
				require.Equal(t, bulkOperations, len(docs))

				for _, ID := range docs {
					require.NotEqual(t, "", ID)
				}

				clean(client)
			})
		})

		t.Run("Find document by property", func(t *testing.T){
			t.Run("Returns at least one document", func(t *testing.T){
				// ensure at least one document
				collection := client.I(testIndex).T(testType)
				body, err := json.Marshal(sampleDocument)
				require.Nil(t, err)
				ID, err := collection.Insert(body)
				require.Nil(t, err)
				assert.NotEqual(t, "", ID)

				// find the document
				docs, err := collection.Search("message:" + testMessage)
				require.Nil(t, err)
				assert.NotEqual(t, 0, len(docs))
				clean(client)
			})
		})

		t.Run("Find document by ID", func(t *testing.T){
			t.Run("Returns at least one document", func(t *testing.T){
				// ensure at least one document
				collection := client.I(testIndex).T(testType)
				body, err := json.Marshal(sampleDocument)
				require.Nil(t, err)
				ID, err := collection.Insert(body)
				require.Nil(t, err)
				assert.NotEqual(t, "", ID)

				// find that document by ID
				result, err := collection.FindById(ID)
				require.Nil(t, err)
				assert.Equal(t, ID, result.ID)
				clean(client)
			})
		})

		t.Run("Update Document by ID", func(t *testing.T){
			t.Run("will update a document by ID", func(t *testing.T){
				// ensure at least one document
				collection := client.I(testIndex).T(testType)
				body, err := json.Marshal(sampleDocument)
				require.Nil(t, err)
				ID, err := collection.Insert(body)
				require.Nil(t, err)
				assert.NotEqual(t, "", ID)

				// update the document
				update, err := json.Marshal(&Example{ Message: testMessageChange })
				require.Nil(t, err)
				err = collection.UpdateById(ID, update)
				require.Nil(t, err)

				result, err := collection.FindById(ID)
				require.Nil(t, err)

				// verify the document
				finalDoc := &Example{}
				err = json.Unmarshal(result.Body, finalDoc)
				require.Nil(t, err)
				assert.Equal(t, testMessageChange, finalDoc.Message)
				clean(client)
			})
		})

		t.Run("Bulk Update Documents", func(t *testing.T){
			t.Run("Bulk update will properly update the given IDs", func(t *testing.T){
				// insert base documents
				collection := client.I(testIndex).T(testType)
				body, err := json.Marshal(sampleDocument)
				require.Nil(t, err)

				inputs := make([][]byte, bulkOperations)
				for i := 0; i < bulkOperations; i++ {
					doc := make([]byte, len(body))
					copy(doc, body)
					inputs[i] = doc
				}

				IDs, err := collection.BulkInsert(inputs)
				require.Nil(t, err)
				require.Equal(t, bulkOperations, len(IDs))

				for _, ID := range IDs {
					require.NotEqual(t, "", ID)
				}

				// update base documents with different field
				updates := make([]*elasticsearch.Document, bulkOperations)
				updatedBody, err := json.Marshal(&Example{ testMessageChange })
				require.Nil(t, err)

				for i := 0; i < len(IDs); i++ {
					update := &elasticsearch.Document { ID: IDs[i], Body: make([]byte, len(updatedBody)) }
					copy(update.Body, updatedBody)
					updates[i] = update
				}

				updated, err := collection.BulkUpdate(updates)

				require.Nil(t, err)
				require.Equal(t, bulkOperations, len(updated))

				for _, ID := range updated {
					require.NotEqual(t, "", ID)
				}

				clean(client)
			})
		})

		t.Run("Delete Document by ID", func(t *testing.T){
			t.Run("will delete a document by id", func(t *testing.T) {
				// ensure at least one document
				collection := client.I(testIndex).T(testType)
				body, err := json.Marshal(sampleDocument)
				require.Nil(t, err)
				ID, err := collection.Insert(body)
				require.Nil(t, err)
				assert.NotEqual(t, "", ID)

				// delete the document
				err = collection.DeleteById(ID)
				require.Nil(t, err)

				// verify it was deleted
				_, err = collection.FindById(ID)

				// FindById will return an error because the document was deleted
				require.Error(t, err)
				clean(client)
			})
		})

		t.Run("Bulk delete Documents", func(t *testing.T){
			// insert base documents
			collection := client.I(testIndex).T(testType)
			body, err := json.Marshal(sampleDocument)
			require.Nil(t, err)

			inputs := make([][]byte, bulkOperations)
			for i := 0; i < bulkOperations; i++ {
				doc := make([]byte, len(body))
				copy(doc, body)
				inputs[i] = doc
			}

			IDs, err := collection.BulkInsert(inputs)
			require.Nil(t, err)
			require.Equal(t, bulkOperations, len(IDs))

			for _, ID := range IDs {
				require.NotEqual(t, "", ID)
			}

			deleted, err := collection.BulkDelete(IDs...)
			require.Nil(t, err)
			require.Equal(t, len(IDs), len(deleted))

			// verify none of the deleted ids still exist
			// todo: bulk findById would be good here

			// todo: write a better test case. because our mock interface does not actually perform
			// todo: Lucene based queries, it actually returns all documents, we test that the updates were
			// todo: performed by effictevly querying all messages and verifying the update on the ids
			// todo: that were initially inserted
			finalDocs, err := collection.Search("*:*")

			require.Nil(t, err)

			for _, finalDoc := range finalDocs {
				for _, id := range deleted {
					require.NotEqual(t, finalDoc.ID, id)
				}
			}

			clean(client)
		})

		t.Run("Search Index", func(t *testing.T){
			t.Run("will return at least one document", func(t *testing.T){
				// ensure at least one document
				collection := client.I(testIndex).T(testType)
				body, err := json.Marshal(sampleDocument)
				require.Nil(t, err)
				ID, err := collection.Insert(body)
				require.Nil(t, err)
				assert.NotEqual(t, "", ID)


				docs, err := collection.Search("*:*")
				require.Nil(t, err)
				assert.NotEqual(t, 0, len(docs))
				clean(client)

			})

			t.Run("searching on a non existent index returns an error", func(t *testing.T){
				collection := mockClient.I("hunter2")
				_, err := collection.Search("*:*")
				require.NotNil(t, err)
				clean(client)
			})
		})

		t.Run("Search Type", func(t *testing.T){
			t.Run("returns at least one document", func(t *testing.T){
				// ensure at least one document
				collection := client.I(testIndex).T(testType)
				body, err := json.Marshal(sampleDocument)
				require.Nil(t, err)
				ID, err := collection.Insert(body)
				require.Nil(t, err)
				assert.NotEqual(t, "", ID)

				docs, err := collection.Search("*:*")
				require.Nil(t, err)
				assert.NotEqual(t, 0, len(docs))
				clean(client)
			})
		})

		t.Run("Drop Index", func(t *testing.T){
			t.Run("will drop an index without error", func(t *testing.T){
				// ensure the index exists to begin with
				collection := client.I(testIndex).T(testType)
				body, err := json.Marshal(sampleDocument)
				require.Nil(t, err)
				ID, err := collection.Insert(body)
				require.Nil(t, err)
				assert.NotEqual(t, "", ID)

				assert.Nil(t, client.I(testIndex).Drop())
				clean(client)
			})
		})

		t.Run("Test sqlSearch", func(t *testing.T){
			// skip this test for the mock client
			if strings.Contains(client.Options.URI, "http://127.0.0.1:9201"){
				t.Skip()
			}

			t.Run("will return matching documents on a simple query", func(t *testing.T){
				collection := client.I(testIndex).T(testType)

				// insert 13 identical documents
				example := &Example{Message:"eureka"}
				body, err := json.Marshal(example)
				require.Nil(t, err)

				inputs := make([][]byte, bulkOperations)
				for i := 0; i < bulkOperations; i++ {
					doc := make([]byte, len(body))
					copy(doc, body)
					inputs[i] = doc
				}

				docs, err := collection.BulkInsert(inputs)
				require.Nil(t, err)
				require.Equal(t, bulkOperations, len(docs))

				sql := `SELECT Message FROM test WHERE Message = 'eureka' LIMIT 14`
				results, err := collection.SearchSQL(sql)
				require.Nil(t, err)
				require.Equal(t, bulkOperations, len(results))
			})
		})
	}
}