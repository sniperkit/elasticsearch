package elasticsearch_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/b3ntly/elasticsearch"
)


// table test
func TestREST_BuildURL(t *testing.T) {
	asserts := assert.New(t)
	r := &elasticsearch.REST{ BaseURL: elasticsearch.DefaultURL, HTTPClient: elasticsearch.DefaultHTTPClient }
	base := r.BaseURL
	asserts.Equal(elasticsearch.DefaultURL, base)

	cases := []struct{
		inputs []string
		expectedOutput string
	}{
		{[]string{"hello", "world"}, base + "/hello/world"},
		{[]string{"_under", "q?query=thing"}, base + "/_under/q?query=thing"},
	}

	for _, test := range cases {
		output := r.BuildURL(test.inputs...)
		asserts.Equal(test.expectedOutput, output)
	}
}

func TestREST_BuildRequest(t *testing.T) {

}

func TestREST_SendRequest(t *testing.T) {

}

func TestREST_Request(t *testing.T) {

}