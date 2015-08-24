package masterbot

import (
   "github.com/robfig/cron"
)

func init() {
   Kairos = cron.New()
   Kairos.Start()
}

