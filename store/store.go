//nolint:revive // it's ok
package store

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
