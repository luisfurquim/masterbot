package masterbot

import (
   "fmt"
   "time"
   "net/http"
)

func (cfg ConfigT) Stop() error {
   var err          error
   var url          string
   var htcli       *http.Client
   var resp        *http.Response
   var botInstance  int

   Goose.Logf(2,"Stopping master bot")

   htcli = cfg.HttpsClient(time.Duration(0))

   for botInstance, _ = range cfg.Host {
      url   = fmt.Sprintf("https://%s%s/%s/stop", cfg.Host[botInstance], cfg.Listen, cfg.Id)
      Goose.Logf(2,"Stopping bot %s@%s via %s",cfg.Id,cfg.Host[botInstance],url)
      resp, err = htcli.Get(url)

      if err != nil {
         Goose.Logf(1,"%s %s@%s (%s)",ErrStoppingBot,cfg.Id,cfg.Host[botInstance],err)
         return ErrStoppingBot
      }

      if resp.StatusCode != http.StatusNoContent {
         Goose.Logf(1,"%s %s@%s (%s)",ErrStatusStoppingBot,cfg.Id,cfg.Host[botInstance],resp.Status)
         return ErrStatusStoppingBot
      }
   }

   return nil
}

