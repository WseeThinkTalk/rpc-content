package snowflake

import (
	"sync"
	"testing"
)

func TestSnowflake_GenerateID_Unique(t *testing.T) {
	node, err := NewNode(1)
	if err != nil {
		t.Fatalf("failed to create node: %v", err)
	}

	const count = 10000
	var wg sync.WaitGroup
	var idMap sync.Map

	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			id := node.Generate()
			if _, loaded := idMap.LoadOrStore(id, true); loaded {
				t.Errorf("duplicate id generated: %d", id)
			}
		}()
	}
	wg.Wait()
}

func TestSnowflake_GenerateID_Monotonic(t *testing.T) {
	node, err := NewNode(1)
	if err != nil {
		t.Fatalf("failed to create node: %v", err)
	}

	var lastId int64
	for i := 0; i < 1000; i++ {
		id := node.Generate()
		if id <= lastId {
			t.Fatalf("id not monotonic: current %d <= last %d", id, lastId)
		}
		lastId = id
	}
}

func TestGenerateID_Default(t *testing.T) {
	id := GenerateID()
	if id <= 0 {
		t.Fatalf("expected positive id, got %d", id)
	}
}
