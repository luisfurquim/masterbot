package masterbot

import (
   "os"
   "fmt"
//   "sync"
//   "bytes"
//   "strings"
//   "net/http"
   "golang.org/x/crypto/ssh"
)

func (s *BotClientT) Start(botId string, botInstance int, cmdline string, cfg *ConfigT, debugLevel int) error {
   var err        error
   var sshclient *ssh.Client
   var session   *ssh.Session
//   var wg         sync.WaitGroup
   var cmd        string
   var sshport    string

   sshport = SSHPort
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

   if s.Host[botInstance].Port != "" {
      sshport = s.Host[botInstance].Port
   }

   sshclient, err = ssh.Dial("tcp", s.Host[botInstance].Name + ":" + sshport, cfg.SshClientConfig)
   if err != nil {
      Goose.StartStop.Logf(1,"%s (%s)",ErrDialingToBot,err)
      return ErrDialingToBot
   }
   defer sshclient.Close()

   Goose.StartStop.Logf(3,"Dialed to bot %s@%s:%s",botId,s.Host[botInstance].Name,sshport)

   session, err = sshclient.NewSession()
   if err != nil {
      Goose.StartStop.Logf(1,"%s (%s)",ErrCreatingSession,err)
      return ErrCreatingSession
   }
   defer session.Close()

   Goose.StartStop.Logf(3,"Session started at bot %s@%s:%s",botId,s.Host[botInstance].Name,sshport)

/*
   wg.Add(1)

   go func() {
      defer wg.Done()
      w, _ := session.StdinPipe()
      defer w.Close()
      Goose.StartStop.Logf(2,"Closing stdin for bot %s",botId)
      //fmt.Fprintf(w, "%s\n", config)
   }()
*/

/*
   // Set up terminal modes
   modes := ssh.TerminalModes{
      ssh.ECHO:          0,     // disable echoing
      ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
      ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
   }
   // Request pseudo terminal
   if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
      Goose.StartStop.Fatalf(1,"request for pseudo terminal failed: %s", err)
   }


   session.Stdout = &bytes.Buffer{}
   session.Stderr = &bytes.Buffer{}
*/

   cmd = fmt.Sprintf("%s%c%s -v %d -path %s %s",s.BinDir, os.PathSeparator, s.BinName, debugLevel, s.WorkDir, cmdline)
//   cmd = fmt.Sprintf("%s%c%s -v %d %s",s.BinDir, os.PathSeparator, s.BinName, debugLevel, cmdline)

   Goose.StartStop.Logf(3,"Will run %s@%s:%s using %s", botId, s.Host[botInstance].Name, sshport, cmd)

   err = session.Start(cmd)
//   err = session.Run(cmd)
   Goose.StartStop.Logf(2,"Running bot %s",botId)
//   wg.Wait()

   if err != nil {
      session.Signal(ssh.SIGKILL)
      Goose.StartStop.Logf(1,"%s (%s)",ErrFailedStartingBot,err)
      return ErrFailedStartingBot
   }

   Goose.StartStop.Logf(2,"Started bot %s with cmd:[%s]",botId,cmd)

   s.Host[botInstance].Status = BotStatRunning
   if s.Host[botInstance].OnStatUpdate != nil {
      s.Host[botInstance].OnStatUpdate(BotStatRunning)
   }

   return nil
}
