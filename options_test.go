package elasticsearch_test

import (
	"github.com/b3ntly/elasticsearch"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOptions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Options.init will set a default URI", func(t *testing.T) {
		options := &elasticsearch.Options{}
		err := options.Init()

		asserts.Nil(err)
		asserts.Equal("http://127.0.0.1:9200{/index,type,suffix}", options.URI)
	})

	t.Run("Options.init will not override a custom URI", func(t *testing.T) {
		const URL = "http://elasticsearch:9200"
		options := &elasticsearch.Options{URI: URL}
		err := options.Init()

		asserts.Nil(err)
		asserts.Equal("http://elasticsearch:9200{/index,type,suffix}", options.URI)
	})

	t.Run("Options.init will return an error if URL does not declare a protocol", func(t *testing.T) {
		const malformedURL = "127.0.0.1:9200"
		options := &elasticsearch.Options{URI: malformedURL}
		err := options.Init()
		asserts.Error(err)
	})
}
