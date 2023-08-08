package example

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func StartWorker() {
	var waiter sync.WaitGroup

	logrus.Infof("starting a worker that listens for SIGINT/SIGTERM. use <ctrl+c> to terminate")
	<-time.After(1 * time.Second)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	isDone := false
	isDoneLock := sync.Mutex{}
	timestampFormat := "2006-01-02 15:04:05 -0700"

	go func() {
		sig := <-sigs
		logrus.Infof("received signal[%s]", sig.String())
		isDoneLock.Lock()
		isDone = true
		isDoneLock.Unlock()
	}()

	waiter.Add(1)
	go func() {
		for {
			var doneValue bool
			isDoneLock.Lock()
			doneValue = isDone
			isDoneLock.Unlock()
			if doneValue {
				break
			}
			logrus.Infof("process 1 running as of %v", time.Now().Format(timestampFormat))
			<-time.After(1 * time.Second)
		}
		logrus.Infof("process 1 is gracefully shut down")
		waiter.Done()
	}()

	waiter.Add(1)
	go func() {
		for {
			var doneValue bool
			isDoneLock.Lock()
			doneValue = isDone
			isDoneLock.Unlock()
			if doneValue {
				break
			}
			logrus.Infof("process 2 running as of %v", time.Now().Format(timestampFormat))
			<-time.After(1 * time.Second)
		}
		logrus.Infof("process 2 is gracefully shut down")
		waiter.Done()
	}()

	waiter.Wait()
	logrus.Infof("worker exited gracefully :D")
}
