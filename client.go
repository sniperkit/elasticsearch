package elasticsearch

type (
	// Client interface for this library
	Client struct {
		Options *Options
		// Convenience wrapper for a rest interface with the elasticsearch server
		REST *rest
	}

	// Reference to an elasticsearch index
	Index struct {
		Client *Client
		Name   string
	}

	// Reference to an elasticsearch type - akin to a MongoDB collection
	Type struct {
		Index *Index
		Name  string
	}
)

// Instantiate a new Client with a passed Options object. Call options.init() to
// replace zero-values with default values where desired.
func New(options *Options) (*Client, error){
	err := options.Init()
	r := &rest{ BaseURL: options.URL, HTTPClient: options.HTTPClient }
	return &Client{ Options: options, REST: r }, err
}

// Index creates a reference to an elasticsearch index.
// It will not create the index as elasticsearch default behavior
// is to create an underlying index if an operation references it and
// it does not already exist.
func (c *Client) I(name string) *Index {
	return &Index{ Client: c, Name: name }
}

// Type creates a reference to an elasticsearch type which is a namespace within
// an index. It will not create the type as the default behavior is to create a type
// if an operation references one that does not already exist.
func (idx *Index) T(name string) *Type {
	return &Type{ Index: idx, Name: name }
}

// Perform a basic elasticsearch query on an index that will return exact matches
// on the passed querystring.
func (idx *Index) Search(querystring string)([]*Document, error){
	return idx.Client.REST.searchIndex(idx.Name, querystring)
}

// Delete an index.
func (idx *Index) Drop() error {
	return idx.Client.REST.deleteIndex(idx.Name)
}

// Perform a basic elasticsearch on a given index-type that will return
// exact string matches on the passed querystring.
func (t *Type) Search(querystring string)([]*Document, error){
	return t.Index.Client.REST.searchType(t.Index.Name, t.Name, querystring)
}

// Insert a document into a given type namespace
func (t *Type) Insert(doc []byte) (string, error){
	return t.Index.Client.REST.insertDocument(t.Index.Name, t.Name, doc)
}

// Insert multiple documents into a given type namespace, not all documents may be inserted
// an error will be returned if any of the operations fail
func (t *Type) BulkInsert(docs [][]byte)([]string, error){
	return t.Index.Client.REST.bulkInsertDocuments(t.Index.Name, t.Name, docs)
}

// Find multiple documents in a given type namespace that match
// key:value pairs in the passed queryString.
func (t *Type) Find(querystring string)([]*Document, error){
	return t.Index.Client.REST.searchType(t.Index.Name, t.Name, querystring)
}

// Return a single document by its ID. If the document is not found
// it will return an error.
func (t *Type) FindById(ID string)(*Document, error){
	return t.Index.Client.REST.getDocument(t.Index.Name, t.Name, ID)

}

// Update a document by its ID. If it is not found it will return an error.
func (t *Type) UpdateById(ID string, doc []byte) error {
	return t.Index.Client.REST.updateDocument(t.Index.Name, t.Name, ID, doc)

}

// Insert multiple documents into a given type namespace, not all updates may be completed
// an error will be returned if any of the operations fail
func (t *Type) BulkUpdate(docs []*Document) ([]string, error) {
	return t.Index.Client.REST.bulkUpdateDocuments(t.Index.Name, t.Name, docs)
}

// Delete a document by its ID. If it is not found it will return an error.
func (t *Type) DeleteById(ID string) error {
	return t.Index.Client.REST.deleteDocument(t.Index.Name, t.Name, ID)
}

// delete a list of documents, not all documents may be deleted
// an error will be returned if any of the operations fail
func (t *Type) BulkDelete(IDs ...string) ([]string, error) {
	return t.Index.Client.REST.bulkDeleteDocuments(t.Index.Name, t.Name, IDs)
}