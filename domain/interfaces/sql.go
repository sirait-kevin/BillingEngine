package interfaces

import (
	"context"
	"database/sql"
)

//go:generate mockgen -build_flags=-mod=mod -destination ../../mocks/domain/atomic_transaction.go -package=mock_domain github.com/sirait-kevin/BillingEngine/domain/interfaces AtomicTransaction
type AtomicTransaction interface {
	Rollback() error
	Commit() error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}
