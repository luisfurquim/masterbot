package masterbot

import (
   "github.com/luisfurquim/stonelizard"
)

func (cfg ConfigT) CertKit() stonelizard.AuthT {
   return cfg.Certkit
}


