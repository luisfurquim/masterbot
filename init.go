package masterbot

import (
   "github.com/wangboo/cron"
)

func init() {
   Kairos = cron.New()
   Kairos.Start()
}

