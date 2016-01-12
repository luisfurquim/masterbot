package masterbot

import (
   "github.com/luisfurquim/stonelizard/certkit"
)

func (cfg ConfigT) CertKit() *certkit.CertKit {
   return cfg.Certkit
}


