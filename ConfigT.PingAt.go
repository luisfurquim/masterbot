package masterbot

import (
   "fmt"
   "time"
   "net/http"
)



func (cfg *ConfigT) PingAt() error {
   var err         error
   var httpclient *http.Client
   var resp       *http.Response
   var url         string


   httpclient = cfg.HttpsClient(cfg.BotPingTimeout * time.Second)
   url   = fmt.Sprintf("https://%s%s/%s/ping", cfg.Host, cfg.Listen, cfg.Id)
   resp, err = httpclient.Get(url)

   if err != nil {
      Goose.Logf(1,"%s (%s) %#v",ErrFailedPingingBot,err,resp)
      return ErrFailedPingingBot
   }

   if resp.StatusCode != http.StatusOK {
      Goose.Logf(1,"%s %s at %s (status code=%d)",ErrFailedPingingBot,cfg.Id,url,resp.StatusCode)
      return ErrFailedPingingBot
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
