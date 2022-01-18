package dispatcher

import "context"

// Resource of worker own.
type Resource interface {
	Check(ctx context.Context) (Resource, error)
	Release() error
}
