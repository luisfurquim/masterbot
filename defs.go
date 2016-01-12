package masterbot


import (
   "time"
   "regexp"
   "net/http"
   "crypto/tls"
   "crypto/x509"
   "golang.org/x/crypto/ssh"
   "github.com/luisfurquim/stonelizard/certkit"
)


type BotClientT struct {
//   PageNotFound  string        `json:"pageNotFound"`
//   Pem           string        `json:"pem"`
   BinDir        string        `json:"bindir"`
   BinName       string        `json:"binname"`
   Listen        string        `json:"listen"`
//   CrlListen     string        `json:"crllisten"`
   Host        []string        `json:"host"`
//   ThisHost      string        `json:"thishost,omitempty"`
   SysUser       string        `json:"sysuser"`
   SearchPath    string        `json:"searchpath"`
   SearchPathRE *regexp.Regexp `json:"searchpath"`
   Status        uint8         `json:"status"`
//   Config        interface{}   `json:"config"`
   CronPingId  []int
   CronPingFn  []func()
}

type BotClientsT map[string]BotClientT

type ConfigT struct {
   Id               string           `json:"id"`
   Host           []string           `json:"host"`
   SysUser          string           `json:"sysuser"`
   WorkDir          string           `json:"workdir"`
   Listen           string           `json:"listen"`
   CrlListen        string           `json:"crllisten"`
   PageNotFoundPath string           `json:"pageNotFound"`
   Pem              string           `json:"pem"`
   BinDir           string           `json:"bindir"`
   BinName          string           `json:"binname"`
   ClientCert       tls.Certificate
   ClientCA        *x509.CertPool
   Bot              BotClientsT      `json:"bot"`
   SshClientConfig *ssh.ClientConfig
   BotPingRate      string           `json:"botpingrate"`
   BotCommTimeout   time.Duration    `json:"botcommtimeout"`
   HttpsPingClient *http.Client
   HttpsStopClient *http.Client
   Certkit         *certkit.CertKit
}

const (
   BotStatStopped = iota
   BotStatRunning
   BotStatPaused
   BotStatUnreachable
)


type ServiceT struct {
   // end point for monitoring
   stop    bool `method:"GET" path:"/stop" ok:"Ends service operation"`

   botId   string
   onStop  func()
   appcfg *ConfigT
}

