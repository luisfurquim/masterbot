package masterbot

import (
   "os"
   "fmt"
   "io/ioutil"
//   "net/http"
   "golang.org/x/crypto/ssh"
)

func (cfg *ConfigT) Start(config []byte, dbgLevel int) error {
   var err        error
   var sshclient *ssh.Client
   var session   *ssh.Session

   err = cfg.PingAt()
   if err == nil {
      Goose.Logf(2,"bot %s is alive",cfg.Id)
      return nil
   }

   Goose.Logf(2,"Starting bot %s",cfg.Id)

   cfg.SshClientConfig.User = cfg.SysUser

   sshclient, err = ssh.Dial("tcp", cfg.Host + ":22", cfg.SshClientConfig)
   if err != nil {
      Goose.Logf(1,"%s (%s)",ErrDialingToBot,err)
      return ErrDialingToBot
   }

   session, err = sshclient.NewSession()
   if err != nil {
      Goose.Logf(1,"%s (%s)",ErrCreatingSession,err)
      return ErrCreatingSession
   }
   defer session.Close()

   go func() {
      Goose.Logf(6,"Sending config")
      w, _ := session.StdinPipe()
      defer w.Close()
      fmt.Fprintf(w, "%s", config)
   }()


   go func() {
      Goose.Logf(6,"getting stdout")
      w, _ := session.StdoutPipe()

      output, err := ioutil.ReadAll(w)
      if err != nil {
         Goose.Logf(1,"Error reading SSH output (%s)",err)
      } else {
         Goose.Logf(6,"SSH stdout Read: %s",output)
      }
   }()

   go func() {
      Goose.Logf(6,"getting stderr")
      w, _ := session.StderrPipe()

      output, err := ioutil.ReadAll(w)
      if err != nil {
         Goose.Logf(1,"Error reading stderr (%s)",err)
      } else if len(output) > 0 {
         Goose.Logf(1,"SSH stderr Read: %s",output)
      }
   }()

   Goose.Logf(6,"SSH starting %s%c%s -v %d",cfg.BinDir, os.PathSeparator, cfg.BinName, dbgLevel)

   if err = session.Start(fmt.Sprintf("%s%c%s",cfg.BinDir, os.PathSeparator, cfg.BinName)); err != nil {
      Goose.Logf(1,"%s (%s)",ErrFailedStartingBot,err)
      return ErrFailedStartingBot
   }

   Goose.Logf(6,"Bot %s started successfully",cfg.Id)

   return nil
}
