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

// ----------------------------------------------------
// TrackInsertOne
// ----------------------------------------------------

func TestTrackClient_InsertOne(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/track/insert-one", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "123", r.Header.Get("X-App-Id"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	c := util.NewFirbolgClient(server.URL+"/track", 123)
	client := TrackClient{FirbolgClient: c}

	err := client.TrackInsertOne(11, "foo", 42, nil, nil)
	assert.NoError(t, err)
}

// ----------------------------------------------------
// TrackGetMany
// ----------------------------------------------------

func TestTrackClient_GetMany(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/track/get-many", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(loadTestData(t, "track_get_many_response.json"))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	c := util.NewFirbolgClient(server.URL+"/track", 123)
	client := TrackClient{FirbolgClient: c}

	req := bullet_interface.TrackGetManyRequest{
		Buckets: []bullet_interface.TrackGetKeys{
			{BucketID: 1, Keys: []string{"a", "b"}},
		},
	}

	resp, err := client.TrackGetMany(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, int64(100), resp.Values[1]["a"].Value)
	assert.Equal(t, int64(200), resp.Values[1]["b"].Value)
}

// ----------------------------------------------------
// TrackGetManyByPrefix
// ----------------------------------------------------

func TestTrackClient_GetManyByPrefix(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/track/get-query", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(loadTestData(t, "track_prefix_response.json"))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	c := util.NewFirbolgClient(server.URL+"/track", 123)
	client := TrackClient{FirbolgClient: c}

	req := bullet_interface.TrackGetItemsByPrefixRequest{
		BucketID: 1,
		Prefix:   "foo",
		Tags:     []int64{10},
		Metric: &bullet_interface.MetricFilter{
			Operator: "gt",
			Value:    3.14,
		},
	}

	resp, err := client.TrackGetManyByPrefix(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	item1 := resp.Values[1]["foo123"]
	assert.NotNil(t, item1.Metric)
	assert.Equal(t, float64(9.9), *item1.Metric)
}

// ----------------------------------------------------
// TrackDeleteMany
// ----------------------------------------------------

func TestTrackClient_DeleteMany(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/track/delete-many", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "123", r.Header.Get("X-App-Id"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	c := util.NewFirbolgClient(server.URL+"/track", 123)
	client := TrackClient{FirbolgClient: c}

	req := bullet_interface.TrackDeleteMany{
		Values: []bullet_interface.TrackDeleteValue{
			{BucketID: 1, Key: "foo"},
			{BucketID: 1, Key: "bar"},
		},
	}

	err := client.TrackDeleteMany(req)
	assert.NoError(t, err)
}
