package masterbot

import (
   "fmt"
   "sync"
   "net/http"
)



func (cfg *ConfigT) PingAt() error {
   var err         error
   var host        string
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
      go func(instance int, host string) {
         var err         error
         var resp       *http.Response
         var url         string

         defer wg.Done()

         url   = fmt.Sprintf("https://%s%s/%s/ping", host, cfg.Listen, cfg.Id)
         Goose.Ping.Logf(6,"Pinging bot at %s using %v",url,cfg.HttpsPingClient)
         resp, err = cfg.HttpsPingClient.Get(url)

         if resp != nil {
            defer resp.Body.Close()
         }

         if err != nil {
            Goose.Ping.Logf(1,"%s %s@%s (%s) %#v",ErrFailedPingingBot,cfg.Id,host,err,resp)
            botError[instance] = ErrFailedPingingBot
            return
         }

         if resp.StatusCode != http.StatusNoContent {
            Goose.Ping.Logf(1,"%s %s@%s at %s (status code=%d)",ErrFailedPingingBot,cfg.Id,host,url,resp.StatusCode)
            botError[instance] = ErrFailedPingingBot
            return
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
