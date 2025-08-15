package ship

import "errors"

var (
	// ErrExecutorNotSet is returned when a container tool has no executor function
	ErrExecutorNotSet = errors.New("tool executor function is not set")

	// ErrToolNotFound is returned when a requested tool is not found in the registry
	ErrToolNotFound = errors.New("tool not found")

	// ErrInvalidParameter is returned when a parameter validation fails
	ErrInvalidParameter = errors.New("invalid parameter")

	// ErrDaggerEngineNotAvailable is returned when Dagger engine is not available
	ErrDaggerEngineNotAvailable = errors.New("dagger engine not available")
)
