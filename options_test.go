package elasticsearch_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/b3ntly/elasticsearch"
)

func TestOptions(t *testing.T){
	asserts := assert.New(t)

	t.Run("Options.init will set a default URL", func(t *testing.T){
		options := &elasticsearch.Options{}
		err := options.Init()

		asserts.Nil(err)
		asserts.Equal(elasticsearch.DefaultURL, options.URL)
	})

	t.Run("Options.init will not override a custom URL", func(t *testing.T){
		const URL = "elasticsearch:9200"
		options := &elasticsearch.Options{ URL: URL }
		err := options.Init()

		asserts.Nil(err)
		asserts.Equal(URL, options.URL)
	})
}