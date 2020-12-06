package types

import "context"

const (
	// StatePending means CRD instance is created; Pod info has been updated into CRD instance;
	// Pod has been accepted by the system, but one or more of the containers has not been started.
	StatePending string = "Pending"
	// StateRunning means Pod has been bound to a node and all of the containers have been started.
	StateRunning string = "Running"
	// StateSucceeded means that all containers in the Pod have voluntarily terminated with a container
	// exit code of 0, and the system is not going to restart any of these containers.
	StateSucceeded string = "Succeeded"
	// StateFailed means that all containers in the Pod have terminated, and at least one container has
	// terminated in a failure (exited with a non-zero exit code or was stopped by the system).
	StateFailed string = "Failed"
)

var (
	ctx context.Context
)

func init() {
	ctx = context.Background()
}

// GetCtx returns global default context.
func GetCtx() context.Context {
	return ctx
}
