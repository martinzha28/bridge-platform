package ws

import (
	"sync"
	"testing"
)

func TestHubCreateTable(t *testing.T) {
	hub := NewHub()
	table := hub.CreateTable()

	if table == nil {
		t.Fatal("CreateTable returned nil")
	}
	if table.ID == "" {
		t.Error("table ID should not be empty")
	}
}

func TestHubGetTable(t *testing.T) {
	hub := NewHub()
	table := hub.CreateTable()

	got, ok := hub.GetTable(table.ID)
	if !ok {
		t.Fatal("GetTable should find the created table")
	}
	if got != table {
		t.Error("GetTable returned a different table")
	}
}

func TestHubGetTableNotFound(t *testing.T) {
	hub := NewHub()

	if _, ok := hub.GetTable("nonexistent"); ok {
		t.Error("GetTable should return false for unknown ID")
	}
}

func TestHubRemoveTable(t *testing.T) {
	hub := NewHub()
	table := hub.CreateTable()

	hub.RemoveTable(table.ID)

	if _, ok := hub.GetTable(table.ID); ok {
		t.Error("table should be gone after RemoveTable")
	}
}

func TestHubUniqueIDs(t *testing.T) {
	hub := NewHub()
	seen := make(map[string]bool)

	for range 100 {
		table := hub.CreateTable()
		if seen[table.ID] {
			t.Fatalf("duplicate table ID: %s", table.ID)
		}
		seen[table.ID] = true
	}
}

func TestHubRemoveNonexistent(t *testing.T) {
	hub := NewHub()
	hub.RemoveTable("does-not-exist") // should not panic
}

func TestHubConcurrentAccess(t *testing.T) {
	hub := NewHub()
	var wg sync.WaitGroup

	// Hammer create/get/remove from 20 goroutines
	for range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			table := hub.CreateTable()
			hub.GetTable(table.ID)
			hub.RemoveTable(table.ID)
		}()
	}
	wg.Wait()
}
