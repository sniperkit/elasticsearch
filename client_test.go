package elasticsearch_test

import (
	"log"
	"testing"
	"net/http"
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

// returns two clients, one configured for an actual elasticsearch service
// and one configured to use the mock service
func setupClean()(*elasticsearch.Client, *elasticsearch.Client, error){
	client, err := elasticsearch.New(&elasticsearch.Options{})

	if err != nil {
		return nil, nil, err
	}

	mockClient, err := elasticsearch.New(&elasticsearch.Options{
		URL: "http://127.0.0.1:9201",
	})

	if err != nil {
		return nil, nil, err
	}

	// dropping a non-existent index returns an error, which is ok
	_ = client.I(testIndex).Drop()
	_ = mockClient.I(testIndex).Drop()
	return client, mockClient, err
}

func TestMain(m *testing.M){
	go setupMockServer()
	retCode := m.Run()
	os.Exit(retCode)
}

func setupPopulated() (*elasticsearch.Client, *elasticsearch.Client, *elasticsearch.Document, *elasticsearch.Document, error){
	client, mockClient, err := setupClean()

	if err != nil {
		return nil, nil, nil, nil, err
	}

	collection := client.I(testIndex).T(testType)

	body, err := json.Marshal(sampleDocument)

	if err != nil {
		return nil, nil, nil, nil, err
	}

	doc, err := collection.Insert(body)

	if err != nil {
		return nil, nil, nil, nil, err
	}

	mockCollection := mockClient.I(testIndex).T(testType)
	mockDoc, err := mockCollection.Insert(body)

	return client, mockClient, doc, mockDoc, err
}

func TestNew(t *testing.T) {
	client, mockClient, err := setupClean()
	require.Nil(t, err)

	t.Run("ELASTICSEARCH with default values", func(t *testing.T){
		assert.Equal(t, elasticsearch.DefaultURL, client.Options.URL)
		assert.Equal(t, elasticsearch.DefaultHTTPClient, client.Options.HTTPClient)
	})

	t.Run("ELASTICSEARCH with custom options", func(t *testing.T){
		const URL = "elasticsearch:9200"
		var CLIENT = &http.Client{}

		client, err := elasticsearch.New(&elasticsearch.Options{ URL: URL, HTTPClient: CLIENT })

		require.Nil(t, err)
		assert.Equal(t, URL, client.Options.URL)
		assert.Equal(t, CLIENT, client.Options.HTTPClient)
	})

	t.Run("ELASTICSEARCH_MOCK with default values", func(t *testing.T){
		assert.Equal(t, "http://127.0.0.1:9201", mockClient.Options.URL)
		assert.Equal(t, elasticsearch.DefaultHTTPClient, mockClient.Options.HTTPClient)
	})

	t.Run("ELASTICSEARCH_MOCK with custom options", func(t *testing.T){
		const URL = "elasticsearch:9201"
		var CLIENT = &http.Client{}

		mockClient, err := elasticsearch.New(&elasticsearch.Options{ URL: URL, HTTPClient: CLIENT })

		require.Nil(t, err)
		assert.Equal(t, URL, mockClient.Options.URL)
		assert.Equal(t, CLIENT, mockClient.Options.HTTPClient)
	})
}

func TestClient_I(t *testing.T) {
	client, mockClient, err := setupClean()
	require.Nil(t, err)

	t.Run("ELASTICSEARCH create from a client", func(t *testing.T){
		index := client.I(testIndex)
		assert.Equal(t, testIndex, index.Name)
	})

	t.Run("ELASTICSEARCH_MOCK create from a client", func(t *testing.T){
		index := mockClient.I(testIndex)
		assert.Equal(t, testIndex, index.Name)
	})
}

func TestIndex_T(t *testing.T) {
	client, mockClient, err := setupClean()
	require.Nil(t, err)

	t.Run("ELASTICSEARCH create from an Index", func(t *testing.T){
		index := client.I(testIndex)
		assert.Equal(t, testIndex, index.Name)
		collection := index.T(testType)
		assert.Equal(t, testType, collection.Name)
	})

	t.Run("ELASTICSEARCH_MOCK create from an Index", func(t *testing.T){
		index := mockClient.I(testIndex)
		assert.Equal(t, testIndex, index.Name)
		collection := index.T(testType)
		assert.Equal(t, testType, collection.Name)
	})
}

func TestIndex_Search(t *testing.T) {
	client, mockClient, _, _, err := setupPopulated()

	require.Nil(t, err)

	t.Run("ELASTICSEARCH Can search an index", func(t *testing.T){
		collection := client.I(testIndex)
		docs, err := collection.Search("*:*")
		require.Nil(t, err)
		assert.NotEqual(t, 0, len(docs))
	})

	t.Run("ELASTICSEARCH_MOCK Can search an index", func(t *testing.T){
		collection := mockClient.I(testIndex)
		docs, err := collection.Search("*:*")
		require.Nil(t, err)
		assert.NotEqual(t, 0, len(docs))
	})
}

func TestIndex_Drop(t *testing.T) {
	client, mockClient, _, _, err := setupPopulated()
	require.Nil(t, err)

	t.Run("ELASTICSEARCH Can drop an index", func(t *testing.T){
		assert.Nil(t, client.I(testIndex).Drop())
	})

	t.Run("ELASTICSEARCH Searching on a non existent index returns an error", func(t *testing.T){
		collection := client.I(testIndex)
		_, err := collection.Search("*:*")
		require.NotNil(t, err)
	})

	t.Run("ELASTICSEARCH_MOCK Can drop an index", func(t *testing.T){
		assert.Nil(t, mockClient.I(testIndex).Drop())
	})

	t.Run("ELASTICSEARCH_MOCK Searching on a non existent index returns an error", func(t *testing.T){
		collection := mockClient.I(testIndex)
		_, err := collection.Search("*:*")
		require.NotNil(t, err)
	})
}

func TestType_Search(t *testing.T) {
	client, mockClient, _, _, err := setupPopulated()

	require.Nil(t, err)

	t.Run("ELASTICSEARCH Can search a type", func(t *testing.T){
		collection := client.I(testIndex).T(testType)
		docs, err := collection.Search("*:*")
		require.Nil(t, err)
		assert.NotEqual(t, 0, len(docs))
	})

	t.Run("ELASTICSEARCH_MOCK Can search a type", func(t *testing.T){
		collection := mockClient.I(testIndex).T(testType)
		docs, err := collection.Search("*:*")
		t.Log(err)
		require.Nil(t, err)
		assert.NotEqual(t, 0, len(docs))
	})
}

func TestType_Insert(t *testing.T) {
	client, mockClient, err := setupClean()
	require.Nil(t, err)

	t.Run("ELASTICSEARCH Can insert a document", func(t *testing.T){
		collection := client.I(testIndex).T(testType)

		body, err := json.Marshal(sampleDocument)
		require.Nil(t, err)

		doc, err := collection.Insert(body)

		require.Nil(t, err)
		assert.NotEqual(t, "", doc.ID)
	})

	t.Run("ELASTICSEARCH_MOCK Can insert a document", func(t *testing.T){
		collection := mockClient.I(testIndex).T(testType)

		body, err := json.Marshal(sampleDocument)
		require.Nil(t, err)

		doc, err := collection.Insert(body)

		require.Nil(t, err)
		assert.NotEqual(t, "", doc.ID)
	})
}

func TestType_Find(t *testing.T) {
	client, mockClient, _, _, err := setupPopulated()
	require.Nil(t, err)

	t.Run("ELASTICSEARCH Can find documents by property value", func(t *testing.T){
		collection := client.I(testIndex).T(testType)
		docs, err := collection.Search("Message:" + testMessage)
		require.Nil(t, err)
		assert.NotEqual(t, 0, len(docs))
	})

	t.Run("ELASTICSEARCH_MOCK Can find documents by property value", func(t *testing.T){
		collection := mockClient.I(testIndex).T(testType)
		docs, err := collection.Search("Message:" + testMessage)
		require.Nil(t, err)
		assert.NotEqual(t, 0, len(docs))
	})
}

func TestType_FindById(t *testing.T) {
	client, mockClient, doc, mockDoc, err := setupPopulated()
	require.Nil(t, err)

	t.Run("ELASTICSEARCH Can find document by id", func(t *testing.T){
		collection := client.I(testIndex).T(testType)
		result, err := collection.FindById(doc.ID)
		require.Nil(t, err)
		assert.Equal(t, doc.ID, result.ID)
	})

	t.Run("ELASTICSEARCH_MOCK Can find document by id", func(t *testing.T){
		collection := mockClient.I(testIndex).T(testType)
		result, err := collection.FindById(mockDoc.ID)
		require.Nil(t, err)
		assert.Equal(t, mockDoc.ID, result.ID)
	})
}

func TestType_UpdateById(t *testing.T) {
	client, mockClient, doc, mockDoc, err := setupPopulated()
	require.Nil(t, err)

	t.Run("ELASTICSEARCH Can update a document by id", func(t *testing.T){
		collection := client.I(testIndex).T(testType)

		update, err := json.Marshal(&Example{ Message: "World" })
		require.Nil(t, err)

		err = collection.UpdateById(doc.ID, update)
		require.Nil(t, err)

		result, err := collection.FindById(doc.ID)
		require.Nil(t, err)

		finalDoc := &Example{}
		err = json.Unmarshal(result.Body, finalDoc)
		require.Nil(t, err)
		assert.Equal(t, "World", finalDoc.Message)
	})

	t.Run("ELASTICSEARCH_MOCK Can update a document by id", func(t *testing.T){
		collection := mockClient.I(testIndex).T(testType)

		update, err := json.Marshal(&Example{ Message: "World" })
		require.Nil(t, err)

		err = collection.UpdateById(mockDoc.ID, update)
		require.Nil(t, err)

		result, err := collection.FindById(mockDoc.ID)
		require.Nil(t, err)

		finalDoc := &Example{}
		err = json.Unmarshal(result.Body, finalDoc)
		require.Nil(t, err)
		assert.Equal(t, "World", finalDoc.Message)
	})
}

func TestType_DeleteById(t *testing.T) {
	t.Run("Can delete a document by id", func(t *testing.T){

	})

	t.Run("Attempting to delete a document by id which does not exist returns an error", func(t *testing.T){

	})
}