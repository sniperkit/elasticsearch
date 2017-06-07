package elasticsearch

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


// table test
func TestREST_BuildURL(t *testing.T) {
	asserts := assert.New(t)
	r := &rest{ BaseURL: DefaultURL, HTTPClient: DefaultHTTPClient }
	base := r.BaseURL
	asserts.Equal(DefaultURL, base)

	cases := []struct{
		inputs []string
		expectedOutput string
	}{
		{[]string{"hello", "world"}, base + "/hello/world"},
		{[]string{"_under", "q?query=thing"}, base + "/_under/q?query=thing"},
	}

	for _, test := range cases {
		output := r.buildURL(test.inputs...)
		asserts.Equal(test.expectedOutput, output)
	}
}

func TestREST_BuildRequest(t *testing.T) {

}

func TestREST_SendRequest(t *testing.T) {

}

func TestREST_Request(t *testing.T) {

}