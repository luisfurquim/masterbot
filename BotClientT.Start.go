package masterbot

import (
   "os"
   "fmt"
   "sync"
//   "net/http"
   "golang.org/x/crypto/ssh"
)

func (s *BotClientT) Start(botId string, botInstance int, config []byte, cfg *ConfigT, debugLevel int) error {
   var err        error
   var sshclient *ssh.Client
   var session   *ssh.Session
   var wg         sync.WaitGroup

   if s.Status == BotStatPaused {
      return nil
   }

   err = s.PingAt(botId, botInstance, cfg)
   if err == nil {
      Goose.Logf(2,"bot %s@%s is alive",botId,s.Host[botInstance])
      return nil
   }

   Goose.Logf(2,"Starting bot %s",botId)

   s.Status = BotStatUnreachable
   cfg.SshClientConfig.User = s.SysUser

   sshclient, err = ssh.Dial("tcp", s.Host[botInstance] + ":22", cfg.SshClientConfig)
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

   wg.Add(1)

   go func() {
      defer wg.Done()
      w, _ := session.StdinPipe()
      defer w.Close()
      fmt.Fprintf(w, "%s\n", config)
   }()

   if err = session.Start(fmt.Sprintf("%s%c%s -v %d",s.BinDir, os.PathSeparator, s.BinName, debugLevel)); err != nil {
      Goose.Logf(1,"%s (%s)",ErrFailedStartingBot,err)
      return ErrFailedStartingBot
   }

   wg.Wait()

   s.Status = BotStatRunning

   return nil
}
