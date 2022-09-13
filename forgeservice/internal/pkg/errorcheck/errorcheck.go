package errorcheck

import "github.com/sirupsen/logrus"

func CheckLogFatal(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}

func CheckLogError(err error) {
	if err != nil {
		logrus.Error(err)
	}
}
