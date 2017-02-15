package masterbot


import (
   "fmt"
   "time"
)


func (t Timeout) MarshalJSON() ([]byte, error) {

   return []byte(fmt.Sprintf("\"%s\"",time.Duration(t))), nil
}

