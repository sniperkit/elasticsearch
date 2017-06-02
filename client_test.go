package elasticsearch_test

import (
	"time"
	"testing"
	"net/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/b3ntly/elasticsearch"
	"encoding/json"
)

type Example struct {
	Message string
}

var (
	ELASTICSEARCH_LAG_DELAY = time.Millisecond * 1
	testIndex = "test"
	testType = "test"
	testMessage = "hello"
	sampleDocument = &Example{ testMessage }
)

func setupClean()(*elasticsearch.Client, error){
	client, err := elasticsearch.New(&elasticsearch.Options{})

	if err != nil {
		return nil, err
	}

	// dropping a non-existent index returns an error, which is ok
	_ = client.I(testIndex).Drop()
	return client, err
}

func setupPopulated() (*elasticsearch.Client, *elasticsearch.Document, error){
	client, err := setupClean()

	if err != nil {
		return nil, nil, err
	}

	collection := client.I(testIndex).T(testType)

	body, err := json.Marshal(sampleDocument)

	if err != nil {
		return nil, nil, err
	}

	doc, err := collection.Insert(body)
	return client, doc, err
}

func TestNew(t *testing.T) {
	client, err := setupClean()
	require.Nil(t, err)

	t.Run("with default values", func(t *testing.T){
		assert.Equal(t, elasticsearch.DefaultURL, client.Options.URL)
		assert.Equal(t, elasticsearch.DefaultHTTPClient, client.Options.HTTPClient)
	})

	t.Run("with custom options", func(t *testing.T){
		const URL = "elasticsearch:9200"
		var CLIENT = &http.Client{}

		client, err := elasticsearch.New(&elasticsearch.Options{ URL: URL, HTTPClient: CLIENT })

		require.Nil(t, err)
		assert.Equal(t, URL, client.Options.URL)
		assert.Equal(t, CLIENT, client.Options.HTTPClient)
	})
}

func TestClient_I(t *testing.T) {
	client, err := setupClean()
	require.Nil(t, err)

	t.Run("create from a client", func(t *testing.T){
		index := client.I(testIndex)
		assert.Equal(t, testIndex, index.Name)
	})
}

func TestIndex_T(t *testing.T) {
	client, err := setupClean()
	require.Nil(t, err)

	t.Run("create from an Index", func(t *testing.T){
		index := client.I(testIndex)
		assert.Equal(t, testIndex, index.Name)
		collection := index.T(testType)
		assert.Equal(t, testType, collection.Name)
	})
}

func TestIndex_Search(t *testing.T) {
	client, _, err := setupPopulated()

	// elasticsearch uses eventual consistence so leave this
	// abomination of an integration test here until I have time
	// to properly mock the elasticsearch interface for testing
	//time.Sleep(ELASTICSEARCH_LAG_DELAY)

	require.Nil(t, err)

	t.Run("Can search an index", func(t *testing.T){
		collection := client.I(testIndex)
		docs, err := collection.Search("*:*")
		require.Nil(t, err)
		assert.NotEqual(t, 0, len(docs))
	})
}

func TestIndex_Drop(t *testing.T) {
	client, _, err := setupPopulated()
	require.Nil(t, err)

	t.Run("Can drop an index", func(t *testing.T){
		assert.Nil(t, client.I(testIndex).Drop())
	})

	t.Run("Searching on a non existent index returns an error", func(t *testing.T){
		collection := client.I(testIndex)
		_, err := collection.Search("*:*")
		require.NotNil(t, err)
	})
}

func TestType_Search(t *testing.T) {
	client, _, err := setupPopulated()

	// elasticsearch uses eventual consistence so leave this
	// abomination of an integration test here until I have time
	// to properly mock the elasticsearch interface for testing
	//time.Sleep(ELASTICSEARCH_LAG_DELAY)

	require.Nil(t, err)

	t.Run("Can search a type", func(t *testing.T){
		collection := client.I(testIndex).T(testType)
		docs, err := collection.Search("*:*")
		require.Nil(t, err)
		assert.NotEqual(t, 0, len(docs))
	})
}

func TestType_Insert(t *testing.T) {
	client, err := setupClean()
	require.Nil(t, err)

	t.Run("Can insert a document", func(t *testing.T){
		collection := client.I(testIndex).T(testType)

		body, err := json.Marshal(sampleDocument)
		require.Nil(t, err)

		doc, err := collection.Insert(body)

		require.Nil(t, err)
		assert.NotEqual(t, "", doc.ID)
	})
}

func TestType_Find(t *testing.T) {
	client, _, err := setupPopulated()

	// elasticsearch uses eventual consistence so leave this
	// abomination of an integration test here until I have time
	// to properly mock the elasticsearch interface for testing
	//time.Sleep(ELASTICSEARCH_LAG_DELAY)

	require.Nil(t, err)

	t.Run("Can find documents by property value", func(t *testing.T){
		collection := client.I(testIndex).T(testType)
		docs, err := collection.Search("Message:" + testMessage)
		require.Nil(t, err)
		assert.NotEqual(t, 0, len(docs))
	})
}

func TestType_FindById(t *testing.T) {
	client, doc, err := setupPopulated()

	// elasticsearch uses eventual consistence so leave this
	// abomination of an integration test here until I have time
	// to properly mock the elasticsearch interface for testing
	//time.Sleep(ELASTICSEARCH_LAG_DELAY)

	require.Nil(t, err)

	t.Run("Can find document by id", func(t *testing.T){
		collection := client.I(testIndex).T(testType)
		result, err := collection.FindById(doc.ID)
		require.Nil(t, err)
		assert.Equal(t, doc.ID, result.ID)
	})
}

func TestType_UpdateById(t *testing.T) {
	client, doc, err := setupPopulated()

	// elasticsearch uses eventual consistence so leave this
	// abomination of an integration test here until I have time
	// to properly mock the elasticsearch interface for testing
	//time.Sleep(ELASTICSEARCH_LAG_DELAY)

	require.Nil(t, err)

	t.Run("Can update a document by id", func(t *testing.T){
		collection := client.I(testIndex).T(testType)

		update, err := json.Marshal(&Example{ Message: "World" })
		require.Nil(t, err)

		err = collection.UpdateById(doc.ID, update)
		require.Nil(t, err)

		// give it time to propagate
		//time.Sleep(ELASTICSEARCH_LAG_DELAY)

		result, err := collection.FindById(doc.ID)
		require.Nil(t, err)

		finalDoc := &Example{}
		err = json.Unmarshal(result.Body, finalDoc)
		require.Nil(t, err)
		assert.Equal(t, "World", finalDoc.Message)
	})
}

func TestType_DeleteById(t *testing.T) {
	//client, doc, err := setupPopulated()

	// elasticsearch uses eventual consistence so leave this
	// abomination of an integration test here until I have time
	// to properly mock the elasticsearch interface for testing
	time.Sleep(ELASTICSEARCH_LAG_DELAY)

	//require.Nil(t, err)

	t.Run("Can delete a document by id", func(t *testing.T){

	})

	t.Run("Attempting to delete a document by id which does not exist returns an error", func(t *testing.T){

	})
}