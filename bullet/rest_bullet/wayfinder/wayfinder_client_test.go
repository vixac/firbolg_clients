package rest_bullet

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	bullet_interface "github.com/vixac/firbolg_clients/bullet/bullet_interface"
	util "github.com/vixac/firbolg_clients/bullet/util"
)

func loadTestData(t *testing.T, filename string) []byte {
	t.Helper()
	data, err := os.ReadFile("testdata/" + filename)
	if err != nil {
		t.Fatalf("failed to read test data file '%s': %v", filename, err)
	}
	return data
}

func TestWayFinderClient_WayFinderInsertOne(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/wayfinder/insert-one", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(loadTestData(t, "insert_one_response.json"))
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	c := util.NewFirbolgClient(server.URL, 123)
	client := WayFinderClient{
		Client: c,
	}
	id, err := client.WayFinderInsertOne(bullet_interface.WayFinderPutRequest{
		BucketId: 1,
		Key:      "foo",
		Payload:  "bar",
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(1234), id)
}

func TestWayFinderClient_WayFinderQueryByPrefix(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/wayfinder/query-by-prefix", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(loadTestData(t, "query_by_prefix_response.json"))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	c := util.NewFirbolgClient(server.URL, 123)
	client := WayFinderClient{
		Client: c,
	}

	items, err := client.WayFinderQueryByPrefix(bullet_interface.WayFinderPrefixQueryRequest{
		BucketId:   1,
		Prefix:     "foo",
		MetricIsGt: false,
	})
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "foo", items[0].Key)
	assert.Equal(t, int64(42), items[0].ItemId)
	assert.Equal(t, "some payload", items[0].Payload)
}

func TestWayFinderClient_WayFinderGetOne(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/wayfinder/get-one", func(w http.ResponseWriter, r *http.Request) {
		// Verify method and content-type if you like here too
		w.WriteHeader(http.StatusOK)
		w.Write(loadTestData(t, "get_one_response.json"))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	c := util.NewFirbolgClient(server.URL, 123)
	client := WayFinderClient{
		Client: c,
	}

	item, err := client.WayFinderGetOne(bullet_interface.WayFinderGetOneRequest{
		BucketId: 1,
		Key:      "foo",
	})
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, "bar", item.Payload)
	assert.NotNil(t, item.Tag)
	assert.Equal(t, *item.Tag, int64(52))

	assert.NotNil(t, item.Metric)
	assert.Equal(t, *item.Metric, 12.3)

}
