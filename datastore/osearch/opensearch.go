package osearch

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

const (
	DefaultMaxRetries = 5
)

type OpenSearchClient struct {
	handle *opensearch.Client
}

func NewOpenSearchClient(username, password string, urls ...string) (*OpenSearchClient, error) {
	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses:     urls,
		MaxRetries:    DefaultMaxRetries,
		Username:      username,
		Password:      password,
		RetryOnStatus: []int{502, 503, 504},
	})

	if err != nil {
		return nil, err
	}

	return &OpenSearchClient{handle: client}, nil
}

func (client *OpenSearchClient) CreateDocument(ctx context.Context, index, id string, obj interface{}) error {
	body, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	req := opensearchapi.CreateRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewBuffer(body),
	}
	resp, err := req.Do(ctx, client.handle)
	if err != nil {
		return err
	}

	if resp.StatusCode > 299 {
		return errors.New(resp.String())
	}

	return nil
}

func (client *OpenSearchClient) UpdateDocument(ctx context.Context, index, id string, obj interface{}) error {
	val, err := json.Marshal(map[string]interface{}{
		"doc": obj,
	})
	if err != nil {
		return err
	}

	req := opensearchapi.UpdateRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewBuffer(val),
	}
	resp, err := req.Do(ctx, client.handle)
	if err != nil {
		return err
	}

	if resp.StatusCode > 299 {
		return errors.New(resp.String())
	}
	return nil
}

func (client *OpenSearchClient) DeleteDocument(ctx context.Context, index, id string) error {
	req := opensearchapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
	}

	resp, err := req.Do(ctx, client.handle)
	if err != nil {
		return err
	}

	if resp.StatusCode > 299 {
		return errors.New(resp.String())
	}
	return nil
}

func (client *OpenSearchClient) GetDocument(ctx context.Context, index, id string) (io.ReadCloser, error) {
	req := opensearchapi.GetRequest{
		Index:      index,
		DocumentID: id,
	}
	resp, err := req.Do(ctx, client.handle)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		return nil, errors.New(resp.String())
	}

	return resp.Body, nil
}

func (client *OpenSearchClient) CreateIndexWithDefaults(ctx context.Context, name string) error {
	settings := strings.NewReader(`{
		"settings": {
			"index": {
				"number_of_shards": 1,
				"number_of_replicas": 0
				}
			}
		}`)
	return client.CreateIndexWithSettings(ctx, name, settings)
}

func (client *OpenSearchClient) CreateIndexWithSettings(ctx context.Context, name string, settings io.Reader) error {
	req := opensearchapi.IndicesCreateRequest{
		Index: name,
		Body:  settings,
	}
	resp, err := req.Do(ctx, client.handle)
	if err != nil {
		return err
	}

	if resp.StatusCode > 299 {
		return errors.New(resp.Status())
	}

	return nil
}

func (client *OpenSearchClient) DeleteIndex(ctx context.Context, name string) error {
	req := opensearchapi.IndicesDeleteRequest{
		Index: []string{name},
	}

	resp, err := req.Do(ctx, client.handle)
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		return errors.New(resp.Status())
	}

	return nil
}
