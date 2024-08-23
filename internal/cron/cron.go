package cron

import (
	"ch-common-package/logger"
	"github.com/robfig/cron/v3"
	"runtime/debug"
)

func NewCron() {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("cron run panic: %v, stack: %s", err, debug.Stack())
		}
	}()

	c := cron.New()

	c.Start()
}
