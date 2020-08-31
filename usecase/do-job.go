package usecase

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/thevan4/golang-ddd-template/domain"
)

// DoJobEntity ...
type DoJobEntity struct {
	jobWorker        domain.Work
	gracefulShutdown *domain.GracefulShutdown
	logging          *logrus.Logger
}

// NewDoJobEntity ...
func NewDoJobEntity(jobWorker domain.Work,
	gracefulShutdown *domain.GracefulShutdown,
	logging *logrus.Logger) *DoJobEntity {
	return &DoJobEntity{
		jobWorker:        jobWorker,
		gracefulShutdown: gracefulShutdown,
		logging:          logging,
	}
}

// DoJob ...
func (doJobEntity *DoJobEntity) DoJob(uuid string) error {
	// graceful shutdown part start
	doJobEntity.gracefulShutdown.Lock()
	if doJobEntity.gracefulShutdown.ShutdownNow {
		defer doJobEntity.gracefulShutdown.Unlock()
		return fmt.Errorf("program got shutdown signal, usecase do job canceled")
	}
	doJobEntity.gracefulShutdown.UsecasesJobs++
	doJobEntity.gracefulShutdown.Unlock()
	defer decreaseJobs(doJobEntity.gracefulShutdown)
	// graceful shutdown part end

	doJobEntity.logging.WithFields(logrus.Fields{"event uuid": uuid}).Info("starting usecase do job")
	err := doJobEntity.jobWorker.WorkHard(uuid)
	doJobEntity.logging.WithFields(logrus.Fields{"event uuid": uuid}).Info("finished usecase do job")
	return err
}
