package elasticsearch

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_buildURI(t *testing.T) {
	cases := []struct {
		base     string
		pathMap  map[string]string
		queryMap map[string]string
		expected string
	}{
		{
			base: "foo.com{/baz,bat,suffix}",
			pathMap: map[string]string{
				"baz":    "batman",
				"bat":    "superman",
				"suffix": "_search",
			},
			queryMap: map[string]string{
				"q": "jobs",
			},
			expected: "foo.com/batman/superman/_search?q=jobs",
		},
		{
			base: "http://127.0.0.1:9200{/index,type,suffix}",
			pathMap: map[string]string{
				"index":  "test",
				"type":   "test",
				"suffix": "_search",
			},
			queryMap: map[string]string{
				"q": "foo",
			},
			expected: "http://127.0.0.1:9200/test/test/_search?q=foo",
		},
		{
			base: "http://127.0.0.1:9200{/index,type,suffix}",
			pathMap: map[string]string{
				"suffix": "_search",
			},
			queryMap: map[string]string{
				"q": "foo",
			},
			expected: "http://127.0.0.1:9200/_search?q=foo",
		},
	}

	for _, test := range cases {
		output, err := buildURI(test.base, test.pathMap, test.queryMap)

		require.Nil(t, err)
		require.Equal(t, test.expected, output)
	}
}

func Test_buildURI2(t *testing.T) {
	malformedURITemplate := "http://127.0.0.1:9200{/foo"
	pathMap := map[string]string{"bar": "baz "}
	output, err := buildURI(malformedURITemplate, pathMap, nil)
	require.Error(t, err)
	require.Equal(t, "", output)
}
