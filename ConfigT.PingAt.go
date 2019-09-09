package masterbot

import (
   "fmt"
   "sync"
   "net/http"
)



func (cfg *ConfigT) PingAt() error {
   var err         error
   var host        Host
   var botError  []error
   var botInstance int
   var wg          sync.WaitGroup

   if len(cfg.Host) == 0 {
      Goose.Ping.Logf(1,"%s",ErrNoBotsToPing)
      return ErrNoBotsToPing
   }
	
   wg.Add(len(cfg.Host))

   botError = make([]error,len(cfg.Host))
   for botInstance, host = range cfg.Host {
      go func(instance int, host Host) {
         var err         error
         var resp       *http.Response
         var url         string

         defer wg.Done()

         if host.Status == BotStatRunning {
            url   = fmt.Sprintf("https://%s%s/%s/ping", host.Name, cfg.Listen, cfg.Id)
            Goose.Ping.Logf(6, "HttpsClient=%p, %#v, %T", cfg.HttpsPingClient, cfg.HttpsPingClient, *cfg.HttpsPingClient)
            Goose.Ping.Logf(6,"Pinging bot at %s using %v",url,cfg.HttpsPingClient)
            resp, err = cfg.HttpsPingClient.Get(url)
            if resp != nil {
               defer resp.Body.Close()
            }

            if err != nil {
               Goose.Ping.Logf(1,"%s %s@%s (%s) %#v",ErrFailedPingingBot,cfg.Id,host.Name,err,resp)
               botError[instance] = ErrFailedPingingBot
               return
            }

            if resp.StatusCode != http.StatusNoContent {
               Goose.Ping.Logf(1,"%s %s@%s at %s (status code=%d)",ErrFailedPingingBot,cfg.Id,host.Name,url,resp.StatusCode)
               botError[instance] = ErrFailedPingingBot
               return
            }
         } else {
            Goose.Ping.Logf(1,"%s %s@%s at %s (not running)",ErrFailedPingingBot,cfg.Id,host.Name,url)
            Goose.Ping.Logf(6,"%#v",cfg.Host)
            botError[instance] = ErrFailedPingingBot
            //Goose.Ping.Logf(1,"Ignore Ping Bot to %s@%s at %s because Status = %s",cfg.Id,host.Name,url,host.Status)
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
