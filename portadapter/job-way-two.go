package portadapter

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// JobWayTwo ...
type JobWayTwo struct {
	mux     sync.Mutex // for disable parallel execution
	logging *logrus.Logger
}

// NewJobWayTwo  ...
func NewJobWayTwo(logging *logrus.Logger) *JobWayTwo {
	return &JobWayTwo{logging: logging}
}

// WorkHard ...
func (jobWayTwo *JobWayTwo) WorkHard(uuid string) error {
	jobWayTwo.mux.Lock()
	defer jobWayTwo.mux.Unlock()
	jobWayTwo.logging.WithFields(logrus.Fields{"event uuid": uuid}).Info("doing job two")
	time.Sleep(3 * time.Second)
	jobWayTwo.logging.WithFields(logrus.Fields{"event uuid": uuid}).Info("job two done")
	return nil
}
