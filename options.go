package elasticsearch

import (
	"net/url"
	"net/http"
	"github.com/hashicorp/go-cleanhttp"
)

type Options struct {
	URL string
	HTTPClient *http.Client
}

var (
	DefaultURL = "http://127.0.0.1:9200"
	DefaultHTTPClient = cleanhttp.DefaultClient()
)

func (opts *Options) Init() error {
	if opts.URL == "" {
		opts.URL = DefaultURL
	}

	uri, err := url.Parse(opts.URL)

	if err != nil {
		return err
	}

	opts.URL = uri.String()

	if opts.HTTPClient == nil {
		opts.HTTPClient = DefaultHTTPClient
	}

	return nil
}