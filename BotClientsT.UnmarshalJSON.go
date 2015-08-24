package masterbot

import (
//   "fmt"
   "encoding/json"
)

type botClientsT BotClientsT


func (bc *BotClientsT) UnmarshalJSON(buf []byte) error {
   var bctemp botClientsT
   var err error

   if err = json.Unmarshal(buf, &bctemp); err == nil {
      *bc = BotClientsT(bctemp)
   }

   return err
}
