package masterbot

import (
   "fmt"
   "time"
   "net/http"
)

func (cfg ConfigT) Stop() error {
   var err    error
   var url    string
   var htcli *http.Client
   var resp  *http.Response

   Goose.Logf(2,"Stopping master bot")

   htcli = cfg.HttpsClient(time.Duration(0))
   url   = fmt.Sprintf("https://%s%s/%s/stop", cfg.Host, cfg.Listen, cfg.Id)
   Goose.Logf(2,"Stopping bot %s via %s",cfg.Id,url)
   resp, err = htcli.Get(url)

   if err != nil {
      Goose.Logf(1,"%s %s (%s)",ErrStoppingBot,cfg.Id,err)
      return ErrStoppingBot
   }

   if resp.StatusCode != http.StatusOK {
      Goose.Logf(1,"%s %s (%s)",ErrStatusStoppingBot,cfg.Id,resp.Status)
      return ErrStatusStoppingBot
   }

   return nil
}

