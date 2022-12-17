package tests

import "github.com/sirupsen/logrus"

func QuickTestLogger() *logrus.Logger {
	return logrus.New()
}
