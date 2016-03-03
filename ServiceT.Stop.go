package masterbot

import (
   "fmt"
   "sync"
   "net/http"
   "github.com/luisfurquim/stonelizard"
)

func (svc ServiceT) Stop() stonelizard.Response {
   var botId       string
   var botCfg      BotClientT
   var wg          sync.WaitGroup
   var botInstance int

   Goose.StartStop.Logf(2,"Stopping slave bots")

   for botId, botCfg = range svc.appcfg.Bot {
      for botInstance, _ = range botCfg.Host {
         wg.Add(len(botCfg.Host)) // wait the stop of each instance of the slavebots
      }
   }
   wg.Add(1) // wait the stop of the masterbot itself

   for botId, botCfg = range svc.appcfg.Bot {
      for botInstance, _ = range botCfg.Host {
         go func(id string, instance int, cfg BotClientT) {
            var err    error
            var url    string
            var resp  *http.Response

            defer wg.Done()

            url   = fmt.Sprintf("https://%s%s/%s/stop", botCfg.Host[instance], botCfg.Listen, id)
            Goose.StartStop.Logf(2,"Stopping bot %s@%s via %s",id,botCfg.Host[instance],url)
            resp, err = svc.appcfg.HttpsStopClient.Get(url)

            if err != nil {
               Goose.StartStop.Logf(1,"Error stopping bot %s@%s (%s)",id,botCfg.Host[instance],err)
               return
            }

            if resp.StatusCode != http.StatusNoContent {
               Goose.StartStop.Logf(1,"Error of status code stopping bot %s@%s (%s)",id,botCfg.Host[instance],resp.Status)
            }

         }(botId,botInstance,botCfg)
      }
   }

   Goose.StartStop.Logf(2,"Stopping masterbot")

   go (func () {
      Kairos.Stop()
      wg.Wait()
      svc.onStop()
   })()

   defer wg.Done()

   return stonelizard.Response{
      Status: http.StatusNoContent,
   }
}

