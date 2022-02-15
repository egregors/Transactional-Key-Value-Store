//nolint:revive // it's ok
package store

import "fmt"

type Store struct {
	db    DBMap
	stack *TKVStoreStack
}

func NewStore() *Store {
	return &Store{
		db:    make(DBMap),
		stack: NewTKVStoreStack(),
	}
}

func (s *Store) Debug() {
	fmt.Println(s.db)
}

func (s *Store) Set(k, v string) {
	db := *s.getTx()
	db[k] = v
}

func (s *Store) Get(k string) (string, bool) {
	db := *s.getTx()
	if v, ok := db[k]; ok {
		return v, true
	}
	return "", false
}

func (s *Store) Delete(k string) {
	db := *s.getTx()
	delete(db, k)
}

func (s *Store) Count(v string) int {
	db := *s.getTx()
	counter := 0
	for _, v2 := range db {
		if v2 == v {
			counter++
		}
	}
	return counter
}

func (s *Store) Begin() {
	if s.stack.IsEmpty() {
		s.stack.Push(s.db.Copy())
	} else {
		s.stack.Push(s.stack.GetLast().Copy())
	}
}

func (s *Store) Commit() error {
	if s.stack.IsEmpty() {
		return fmt.Errorf("no transaction")
	}
	db := s.stack.Pop()
	if s.stack.IsEmpty() {
		s.db = *db
	} else {
		s.stack.UpdateLast(db)
	}
	return nil
}

func (s *Store) Rollback() error {
	if s.stack.IsEmpty() {
		return fmt.Errorf("no transaction")
	}
	s.stack.Pop()
	return nil
}

func (s *Store) getTx() *DBMap {
	if s.stack.IsEmpty() {
		return &s.db
	}
	return s.stack.GetLast()
}
