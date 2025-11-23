package bullet_stl

import "fmt"

type NamespacedId struct {
	Namespace int32
	Id        int32
}

/**
The motivation here is to use a single int64 field to hold the namespace in the first 32 bits, and the value in the last 32 bits.
*/
// MakeID combines namespaceId (int32) and localID (int32) into an int64.
func MakeNamespacedId(namespace int32, localID int32) int64 {
	return (int64(namespace) << 32) | (int64(localID) & 0xffffffff)
}

// ParseID splits a combined int64 ID into bucketID and localID.
func ParseNamespacedId(namespacedId int64) NamespacedId {
	namespace := int32(namespacedId >> 32)
	id := int32(namespacedId & 0xffffffff)
	return NamespacedId{
		Namespace: namespace,
		Id:        id,
	}
}

// IsValid checks whether the ID was generated using MakeID.
// (Useful if you later add reserved bits.)
func IsValid(id int64) bool {
	// Currently all bit patterns are valid.
	// This exists in case you add sanity checks later.
	return true
}

// NextLocalID safely increments a local counter and
// protects against int32 rollover.
func NextLocalID(curr int32) (int32, error) {
	if curr == 0x7fffffff {
		return 0, fmt.Errorf("localID overflow for bucket")
	}
	return curr + 1, nil
}

// FormatID renders an ID nicely for logging/debugging.
func FormatID(id int64) string {
	spacedId := ParseNamespacedId(id)
	return fmt.Sprintf("DepotID(%d:%d)", spacedId.Namespace, spacedId.Id)
}
