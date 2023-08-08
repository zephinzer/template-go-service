package example

import (
	"time"

	"github.com/sirupsen/logrus"
)

func StartJob(waitFor time.Duration) {
	logrus.Infof("simulating a job run for %v...", waitFor)
	<-time.After(3 * time.Second)
	logrus.Infof("job run complete!")
}
