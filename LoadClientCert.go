package masterbot

import (
   "os"
   "fmt"
//   "io/ioutil"
   "crypto/tls"
//   "crypto/rsa"
//   "crypto/x509"
   "archive/zip"
//   "encoding/pem"
)

func (cfg *ConfigT) LoadClientCert() error {
   var err    error
   var hn     string
   var cert []byte
   var key  []byte
//   var key   *rsa.PrivateKey

   hn, err = os.Hostname()
   if err != nil {
      Goose.ClientCfg.Logf(1,"Error checking hostname: %s",err)
      return err
   }

   r, err := zip.OpenReader(fmt.Sprintf("%s%c%s.ck",cfg.Pem, os.PathSeparator,hn))
   if err != nil {
      Goose.ClientCfg.Logf(1,"Error decompressing certificate archive: %s",err)
      return err
   }
   defer r.Close()

   // Iterate through the files in the archive.
   for _, f := range r.File {
      rc, err := f.Open()
      if err != nil {
         Goose.ClientCfg.Logf(1,"Error opening %s: %s",f.Name,err)
         return err
      }

      switch f.Name {
         case "client.crt":
            _, cert, err = cfg.Certkit.ReadCertFromReader(rc)
         case "client.key":
            _, key, err = cfg.Certkit.ReadDecryptRsaPrivKeyFromReader(rc)
      }
      rc.Close()

      if err != nil {
         Goose.ClientCfg.Logf(1,"Error loading %s: %s",f.Name,err)
         return err
      }
   }
   cfg.ClientCert, err = tls.X509KeyPair(cert, key)
   if err != nil {
      Goose.ClientCfg.Logf(1,"Error setting client keypair: %s",err)
      return err
   }

   cfg.ClientCA   = cfg.Certkit.GetCertPool()

   return nil
}

