package masterbot

import (
   "os"
   "fmt"
   "sync"
   "io/ioutil"
//   "net/http"
   "golang.org/x/crypto/ssh"
)

func (cfg *ConfigT) Start(config []byte, dbgLevel int) error {
   var err         error
   var sshclient  *ssh.Client
   var session    *ssh.Session
   var wg, subwg   sync.WaitGroup
   var host        string
   var multiErr  []error
   var botInstance int

   err = cfg.PingAt()
   if err == nil {
      Goose.Logf(2,"bot %s is alive",cfg.Id)
      return nil
   }

   Goose.Logf(2,"Starting bot %s",cfg.Id)

   cfg.SshClientConfig.User = cfg.SysUser


   wg.Add(len(cfg.Host))

   multiErr = make([]error,len(cfg.Host))
   for botInstance, host = range cfg.Host {
      go func(instance int) {
         defer wg.Done()

         sshclient, err = ssh.Dial("tcp", host + ":22", cfg.SshClientConfig)
         if err != nil {
            Goose.Logf(1,"%s (%s)",ErrDialingToBot,err)
            multiErr[instance] = ErrDialingToBot
            return
         }

         session, err = sshclient.NewSession()
         if err != nil {
            Goose.Logf(1,"%s (%s)",ErrCreatingSession,err)
            multiErr[instance] = ErrCreatingSession
            return
         }
         defer session.Close()

         subwg.Add(1)

         go func() {
            defer subwg.Done()
            Goose.Logf(6,"Sending config")
            w, _ := session.StdinPipe()
            defer w.Close()
            fmt.Fprintf(w, "%s\n", config)
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

         if err = session.Start(fmt.Sprintf("%s%c%s -v %d",cfg.BinDir, os.PathSeparator, cfg.BinName, dbgLevel)); err != nil {
            Goose.Logf(1,"%s (%s)",ErrFailedStartingBot,err)
            multiErr[instance] = ErrFailedStartingBot
            return
         }

         subwg.Wait()

      }(botInstance)
   }

   wg.Wait()

   for _, err = range multiErr {
      if err != nil {
         return err
      }
   }

   Goose.Logf(2,"Bot %s started successfully",cfg.Id)
   return nil
}
