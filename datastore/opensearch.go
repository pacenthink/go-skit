package datastore

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

var OpenSearchClient *opensearch.Client

func NewOpenSearchClient(urls ...string) (*opensearch.Client, error) {
	if len(urls) == 0 {
		urls = []string{"https://127.0.0.1:9200"}
	}
	return opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses:     urls,
		MaxRetries:    5,
		Username:      "admin",
		Password:      "admin",
		RetryOnStatus: []int{502, 503, 504},
	})
}

func CreateDocument(ctx context.Context, index, id string, obj interface{}) error {
	body, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	req := opensearchapi.CreateRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewBuffer(body),
	}
	resp, err := req.Do(ctx, OpenSearchClient)
	if err != nil {
		return err
	}

	if resp.StatusCode > 299 {
		return errors.New(resp.String())
	}

	return nil
}

func UpdateDocument(ctx context.Context, index, id string, obj interface{}) error {
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
	resp, err := req.Do(ctx, OpenSearchClient)
	if err != nil {
		return err
	}

	if resp.StatusCode > 299 {
		return errors.New(resp.String())
	}
	return nil
}

func DeleteDocument(ctx context.Context, index, id string) error {
	req := opensearchapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
	}

	resp, err := req.Do(ctx, OpenSearchClient)
	if err != nil {
		return err
	}

	if resp.StatusCode > 299 {
		return errors.New(resp.String())
	}
	return nil
}

func GetDocument(ctx context.Context, index, id string) (io.ReadCloser, error) {
	req := opensearchapi.GetRequest{
		Index:      index,
		DocumentID: id,
	}
	resp, err := req.Do(ctx, OpenSearchClient)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		return nil, errors.New(resp.String())
	}

	return resp.Body, nil
}

func CreateIndexWithDefaults(ctx context.Context, name string) error {
	settings := strings.NewReader(`{
		"settings": {
			"index": {
				"number_of_shards": 1,
				"number_of_replicas": 0
				}
			}
		}`)
	return CreateIndexWithSettings(ctx, name, settings)
}

func CreateIndexWithSettings(ctx context.Context, name string, settings io.Reader) error {
	req := opensearchapi.IndicesCreateRequest{
		Index: name,
		Body:  settings,
	}
	resp, err := req.Do(ctx, OpenSearchClient)
	if err != nil {
		return err
	}

	log.Println(resp.String())
	if resp.StatusCode > 299 {
		return errors.New(resp.Status())
	}

	return nil
}

func DeleteIndex(ctx context.Context, name string) error {
	req := opensearchapi.IndicesDeleteRequest{
		Index: []string{name},
	}

	resp, err := req.Do(ctx, OpenSearchClient)
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		return errors.New(resp.Status())
	}

	return nil
}
