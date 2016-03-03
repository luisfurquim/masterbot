package masterbot

import (
   "io/ioutil"
   "golang.org/x/crypto/ssh"
)

func (cfg *ConfigT) NewSshClientConfig(privKeyPath string) error {
   var err         error
   var clientKey []byte
   var signer      ssh.Signer

   Goose.ClientCfg.Logf(4,"Reading SSH private key from %s",privKeyPath)

   clientKey, err = ioutil.ReadFile(privKeyPath)
   if err != nil {
      Goose.ClientCfg.Logf(1,"%s (%s)",ErrReadingSSHKeys,err)
      return ErrReadingSSHKeys
   }

   signer, err = ssh.ParsePrivateKey(clientKey)
   if err != nil {
      Goose.ClientCfg.Logf(1,"%s (%s)",ErrParsingSSHKeys,err)
      return ErrReadingSSHKeys
   }

   cfg.SshClientConfig = &ssh.ClientConfig{
//      User: user,
      Auth: []ssh.AuthMethod{
         ssh.PublicKeys(signer),
      },
   }

   return nil
}
