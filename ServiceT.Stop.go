package masterbot

import (
   "fmt"
   "sync"
   "time"
   "net/http"
   "github.com/luisfurquim/stonelizard"
)

func (svc ServiceT) Stop() stonelizard.Response {
   var botId       string
   var botCfg      BotClientT
   var wg          sync.WaitGroup
   var botInstance int

   Goose.Logf(2,"Stopping slave bots")

   for botId, botCfg = range svc.appcfg.Bot {
      for botInstance, _ = range botCfg.Host {
         wg.Add(len(botCfg.Host)) // stop each instance of the slavebots
      }
   }
   wg.Add(1) // stop the masterbot itself

   for botId, botCfg = range svc.appcfg.Bot {
      for botInstance, _ = range botCfg.Host {
         go func(id string, instance int, cfg BotClientT) {
            var err    error
            var url    string
            var htcli *http.Client
            var resp  *http.Response

            defer wg.Done()

            htcli = svc.appcfg.HttpsClient(5 * time.Second)
            url   = fmt.Sprintf("https://%s%s/%s/stop", botCfg.Host[instance], botCfg.Listen, id)
            Goose.Logf(2,"Stopping bot %s@%s via %s",id,botCfg.Host[instance],url)
            resp, err = htcli.Get(url)

            if err != nil {
               Goose.Logf(1,"Error stopping bot %s@%s (%s)",id,botCfg.Host[instance],err)
               return
            }

            if resp.StatusCode != http.StatusOK {
               Goose.Logf(1,"Error status stopping bot %s@%s (%s)",id,botCfg.Host[instance],resp.Status)
            }

         }(botId,botInstance,botCfg)
      }
   }

   Goose.Logf(2,"Stopping masterbot")

   go (func () {
      Kairos.Stop()
      wg.Wait()
      svc.onStop()
   })()

   defer wg.Done()

   return stonelizard.Response{
      Status: http.StatusOK,
      Body: "OK",
   }
}

