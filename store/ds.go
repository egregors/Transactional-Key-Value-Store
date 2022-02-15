//nolint:revive // it's ok
package store

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
