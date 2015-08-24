package masterbot

import (
//   "fmt"
   "io/ioutil"
   "crypto/tls"
   "crypto/rsa"
   "crypto/x509"
   "encoding/pem"
)

func (cfg *ConfigT) LoadClientCert() error {
   var err error
   var caCert []byte
   var cert, key, plainkey []byte
   var pemblockkey, pemblockcert *pem.Block
   var rsakey *rsa.PrivateKey

   key, err = ioutil.ReadFile(cfg.Pem + "/client.key")
   if err != nil {
      return err
   }

   cert, err = ioutil.ReadFile(cfg.Pem + "/client.crt")
   if err != nil {
      return err
   }


   pemblockkey, _ = pem.Decode(key)

   plainkey, err = x509.DecryptPEMBlock(pemblockkey,[]byte{})
   if err != nil {
      return err
   }

   rsakey, err = x509.ParsePKCS1PrivateKey(plainkey)
   if err != nil {
      return err
   }

   pemblockcert, _ = pem.Decode(cert)

//   fmt.Printf("%s\n\n%#v\n",plainkey,pemblockkey)

/*
   // Load client cert
   cfg.ClientCert, err = tls.X509KeyPair(cert, plainkey)
   if err != nil {
      return err
   }
*/

   cfg.ClientCert = tls.Certificate{
      Certificate: [][]byte{pemblockcert.Bytes},
      PrivateKey: rsakey,
   }

   // Load CA cert
   caCert, err = ioutil.ReadFile(cfg.Pem + "/rootCA.crt")
   if err != nil {
      return err
   }
   cfg.ClientCA = x509.NewCertPool()
   cfg.ClientCA.AppendCertsFromPEM(caCert)

//   fmt.Printf("%#v\n\n%#v\n",cfg.ClientCert,cfg.ClientCA)

   return nil
}

