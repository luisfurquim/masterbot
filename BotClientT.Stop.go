package masterbot

import (
   "fmt"
   "time"
   "net/http"
)

func (svc BotClientT) Stop(botId string, cfg *ConfigT) error {
   var err    error
   var url    string
   var htcli *http.Client
   var resp  *http.Response

   Goose.Logf(2,"Stopping slave bot %s",botId)

   htcli = cfg.HttpsClient(time.Duration(0))
   url   = fmt.Sprintf("https://%s%s/%s/stop", svc.Host, svc.Listen, botId)
   Goose.Logf(2,"Stopping bot %s via %s",botId,url)
   resp, err = htcli.Get(url)

   if err != nil {
      Goose.Logf(1,"%s %s (%s)",ErrStoppingBot,botId,err)
      return ErrStoppingBot
   }

   if resp.StatusCode != http.StatusOK {
      Goose.Logf(1,"%s %s (%s)",ErrStatusStoppingBot,botId,resp.Status)
      return ErrStatusStoppingBot
   }

   return nil
}

