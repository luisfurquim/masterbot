package masterbot


import (
   "time"
)


func (t *Timeout) Set(val string) {
   var d    time.Duration
   var err  error

   d, err = time.ParseDuration(val)
   Goose.ClientCfg.Logf(5,"BotCommTimeout: %s (err: %s)",d, err)
   if err != nil {
      d, err = time.ParseDuration(val + "s")
   }

   (*t) = Timeout(d)
}

