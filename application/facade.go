package application

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/thevan4/golang-ddd-template/domain"
	"github.com/thevan4/golang-ddd-template/portadapter"
	"github.com/thevan4/golang-ddd-template/usecase"
)

// Facade for interacting with components
type Facade struct {
	JobWayOne        *portadapter.JobWayOne
	JobWayTwo        *portadapter.JobWayTwo
	GracefulShutdown *domain.GracefulShutdown
	UUIDgenerator    domain.UUIDgenerator
	Logging          *logrus.Logger
}

// NewFacade ...
func NewFacade(jobWayOne *portadapter.JobWayOne,
	jobWayTwo *portadapter.JobWayTwo,
	gracefulShutdown *domain.GracefulShutdown,
	uuidGenerator domain.UUIDgenerator,
	logging *logrus.Logger) *Facade {

	return &Facade{
		JobWayOne:        jobWayOne,
		JobWayTwo:        jobWayTwo,
		GracefulShutdown: gracefulShutdown,
		UUIDgenerator:    uuidGenerator,
		Logging:          logging,
	}
}

// DoJob example of facade instance
func (facade *Facade) DoJob(jobType string, uuid string) error {
	switch jobType {
	// selecting an interface implementation
	case "one":
		facade.Logging.WithFields(logrus.Fields{"event uuid": uuid}).Infof("take portadapter for job type %v", jobType)
		doJobEntity := usecase.NewDoJobEntity(facade.JobWayOne, facade.GracefulShutdown, facade.Logging)
		return doJobEntity.DoJob(uuid)
	case "two":
		facade.Logging.WithFields(logrus.Fields{"event uuid": uuid}).Infof("take portadapter for job type %v", jobType)
		doJobEntity := usecase.NewDoJobEntity(facade.JobWayTwo, facade.GracefulShutdown, facade.Logging)
		return doJobEntity.DoJob(uuid)
	default:
		return fmt.Errorf("unknown job type for do job usecase: %v", jobType) // also must be checked to the logic of the facade
	}
}
