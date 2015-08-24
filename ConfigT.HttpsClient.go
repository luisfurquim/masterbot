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
   tlsConfig.BuildNameToCertificate()

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

   return &httpclient
}