package masterbot

import (
   "fmt"
   "time"
   "sync"
   "net/http"
)



func (cfg *ConfigT) PingAt() error {
   var err         error
   var host        string
   var multiErr  []error
   var botInstance int
   var wg          sync.WaitGroup

   wg.Add(len(cfg.Host))

   multiErr = make([]error,len(cfg.Host))
   for botInstance, host = range cfg.Host {
      go func(instance int, host string) {
         var err         error
         var httpclient *http.Client
         var resp       *http.Response
         var url         string

         defer wg.Done()

         httpclient = cfg.HttpsClient(cfg.BotPingTimeout * time.Second)
         url   = fmt.Sprintf("https://%s%s/%s/ping", host, cfg.Listen, cfg.Id)
         resp, err = httpclient.Get(url)

         if err != nil {
            Goose.Logf(1,"%s %s@%s (%s) %#v",ErrFailedPingingBot,cfg.Id,host,err,resp)
            multiErr[instance] = ErrFailedPingingBot
            return
         }

         if resp.StatusCode != http.StatusOK {
            Goose.Logf(1,"%s %s@%s at %s (status code=%d)",ErrFailedPingingBot,cfg.Id,host,url,resp.StatusCode)
            multiErr[instance] = ErrFailedPingingBot
            return
         }
      }(botInstance,host)
   }

   wg.Wait()

   for _, err = range multiErr {
      if err != nil {
         return err
      }
   }

   return nil
}

/*
   tr = &http.Transport{
      TLSClientConfig:    &tls.Config{
         //RootCAs: pool // crypto/x509
      },
      DisableCompression: true,
   }

   Client: &http.Client{
      Transport:  tr,
   },

   req, err = http.NewRequest("POST",cfg.Auth[0].URL, strings.NewReader(form_data))
   if err != nil {
      return nil, err
   }

   resp, err = csi.Client.Do(req)
   if err != nil {
      return nil, err
   }

   body, err = ioutil.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   resp.Body.Close()
*/
