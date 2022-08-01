package facades

import (
	"github.com/995933447/log-go/impls/loggerwriters"
	"time"
)

func CheckTimeToOpenNewFileHandlerForFileLogger() loggerwriters.CheckTimeToOpenNewFileFunc {
	return func(lastOpenFileTime *time.Time, isNeverOpenFile bool) (string, bool) {
		if isNeverOpenFile {
			return time.Now().Format("2006010215.log"), true
		}

		if lastOpenFileTime.Hour() != time.Now().Hour() {
			return time.Now().Format("2006010215.log"), true
		}

		lastOpenYear, lastOpenMonth, lastOpenDay := lastOpenFileTime.Date()
		nowYear, nowMonth, nowDay := time.Now().Date()
		if lastOpenDay != nowDay || lastOpenMonth != nowMonth || lastOpenYear != nowYear {
			return time.Now().Format("2006010215.log"), true
		}

		return "", false
	}
}