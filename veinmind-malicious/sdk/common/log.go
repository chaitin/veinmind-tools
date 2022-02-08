package common

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log = logrus.New()

func InitLogger() {
	Log.Out = os.Stdout
	Log.Level = logrus.InfoLevel
}
