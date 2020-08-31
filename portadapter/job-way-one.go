package portadapter

import (
	"time"

	"github.com/sirupsen/logrus"
)

// JobWayOne ...
type JobWayOne struct {
	logging *logrus.Logger
}

// NewJobWayOne ...
func NewJobWayOne(logging *logrus.Logger) *JobWayOne {
	return &JobWayOne{logging: logging}
}

// WorkHard ...
func (jobWayOne *JobWayOne) WorkHard(uuid string) error {
	jobWayOne.logging.WithFields(logrus.Fields{"event uuid": uuid}).Info("doing job one")
	time.Sleep(3 * time.Second)
	jobWayOne.logging.WithFields(logrus.Fields{"event uuid": uuid}).Info("job one done")
	return nil
}
