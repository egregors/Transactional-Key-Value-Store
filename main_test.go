package main

import (
	"testing"

	"github.com/egregors/Transactional-Key-Value-Store/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTKVStoreSetAndGetAValue(t *testing.T) {
	// > SET foo 123
	// > GET foo
	// 123
	s := store.NewStore()
	s.Set("foo", "123")
	r, ok := s.Get("foo")
	require.True(t, ok)
	assert.Equal(t, "123", r)
}

func TestTKVStoreDeleteAValue(t *testing.T) {
	// > DELETE foo
	// > GET foo
	// key not set
	s := store.NewStore()
	s.Set("foo", "123")
	r, ok := s.Get("foo")
	require.True(t, ok)
	assert.Equal(t, "123", r)
	s.Delete("foo")
	r, ok = s.Get("foo")
	require.False(t, ok)
	assert.Equal(t, "", r)
}
func TestTKVStoreCountTheNumberOfOccurrencesOfAValue(t *testing.T) {
	// > SET foo 123
	// > SET bar 456
	// > SET baz 123
	// > COUNT 123
	// 2
	// > COUNT 456
	// 1
	s := store.NewStore()
	s.Set("foo", "123")
	s.Set("bar", "456")
	s.Set("baz", "123")
	assert.Equal(t, 2, s.Count("123"))
	assert.Equal(t, 1, s.Count("456"))
}

func TestTKVStoreCommitATransaction(t *testing.T) {
	// > BEGIN
	// > SET foo 456
	// > COMMIT
	// > ROLLBACK
	// no transaction
	// > GET foo
	// 456
	s := store.NewStore()
	s.Begin()
	s.Set("foo", "456")
	err := s.Commit()
	require.NoError(t, err)
	r, ok := s.Get("foo")
	require.True(t, ok)
	assert.Equal(t, "456", r)
}

func TestTKVStoreRollbackATransaction(t *testing.T) {
	// > SET foo 123
	// > SET bar abc
	// > BEGIN
	// > SET foo 456
	// > GET foo
	// 456
	// > SET bar def
	// > GET bar
	// def
	// > ROLLBACK
	// > GET foo
	// 123
	// > GET bar
	// abc
	// > COMMIT
	// no transaction
	s := store.NewStore()
	s.Set("foo", "123")
	s.Set("bar", "abc")
	s.Begin()
	s.Set("foo", "456")
	r, ok := s.Get("foo")
	require.True(t, ok)
	assert.Equal(t, "456", r)
	s.Set("bar", "def")
	r, ok = s.Get("bar")
	require.True(t, ok)
	assert.Equal(t, "def", r)
	err := s.Rollback()
	require.NoError(t, err)
	r, ok = s.Get("foo")
	require.True(t, ok)
	assert.Equal(t, "123", r)
	r, ok = s.Get("bar")
	require.True(t, ok)
	assert.Equal(t, "abc", r)
	err = s.Commit()
	require.Error(t, err)
}

func TestTKVStoreNestedTransactions(t *testing.T) {
	// > SET foo 123
	// > BEGIN
	// > SET foo 456
	// > BEGIN
	// > SET foo 789
	// > GET foo
	// 789
	// > ROLLBACK
	// > GET foo
	// 456
	// > ROLLBACK
	// > GET foo
	// 123
	s := store.NewStore()
	s.Set("foo", "123")
	s.Begin()
	s.Set("foo", "456")
	s.Begin()
	s.Set("foo", "789")
	r, ok := s.Get("foo")
	require.True(t, ok)
	assert.Equal(t, "789", r)
	err := s.Rollback()
	require.NoError(t, err)
	r, ok = s.Get("foo")
	require.True(t, ok)
	assert.Equal(t, "456", r)
	err = s.Rollback()
	require.NoError(t, err)
	r, ok = s.Get("foo")
	require.True(t, ok)
	assert.Equal(t, "123", r)
}
