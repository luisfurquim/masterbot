package masterbot


import (
   "time"
   "regexp"
   "net/http"
   "crypto/tls"
   "crypto/x509"
   "golang.org/x/crypto/ssh"
   "github.com/luisfurquim/stonelizard"
)

type StatusT struct {
   Status       string             `json:"status"`
   OnStatUpdate func(string) error `json:"-"`
}


type Host struct {
   Name    string        `json:"name"`
   Port    string        `json:"port"` // Used when the SSH port is non standard
   StatusT
}

type Hosts []Host

type BotClientT struct {
   BinDir        string        `json:"bindir"`
   BinName       string        `json:"binname"`
   Listen        string        `json:"listen"`
   Host          Hosts         `json:"host"`
   SysUser       string        `json:"sysuser"`
   WorkDir       string        `json:"workdir"`
   SearchPath    string        `json:"searchpath"`
   SearchPathRE *regexp.Regexp `json:"-"`
   StatusT
   CronPingId  []int           `json:"-"`
   CronPingFn  []func()        `json:"-"`
}

type BotClientPtr           *BotClientT
type BotClientsT map[string]BotClientPtr

type Timeout time.Duration

type ConfigT struct {
   Id               string            `json:"id"`
   Host             Hosts             `json:"host"`
   SysUser          string            `json:"sysuser"`
   WorkDir          string            `json:"workdir"`
   Listen           string            `json:"listen"`
   CrlListen        string            `json:"crllisten"`
   PageNotFoundPath string            `json:"pageNotFound"`
   Pem              string            `json:"pem"`
   BinDir           string            `json:"bindir"`
   BinName          string            `json:"binname"`
   ClientCert       tls.Certificate   `json:"-"`
   ClientCA        *x509.CertPool     `json:"-"`
   Bot              BotClientsT       `json:"bot,omitempty"`
   SshClientConfig *ssh.ClientConfig  `json:"-"`
   BotPingRate      string            `json:"botpingrate,omitempty"`
   BotCommTimeout   Timeout           `json:"botcommtimeout,omitempty"`
   HttpsPingClient *http.Client       `json:"-"`
   HttpsStopClient *http.Client       `json:"-"`
   Certkit          stonelizard.AuthT `json:"-"`
   Transport		 *http.Transport		`rset:"-"`
}

const (
   BotStatStopped     string = "S"
   BotStatRunning     string = "R"
   BotStatPaused      string = "P" //Used to disable
   BotStatUnreachable string = "U"
)

const SSHPort string = "22"

type ServiceT struct {
   // end point for monitoring
   stop    bool `method:"GET" path:"/stop" ok:"Ends service operation"`

   botId   string
   onStop  func()
   appcfg *ConfigT
}

