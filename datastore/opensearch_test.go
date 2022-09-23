package datastore

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Create_Delete_Index(t *testing.T) {

	ctx := context.Background()

	err := CreateIndexWithDefaults(ctx, "project-test")
	assert.Nil(t, err)

	err = CreateDocument(ctx, "project-test", "test-id", strings.NewReader(`{"foo":"bar"}`))
	assert.Nil(t, err)

	err = UpdateDocument(ctx, "project-test", "test-id", strings.NewReader(`{"foo": "flv"}`))
	assert.Nil(t, err)

	err = DeleteDocument(ctx, "project-test", "test-id")
	assert.Nil(t, err)

	_, err = GetDocument(ctx, "project-test", "test-id")
	assert.NotNil(t, err)

	err = DeleteIndex(ctx, "project-test")
	assert.Nil(t, err)
}
