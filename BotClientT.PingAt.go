package masterbot

import (
   "fmt"
   "net"
   "time"
   "net/http"
   "crypto/tls"
)



func (s *BotClientT) PingAt(botId string, botInstance int, cfg *ConfigT) error {
   var err         error
   var httpclient  http.Client
   var resp       *http.Response
   var tlsConfig  *tls.Config
   var url         string

   tlsConfig = &tls.Config{
      Certificates: []tls.Certificate{cfg.ClientCert},
      RootCAs:      cfg.ClientCA,
      InsecureSkipVerify: true,
   }
//   tlsConfig.BuildNameToCertificate()

   httpclient = http.Client{
      Transport: &http.Transport{
//         Timeout:             BOTPINGTMOUT,
         Dial: func (network, addr string) (net.Conn, error) {
            return net.DialTimeout(network, addr, cfg.BotPingTimeout * time.Second)
         },
         TLSClientConfig:     tlsConfig,
         DisableCompression:  true,
      },
   }

   url = fmt.Sprintf("https://%s%s/%s/ping", s.Host[botInstance], s.Listen, botId)
   resp, err = httpclient.Get(url)

   if err != nil {
      Goose.Logf(1,"%s (%s) %#v",ErrFailedPingingBot,err,resp)
      return ErrFailedPingingBot
   }

   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      Goose.Logf(1,"%s %s at %s (status code=%d)",ErrFailedPingingBot,botId,url,resp.StatusCode)
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
