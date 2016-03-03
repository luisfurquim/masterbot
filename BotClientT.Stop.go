package masterbot

import (
   "fmt"
   "time"
   "net/http"
)

func (svc BotClientT) Stop(botId string, cfg *ConfigT) error {
   var err     error
   var errHist error
   var url     string
   var htcli  *http.Client
   var resp   *http.Response
   var host    string

   Goose.StartStop.Logf(2,"Stopping slave bot %s",botId)

   htcli = cfg.HttpsClient(time.Duration(0))

   for _, host = range svc.Host{
      url   = fmt.Sprintf("https://%s%s/%s/stop", host, svc.Listen, botId)
      Goose.StartStop.Logf(2,"Stopping bot %s via %s",botId,url)
      resp, err = htcli.Get(url)

      if err != nil {
         Goose.StartStop.Logf(1,"%s %s (%s)",ErrStoppingBot,botId,err)
         errHist = ErrStoppingBot
         continue
      }

      if resp.StatusCode != http.StatusNoContent {
         Goose.StartStop.Logf(1,"%s %s (%s)",ErrStatusStoppingBot,botId,resp.Status)
         errHist = ErrStatusStoppingBot
         continue
      }
   }

   return errHist
}
