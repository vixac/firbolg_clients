package bullet_stl

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/vixac/bullet/client"
	"github.com/vixac/bullet/model"
)

// AtomicMesh implements Mesh using one atomic Track operation per Mesh method.
//
// Each relationship is still indexed twice: once by subject and once by object.
// TrackMutate keeps those two index entries in sync. The multi-subject and
// multi-object reads use one multi-prefix query, so their results also come from
// one Track snapshot.
//
// Atomicity is per Mesh call. A caller that reads pairs and subsequently passes
// them to RemovePairs is performing two independent atomic operations.
type AtomicMesh struct {
	TrackStore        client.Track
	BucketId          int32
	MeshName          string
	ForwardSeparator  string
	BackwardSeparator string
}

var _ Mesh = (*AtomicMesh)(nil)

func NewAtomicMesh(store client.Track, bucketID int32, meshName, forwardSeparator, backwardSeparator string) (Mesh, error) {
	return &AtomicMesh{
		TrackStore:        store,
		BucketId:          bucketID,
		MeshName:          meshName,
		ForwardSeparator:  forwardSeparator,
		BackwardSeparator: backwardSeparator,
	}, nil
}

func (m *AtomicMesh) AppendPairs(pairs []ManyToManyPair) error {
	if len(pairs) == 0 {
		return nil
	}

	mutationID, err := atomicMeshMutationID()
	if err != nil {
		return err
	}
	mutation := model.TrackMutation{MutationID: mutationID, Puts: make([]model.TrackPut, 0, len(pairs)*2)}
	for _, pair := range pairs {
		object := pair.Object.Value
		rank := float64(pair.Rank)
		mutation.Puts = append(mutation.Puts,
			model.TrackPut{BucketID: m.BucketId, Key: buildKey(m.MeshName, m.ForwardSeparator, pair.Subject.Value, &object, false), Metric: &rank},
			model.TrackPut{BucketID: m.BucketId, Key: buildKey(m.MeshName, m.BackwardSeparator, pair.Object.Value, &pair.Subject.Value, false), Metric: &rank},
		)
	}
	return m.applyMutation(mutation)
}

func (m *AtomicMesh) RemovePairs(pairs []ManyToManyPair) error {
	if len(pairs) == 0 {
		return nil
	}

	mutationID, err := atomicMeshMutationID()
	if err != nil {
		return err
	}
	mutation := model.TrackMutation{MutationID: mutationID, Deletes: make([]model.TrackKey, 0, len(pairs)*2)}
	for _, pair := range pairs {
		object := pair.Object.Value
		mutation.Deletes = append(mutation.Deletes,
			model.TrackKey{BucketID: m.BucketId, Key: buildKey(m.MeshName, m.ForwardSeparator, pair.Subject.Value, &object, false)},
			model.TrackKey{BucketID: m.BucketId, Key: buildKey(m.MeshName, m.BackwardSeparator, pair.Object.Value, &pair.Subject.Value, false)},
		)
	}
	return m.applyMutation(mutation)
}

func (m *AtomicMesh) AllPairsForManySubjects(subjects []ListSubject) (*PairFetchResponse, error) {
	prefixes := make([]string, 0, len(subjects))
	for _, subject := range subjects {
		prefixes = append(prefixes, buildKey(m.MeshName, m.ForwardSeparator, subject.Value, nil, false))
	}
	items, err := m.TrackStore.GetItemsByKeyPrefixes(m.BucketId, prefixes, nil, nil, false)
	if err != nil {
		return nil, err
	}
	pairs, err := m.forwardPairs(items)
	if err != nil || len(pairs) == 0 {
		return nil, err
	}
	return &PairFetchResponse{Pairs: pairs}, nil
}

func (m *AtomicMesh) AllPairsForSubject(subject ListSubject) (*PairFetchResponse, error) {
	return m.allPairsForSubject(subject, false)
}

func (m *AtomicMesh) AllPairsForPrefixSubject(subject ListSubject) (*PairFetchResponse, error) {
	return m.allPairsForSubject(subject, true)
}

func (m *AtomicMesh) AllPairsForObject(object ListObject) (*PairFetchResponse, error) {
	return m.AllPairsForManyObjects([]ListObject{object})
}

func (m *AtomicMesh) AllPairsForManyObjects(objects []ListObject) (*PairFetchResponse, error) {
	prefixes := make([]string, 0, len(objects))
	for _, object := range objects {
		prefixes = append(prefixes, buildKey(m.MeshName, m.BackwardSeparator, object.Value, nil, false))
	}
	items, err := m.TrackStore.GetItemsByKeyPrefixes(m.BucketId, prefixes, nil, nil, false)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}

	pairs := make([]ManyToManyPair, 0, len(items))
	for _, item := range items {
		parts := strings.Split(item.Key, m.BackwardSeparator)
		if len(parts) != 3 || parts[0] != m.MeshName {
			return nil, errors.New("expected <meshname><separator><object><separator><subject>")
		}
		pairs = append(pairs, ManyToManyPair{
			Subject: ListSubject{Value: parts[2]},
			Object:  ListObject{Value: parts[1]},
			Rank:    metricRank(item.Value.Metric),
		})
	}
	sortPairs(pairs)
	return &PairFetchResponse{Pairs: pairs}, nil
}

func (m *AtomicMesh) applyMutation(mutation model.TrackMutation) error {
	result, err := m.TrackStore.TrackMutate(mutation)
	if err != nil {
		return err
	}
	if !result.Applied {
		return errors.New("atomic mesh mutation ID was already applied")
	}
	return nil
}

func (m *AtomicMesh) allPairsForSubject(subject ListSubject, subjectIsPrefix bool) (*PairFetchResponse, error) {
	prefix := buildKey(m.MeshName, m.ForwardSeparator, subject.Value, nil, subjectIsPrefix)
	items, err := m.TrackStore.GetItemsByKeyPrefix(m.BucketId, prefix, nil, nil, false)
	if err != nil {
		return nil, err
	}
	pairs, err := m.forwardPairs(items)
	if err != nil || len(pairs) == 0 {
		return nil, err
	}
	return &PairFetchResponse{Pairs: pairs}, nil
}

func (m *AtomicMesh) forwardPairs(items []model.TrackKeyValueItem) ([]ManyToManyPair, error) {
	if len(items) == 0 {
		return nil, nil
	}
	pairs := make([]ManyToManyPair, 0, len(items))
	for _, item := range items {
		parts := strings.Split(item.Key, m.ForwardSeparator)
		if len(parts) != 3 || parts[0] != m.MeshName {
			return nil, errors.New("expected <meshname><separator><subject><separator><object>")
		}
		pairs = append(pairs, ManyToManyPair{
			Subject: ListSubject{Value: parts[1]},
			Object:  ListObject{Value: parts[2]},
			Rank:    metricRank(item.Value.Metric),
		})
	}
	sortPairs(pairs)
	return pairs, nil
}

func atomicMeshMutationID() (model.MutationID, error) {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", fmt.Errorf("generate atomic mesh mutation ID: %w", err)
	}
	return model.MutationID("atomic-mesh-" + hex.EncodeToString(bytes[:])), nil
}

func metricRank(metric *float64) int32 {
	if metric == nil {
		return 0
	}
	return int32(*metric)
}

func sortPairs(pairs []ManyToManyPair) {
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Subject.Value == pairs[j].Subject.Value {
			return pairs[i].Object.Value < pairs[j].Object.Value
		}
		return pairs[i].Subject.Value < pairs[j].Subject.Value
	})
}
