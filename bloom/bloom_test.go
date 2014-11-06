package bloom

import (
	"code.google.com/p/go-uuid/uuid"
	"testing"
)

func TestBloom2(t *testing.T) {

	filter := NewSimpleBloomFilter(2)

	uuid1 := uuid.New()
	filter.Add(uuid1)
	if !filter.Test(uuid1) {
		t.Fatalf("%s should exist in filter", uuid1)
	}

	uuid2 := uuid.New()
	filter.Add(uuid2)
	if !filter.Test(uuid2) {
		t.Fatalf("%s should exist in filter", uuid2)
	}

	uuid3 := uuid.New()
	if filter.Test(uuid3) {
		t.Fatalf("%s should not exist in filter", uuid3)
	}

	if filter.Count() != 2 {
		t.Fatalf("filter should contain %d records", 2)
	}
}
