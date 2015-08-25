package masterbot

import (
   "io"
   "github.com/luisfurquim/slavebot"
)

func (svc ServiceT) GetConfig() (io.Reader, error) {
   return slavebot.ConfigReader(), nil
}
