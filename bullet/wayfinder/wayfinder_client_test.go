package bullet

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	util "github.com/vixac/firbolg_clients/util"
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
	id, err := client.WayFinderInsertOne(WayFinderPutRequest{
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

	items, err := client.WayFinderQueryByPrefix(WayFinderPrefixQueryRequest{
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
