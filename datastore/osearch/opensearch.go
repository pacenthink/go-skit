package osearch

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

const (
	DefaultMaxRetries = 5
)

// These are all non-prod values
const (
	defaultOpensearchAddr = "https://127.0.0.1:9200"
	defaultUsername       = "admin"
	defaultPasswd         = "admin"
)

type OpenSearchClient struct {
	handle *opensearch.Client
}

func DefaultClient() (*OpenSearchClient, error) {
	username := os.Getenv("OPENSEARCH_USERNAME")
	if username == "" {
		username = defaultUsername
	}
	password := os.Getenv("OPENSEARCH_SECRET")
	if password == "" {
		password = defaultPasswd
	}

	urls := getUrls()

	log.Printf("INF OpenSearch urls: %v", urls)
	return NewOpenSearchClient(username, password, urls...)
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

func (client *OpenSearchClient) Raw() *opensearch.Client {
	return client.handle
}

func getUrls() []string {
	urlstr := os.Getenv("OPENSEARCH_URLS")
	if urlstr == "" {
		return []string{defaultOpensearchAddr}
	}

	urls := strings.Split(urlstr, ",")
	out := make([]string, 0, len(urls))
	for _, u := range urls {
		cu := strings.TrimSpace(u)
		if cu != "" {
			out = append(out, cu)
		}
	}
	if len(out) == 0 {
		return []string{defaultOpensearchAddr}
	}
	return out
}
