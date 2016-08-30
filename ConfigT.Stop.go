package masterbot

import (
   "fmt"
   "sync"
   "net/http"
)

func (cfg ConfigT) Stop() error {
   var err         error
   var botError  []error
   var botInstance int
   var wg          sync.WaitGroup
   var host        Host

   Goose.StartStop.Logf(2,"Stopping master bot")

   wg.Add(len(cfg.Host))

   botError = make([]error,len(cfg.Host))
   for botInstance, host = range cfg.Host {
      go func(instance int, h Host) {
         var err         error
         var resp       *http.Response
         var url         string

         defer wg.Done()

         if h.Status == BotStatRunning {
            url   = fmt.Sprintf("https://%s%s/%s/stop", h.Name, cfg.Listen, cfg.Id)
            Goose.StartStop.Logf(2,"Stopping bot %s@%s via %s",cfg.Id,h.Name,url)
            resp, err = cfg.HttpsStopClient.Get(url)

            if resp != nil {
               defer resp.Body.Close()
            }

            if err != nil {
               Goose.StartStop.Logf(1,"%s %s@%s (%s)",ErrStoppingBot,cfg.Id,h.Name,err)
               botError[instance] = ErrStoppingBot
               return
            }

            if resp.StatusCode != http.StatusNoContent {
               Goose.StartStop.Logf(1,"%s %s@%s (%s)",ErrStatusStoppingBot,cfg.Id,h.Name,resp.Status)
               botError[instance] = ErrStatusStoppingBot
               return
            }

            h.Status = BotStatStopped
            if h.OnStatUpdate != nil {
               h.OnStatUpdate(BotStatStopped)
            }
         }
      }(botInstance,host)
   }

   wg.Wait()

   for _, err = range botError {
      if err != nil {
         return err
      }
   }

   return nil
}

