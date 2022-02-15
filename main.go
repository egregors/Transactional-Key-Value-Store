//nolint:revive // it's ok
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type DBMap map[string]string

func (db *DBMap) Copy() *DBMap {
	cp := make(DBMap, len(*db))
	for k, v := range *db {
		cp[k] = v
	}
	return &cp
}

type TKVStoreStack struct {
	dbs []*DBMap
}

func NewTKVStoreStack() *TKVStoreStack {
	return &TKVStoreStack{
		dbs: []*DBMap{},
	}
}

func (st *TKVStoreStack) IsEmpty() bool     { return len(st.dbs) == 0 }
func (st *TKVStoreStack) Push(store *DBMap) { st.dbs = append(st.dbs, store) }
func (st *TKVStoreStack) Pop() *DBMap {
	store := st.GetLast()
	if store != nil {
		st.dbs = st.dbs[:len(st.dbs)-1]
		return store
	}
	return nil
}

func (st *TKVStoreStack) GetLast() *DBMap {
	if !st.IsEmpty() {
		store := st.dbs[len(st.dbs)-1]
		return store
	}
	return nil
}

func (st *TKVStoreStack) UpdateLast(db *DBMap) {
	if !st.IsEmpty() {
		st.dbs[len(st.dbs)-1] = db
	}
}

type TransactionalKVStorer interface {
	Set(k, v string)
	Get(k string) (string, bool)
	Delete(k string)
	Count(v string) int
	Begin()
	Commit() error
	Rollback() error

	Debug()
}

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
	db := *s.getCurrentDB()
	db[k] = v
}

func (s *Store) Get(k string) (string, bool) {
	db := *s.getCurrentDB()
	if v, ok := db[k]; ok {
		return v, true
	}
	return "", false
}

func (s *Store) Delete(k string) {
	db := *s.getCurrentDB()
	delete(db, k)
}

func (s *Store) Count(v string) int {
	db := *s.getCurrentDB()
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

func (s *Store) getCurrentDB() *DBMap {
	if s.stack.IsEmpty() {
		return &s.db
	}
	return s.stack.GetLast()
}

func main() {
	fmt.Println("TKVS :3")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	store := NewStore()
	go runStoreAndCLI(store)

	<-stop
	fmt.Println("shutting down...")
}

func runStoreAndCLI(store TransactionalKVStorer) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		rawCmd, _ := reader.ReadString('\n')
		rawCmd = strings.TrimSuffix(rawCmd, "\n")
		cmds := strings.Split(rawCmd, " ")

		switch cmds[0] {
		case "SET":
			store.Set(cmds[1], cmds[2])
		case "GET":
			if v, ok := store.Get(cmds[1]); ok {
				fmt.Printf("%s\n", v)
			} else {
				fmt.Printf("key not set\n")
			}
		case "DELETE":
			store.Delete(cmds[1])
		case "COUNT":
			fmt.Printf("%d\n", store.Count(cmds[1]))
		case "BEGIN":
			store.Begin()
		case "COMMIT":
			err := store.Commit()
			if err != nil {
				fmt.Println(err.Error())
			}
		case "ROLLBACK":
			err := store.Rollback()
			if err != nil {
				fmt.Println(err.Error())
			}
		default:
			fmt.Printf(
				"wrong command, expecting: SET, GET, DELETE, COUNT, BEGIN, COMMIT, ROLLBACK; got: %s\n",
				rawCmd,
			)
		}
	}
}
