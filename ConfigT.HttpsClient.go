package masterbot

import(
   "net"
   "time"
   "net/http"
   "crypto/tls"
)

func (cfg ConfigT) HttpsClient(tmout time.Duration) *http.Client {
   var httpclient  http.Client
   var tlsConfig  *tls.Config

   tlsConfig = &tls.Config{
      Certificates: []tls.Certificate{cfg.ClientCert},
      RootCAs:      cfg.ClientCA,
      InsecureSkipVerify: true,
   }
   
  	Goose.Ping.Logf(3, "Construindo NametoCertificate \n")
   tlsConfig.BuildNameToCertificate()

   Goose.Ping.Logf(3, "Definicao do httpclient apos definicao tlsCfg\n") 
   /*
   httpclient = http.Client{
     Transport: &http.Transport{
        TLSClientConfig:     tlsConfig,
        DisableCompression:  true,
        Dial: (&net.Dialer{
           Timeout:   tmout,
           KeepAlive: 30 * time.Second,
        }).Dial,
        TLSHandshakeTimeout: 10 * time.Second,
     },
     Timeout:	tmout,
  }
   
*/
   httpclient = http.Client{
      Transport: &http.Transport{
         TLSClientConfig:     tlsConfig,
         DisableCompression:  true,
      },
   }

   if tmout > time.Duration(0) {
      httpclient.Transport.(*http.Transport).Dial = func (network, addr string) (net.Conn, error) {
         return net.DialTimeout(network, addr, tmout)
      }
   }

	Goose.Ping.Logf(6, "HttpsClient=%p, %#v, %T", &httpclient, &httpclient, httpclient)
   return &httpclient
}