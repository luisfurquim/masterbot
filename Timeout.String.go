package masterbot


import (
   "fmt"
   "time"
)


func (t Timeout) String() string {
   return fmt.Sprintf("%s",time.Duration(t))
}

