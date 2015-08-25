package masterbot


import (
   "time"
   "crypto/tls"
   "crypto/x509"
   "golang.org/x/crypto/ssh"
)


type QueryT struct {
   Resource string
   SearchBy []interface{}
   SearchFor []interface{}
}

/*
type BotClientT interface {
   SetConfig(io.Reader) error
   Ping() error
   Auth(user, pw string) error
   Restart() error
   Stop() error
   ListResources() ([]ResourceT,error)
   Query(qry QueryT) (interface{},error)
}
*/

type BotClientT struct {
   PageNotFound string       `json:"pageNotFound"`
   Pem          string       `json:"pem"`
   BinDir       string       `json:"bindir"`
   BinName      string       `json:"binname"`
   Listen       string       `json:"listen"`
   CrlListen    string       `json:"crllisten"`
   Host         string       `json:"host"`
   SysUser      string       `json:"sysuser"`
   Status       uint8        `json:"status"`
   Config       interface{}  `json:"config"`
}

type BotClientsT map[string]BotClientT

type ConfigT struct {
   Id               string           `json:"id"`
   Host             string           `json:"host"`
   SysUser          string           `json:"sysuser"`
   WorkDir          string           `json:"workdir"`
   Listen           string           `json:"listen"`
   CrlListen        string           `json:"crllisten"`
   PageNotFound     string           `json:"pageNotFound"`
   Pem              string           `json:"pem"`
   BinDir           string           `json:"bindir"`
   BinName          string           `json:"binname"`
   ClientCert       tls.Certificate
   ClientCA        *x509.CertPool
   Bot              BotClientsT      `json:"bot"`
   SshClientConfig *ssh.ClientConfig
   BotPingRate      string           `json:"botpingrate"`
   BotPingTimeout   time.Duration    `json:"botpingtimeout"`
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

