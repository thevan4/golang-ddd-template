package application

import (
	"github.com/sirupsen/logrus"
)

// TaskReceiver gets tasks for work
type TaskReceiver struct {
	GracefulShutdownListen chan struct{}
	TaskReceiverIsDone     chan struct{}
	Facade                 *Facade
}

// NewTaskReceiver ...
func NewTaskReceiver(facade *Facade, gracefulShutdownListen chan struct{}, taskReceiverIsDone chan struct{}) *TaskReceiver {
	return &TaskReceiver{
		GracefulShutdownListen: gracefulShutdownListen,
		TaskReceiverIsDone:     taskReceiverIsDone,
		Facade:                 facade,
	}
}

// DoJob example of a task receiver
func (taskReceiver *TaskReceiver) DoJob(jobType string) {
	receivingAndProcessingUUID := taskReceiver.Facade.UUIDgenerator.NewUUID().UUID.String()
	taskReceiver.Facade.Logging.WithFields(logrus.Fields{"event uuid": receivingAndProcessingUUID}).Info("got new do job request")
	// try unmarshall request
	// validate request
	if err := taskReceiver.Facade.DoJob(jobType, receivingAndProcessingUUID); err != nil {
		taskReceiver.Facade.Logging.WithFields(logrus.Fields{"event uuid": receivingAndProcessingUUID}).Errorf("do job request error: %v", err)
		// response error
	}
	taskReceiver.Facade.Logging.WithFields(logrus.Fields{"event uuid": receivingAndProcessingUUID}).Info("do job request done")
	// response okay
}
