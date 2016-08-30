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

   if s.Host[botInstance].Status == BotStatPaused {
      return nil
   }

   if s.Host[botInstance].Status == BotStatRunning {
      err = s.PingAt(botId, botInstance, cfg)
      if err == nil {
         Goose.Ping.Logf(2,"bot %s@%s is alive",botId,s.Host[botInstance].Name)
         return nil
      }
   }

   Goose.StartStop.Logf(2,"Starting bot %s@%s",botId,s.Host[botInstance].Name)

   s.Host[botInstance].Status = BotStatUnreachable
   if s.Host[botInstance].OnStatUpdate != nil {
      s.Host[botInstance].OnStatUpdate(BotStatUnreachable)
   }

   cfg.SshClientConfig.User   = s.SysUser

   sshclient, err = ssh.Dial("tcp", s.Host[botInstance].Name + ":22", cfg.SshClientConfig)
   if err != nil {
      Goose.StartStop.Logf(1,"%s (%s)",ErrDialingToBot,err)
      return ErrDialingToBot
   }

   session, err = sshclient.NewSession()
   if err != nil {
      Goose.StartStop.Logf(1,"%s (%s)",ErrCreatingSession,err)
      return ErrCreatingSession
   }
   defer session.Close()


   wg.Add(1)

   go func() {
      defer wg.Done()
      w, _ := session.StdinPipe()
      defer w.Close()
      Goose.StartStop.Logf(2,"Closing stdin for bot %s",botId)
      //fmt.Fprintf(w, "%s\n", config)
   }()

   cmd = fmt.Sprintf("%s%c%s -v %d %s",s.BinDir, os.PathSeparator, s.BinName, debugLevel, cmdline)
   if err = session.Start(cmd); err != nil {
      Goose.StartStop.Logf(1,"%s (%s)",ErrFailedStartingBot,err)
      return ErrFailedStartingBot
   }

   wg.Wait()

   Goose.StartStop.Logf(2,"Started bot %s with cmd:[%s]",botId,cmd)

   s.Host[botInstance].Status = BotStatRunning
   if s.Host[botInstance].OnStatUpdate != nil {
      s.Host[botInstance].OnStatUpdate(BotStatRunning)
   }

   return nil
}
