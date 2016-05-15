package job

import (
	"time"
)

type Job struct {
	Start time.Time
	Data  []byte
}
