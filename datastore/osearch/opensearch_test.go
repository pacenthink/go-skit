package osearch

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Create_Delete_Index(t *testing.T) {
	client, _ := NewOpenSearchClient(defaultUsername, defaultPasswd, defaultOpensearchAddr)

	ctx := context.Background()

	err := client.CreateIndexWithDefaults(ctx, "project-test")
	assert.Nil(t, err)

	err = client.CreateDocument(ctx, "project-test", "test-id", strings.NewReader(`{"foo":"bar"}`))
	assert.Nil(t, err)

	err = client.UpdateDocument(ctx, "project-test", "test-id", strings.NewReader(`{"foo": "flv"}`))
	assert.Nil(t, err)

	err = client.DeleteDocument(ctx, "project-test", "test-id")
	assert.Nil(t, err)

	_, err = client.GetDocument(ctx, "project-test", "test-id")
	assert.NotNil(t, err)

	err = client.DeleteIndex(ctx, "project-test")
	assert.Nil(t, err)
}
