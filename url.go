package elasticsearch

import (
	"github.com/b3ntly/uritemplates"
	"net/url"
)

func constructPath(uriTemplate string, pathMap map[string]string) (string, error) {
	escapedPath, _, err := uritemplates.Expand(uriTemplate, pathMap)
	return escapedPath, err
}

func injectQuerystring(uri string, queryMap map[string]string) string {
	values := &url.Values{}

	for k, v := range queryMap {
		values.Add(k, v)
	}

	return uri + "?" + values.Encode()
}

// Concatenate a URI from an array of strings.
func buildURI(baseURI string, pathMap map[string]string, queryMap map[string]string) (string, error) {
	uri, err := constructPath(baseURI, pathMap)

	if err != nil {
		return "", err
	}

	if queryMap != nil {
		return injectQuerystring(uri, queryMap), nil
	}

	return uri, err
}
