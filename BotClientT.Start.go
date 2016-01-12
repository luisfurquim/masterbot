package masterbot

import (
   "os"
   "fmt"
   "sync"
//   "strings"
//   "net/http"
   "golang.org/x/crypto/ssh"
)

func (s *BotClientT) Start(botId string, botInstance int, cmdline string, cfg *ConfigT, debugLevel int) error {
   var err        error
   var sshclient *ssh.Client
   var session   *ssh.Session
   var wg         sync.WaitGroup
   var cmd        string

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
      Goose.Logf(2,"Closing stdin for bot %s",botId)
      //fmt.Fprintf(w, "%s\n", config)
   }()

   cmd = fmt.Sprintf("%s%c%s -v %d %s",s.BinDir, os.PathSeparator, s.BinName, debugLevel, cmdline)
   if err = session.Start(cmd); err != nil {
      Goose.Logf(1,"%s (%s)",ErrFailedStartingBot,err)
      return ErrFailedStartingBot
   }

   wg.Wait()

   Goose.Logf(2,"Started bot %s with cmd:[%s]",botId,cmd)

   s.Status = BotStatRunning

   return nil
}
