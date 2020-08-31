package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thevan4/golang-ddd-template/application"
	"github.com/thevan4/golang-ddd-template/domain"
	"github.com/thevan4/golang-ddd-template/portadapter"
)

var rootCmd = &cobra.Command{
	Use:   "run",
	Short: "example-program 😉",
	Run: func(cmd *cobra.Command, args []string) {
		uuidGenerator := portadapter.NewUUIDGenerator()
		uuidForRootProcess := uuidGenerator.NewUUID().UUID.String()

		// validate fields
		logging.WithFields(logrus.Fields{
			"event uuid":       uuidForRootProcess,
			"config file path": viperConfig.GetString(configFilePathName),
			"log format":       viperConfig.GetString(logFormatName),
			"log level":        viperConfig.GetString(logLevelName),
			"log output":       viperConfig.GetString(logOutputName),
			"syslog tag":       viperConfig.GetString(syslogTagName),

			"max shutdown work receiver time": viperConfig.GetDuration(maxShutdownTaskReceiverTimeName),
			// "version": version,
			// "commit":  commit,
			// "branch":  branch,
		}).Info("")

		// getting ready to catch system signals
		signalChan := make(chan os.Signal, 2)
		signal.Notify(signalChan, syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)

		gracefulShutdown := &domain.GracefulShutdown{} // Somehow
		// portadapters init
		doJobWayOne := portadapter.NewJobWayOne(logging)
		doJobWayTwo := portadapter.NewJobWayTwo(logging)
		// portadapters init end

		facade := application.NewFacade(doJobWayOne,
			doJobWayTwo,
			gracefulShutdown,
			uuidGenerator,
			logging)

		// init and start task receiver
		gracefulShutdownForTaskReceiver := make(chan struct{}, 1)
		taskReceiverIsDone := make(chan struct{}, 1)
		taskReceiver := application.NewTaskReceiver(facade, gracefulShutdownForTaskReceiver, taskReceiverIsDone)
		// go taskReceiver.TaskReceiver()

		logging.WithFields(logrus.Fields{"event uuid": uuidForRootProcess}).Info("program running")
		for i := 0; i < 10; i++ { // TODO: more elegant somehow
			go taskReceiver.DoJob("two")
		}

		<-signalChan // shutdown signal
		logging.WithFields(logrus.Fields{"event uuid": uuidForRootProcess}).Info("got shutdown signal")

		gracefulShutdownForTaskReceiver <- struct{}{}                                                                 // FIXME: some broken here?
		gracefulShutdownUsecases(gracefulShutdown, viperConfig.GetDuration(maxShutdownTaskReceiverTimeName), logging) // TODO: rework that
		<-taskReceiverIsDone
		logging.WithFields(logrus.Fields{"event uuid": uuidForRootProcess}).Info("task receiver is done is Done")
		// TODO: usecases done wait
		logging.WithFields(logrus.Fields{"event uuid": uuidForRootProcess}).Info("program stopped")
	},
}

func gracefulShutdownUsecases(gracefulShutdown *domain.GracefulShutdown, maxWaitTimeForJobsIsDone time.Duration, logging *logrus.Logger) {
	gracefulShutdown.Lock()
	gracefulShutdown.ShutdownNow = true
	gracefulShutdown.Unlock()

	ticker := time.NewTicker(time.Duration(100 * time.Millisecond)) // hardcode
	defer ticker.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), maxWaitTimeForJobsIsDone)
	defer cancel()
	for {
		select {
		case <-ticker.C:
			gracefulShutdown.Lock()
			if gracefulShutdown.UsecasesJobs <= 0 {
				logging.Info("All jobs is done")
				defer gracefulShutdown.Unlock()
				return
			}
			gracefulShutdown.Unlock()
			continue
		case <-ctx.Done():
			gracefulShutdown.Lock()
			logging.Warnf("%v jobs is fail when program stop", gracefulShutdown.UsecasesJobs)
			defer gracefulShutdown.Unlock()
			return
		}
	}
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
