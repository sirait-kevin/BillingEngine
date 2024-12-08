package interfaces

import "time"

//go:generate mockgen -build_flags=-mod=mod -destination ../../mocks/domain/clock.go -package=mock_domain github.com/sirait-kevin/BillingEngine/domain/interfaces Clock
type Clock interface {
	Now() time.Time
}
