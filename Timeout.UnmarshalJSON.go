package masterbot


import (
   "time"
)


func (t *Timeout) UnmarshalJSON(b []byte) error {
   var d time.Duration
   var err error

   Goose.StartStop.Fatalf(0,"%s",time.Duration(*t))
   d, err = time.ParseDuration(string(b))
   if err != nil {
      return err
   }

   *t = Timeout(d)
   return nil
}

