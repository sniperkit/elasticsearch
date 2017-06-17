package elasticsearch

import (
	"github.com/hashicorp/go-cleanhttp"
	"net/http"
	"net/url"
)

type Options struct {
	URI        string
	HTTPClient *http.Client
}

var (
	DefaultURL        = "http://127.0.0.1:9200"
	DefaultHTTPClient = cleanhttp.DefaultClient()
)

func (opts *Options) Init() error {
	if opts.URI == "" {
		opts.URI = DefaultURL
	} else {
		uri, err := url.Parse(opts.URI)

		if err != nil {
			return err
		}

		opts.URI = uri.String()
	}

	// add templating suffix
	opts.URI = opts.URI + "{/index,type,suffix}"

	if opts.HTTPClient == nil {
		opts.HTTPClient = DefaultHTTPClient
	}

	return nil
}
