package elasticsearch

type (
	// Client interface for this library
	Client struct {
		Options *Options
		// Convenience wrapper for a REST interface with the elasticsearch server
		REST *REST
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
	r := &REST{ BaseURL: options.URL, HTTPClient: options.HTTPClient }
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
	return idx.Client.REST.SearchIndex(idx.Name, querystring)
}

// Delete an index.
func (idx *Index) Drop() error {
	return idx.Client.REST.DeleteIndex(idx.Name)
}

// Perform a basic elasticsearch on a given index-type that will return
// exact string matches on the passed querystring.
func (t *Type) Search(querystring string)([]*Document, error){
	return t.Index.Client.REST.SearchType(t.Index.Name, t.Name, querystring)
}

// Insert a document into a given type namespace
func (t *Type) Insert(doc []byte) (*Document, error){
	return t.Index.Client.REST.InsertDocument(t.Index.Name, t.Name, doc)
}

// Find multiple documents in a given type namespace that match
// key:value pairs in the passed queryString.
func (t *Type) Find(querystring string)([]*Document, error){
	return t.Index.Client.REST.SearchType(t.Index.Name, t.Name, querystring)
}

// Return a single document by its ID. If the document is not found
// it will return an error.
func (t *Type) FindById(ID string)(*Document, error){
	return t.Index.Client.REST.GetDocument(t.Index.Name, t.Name, ID)

}

// Update a document by its ID. If it is not found it will return an error.
func (t *Type) UpdateById(ID string, doc []byte) error {
	return t.Index.Client.REST.UpdateDocument(t.Index.Name, t.Name, ID, doc)

}

// Delete a document by its ID. If it is not found it will return an error.
func (t *Type) DeleteById(ID string) error {
	return t.Index.Client.REST.DeleteDocument(t.Index.Name, t.Name, ID)
}



