package masterbot

import (
   "fmt"
   "net/http"
)



func (s *BotClientT) PingAt(botId string, botInstance int, cfg *ConfigT) error {
   var err         error
   var resp       *http.Response
   var url         string

   Goose.Ping.Logf(2,"Pinging %s@%s",botId,s.Host[botInstance].Name)
   url = fmt.Sprintf("https://%s%s/%s/ping", s.Host[botInstance].Name, s.Listen, botId)
   resp, err = cfg.HttpsPingClient.Get(url)

   if err != nil {
      Goose.Ping.Logf(1,"%s (%s) %#v",ErrFailedPingingBot,err,resp)
      return ErrFailedPingingBot
   }

   if resp != nil {
      defer resp.Body.Close()
   }

   if resp.StatusCode != http.StatusNoContent {
      Goose.Ping.Logf(1,"%s %s at %s (status code=%d)",ErrFailedPingingBot,botId,url,resp.StatusCode)
      return ErrFailedPingingBot
   }

   return nil
}

