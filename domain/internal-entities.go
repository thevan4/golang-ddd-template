package domain

import "sync"

// GracefulShutdown used for a graceful finish usecases
type GracefulShutdown struct {
	sync.Mutex
	ShutdownNow  bool
	UsecasesJobs int
}
