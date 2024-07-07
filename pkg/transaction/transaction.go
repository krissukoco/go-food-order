package transaction

import (
	"context"
	"database/sql"
	"log"
)

type RunnerFunc = func(ctx context.Context) error

type Transactioner interface {
	WithTx(ctx context.Context, fn RunnerFunc) error
}

type sqlController struct {
	db *sql.DB
}

var _ Transactioner = (*sqlController)(nil)

func (ctl *sqlController) WithTx(ctx context.Context, fn RunnerFunc) error {
	var err error
	// Begin transaction
	tx, err := ctl.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Build context
	c := buildContext(ctx, tx)

	// Handle errors and panics on defer
	defer SafelyCommit(tx, err)

	err = fn(c)
	return err
}

type contextKey struct {
	key string
}

var (
	sqlCtxKey = &contextKey{"database/sql"}
)

func buildContext(parent context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(parent, sqlCtxKey, tx)
}

func FromContext(ctx context.Context) (*sql.Tx, bool) {
	v := ctx.Value(sqlCtxKey)
	tx, ok := v.(*sql.Tx)
	return tx, ok
}

func SafelyCommit(tx *sql.Tx, err error) error {
	if r := recover(); r != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("ROLLBACK ERROR: %v\n", rbErr)
		}
		// re-throw panic
		panic(r)
	}
	// Rollback if error is not nil
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("ROLLBACK ERROR: %v\n", rbErr)
		}
		return err
	}
	// Finally commit transaction
	return tx.Commit()
}
