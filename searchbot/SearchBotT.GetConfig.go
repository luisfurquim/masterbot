package searchbot

import (
   "io"
   "github.com/luisfurquim/slavebot"
)

func (sb SearchBotT) GetConfig() (io.Reader, error) {
   return slavebot.ConfigReader(), nil
}
