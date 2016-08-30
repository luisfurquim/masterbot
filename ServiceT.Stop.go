package masterbot

import (
   "fmt"
   "time"
   "sync"
   "net/http"
   "github.com/luisfurquim/stonelizard"
)

func (svc ServiceT) Stop() stonelizard.Response {
   var botId       string
   var botCfg     *BotClientT
   var wg          sync.WaitGroup
   var botInstance int
   var host        Host

   Goose.StartStop.Logf(2,"Stopping kairos")
   Kairos.Stop()

   Goose.StartStop.Logf(2,"Stopping slave bots")

   for botId, botCfg = range svc.appcfg.Bot {
      wg.Add(len(botCfg.Host)) // wait the stop of each instance of the slavebots
   }

   for botId, botCfg = range svc.appcfg.Bot {
      for botInstance, host = range botCfg.Host {
         go func(id string, instance int, cfg *BotClientT, h Host) {
            var err    error
            var url    string
            var resp  *http.Response

            defer wg.Done()

            if h.Status == BotStatRunning {
               url   = fmt.Sprintf("https://%s%s/%s/stop", h.Name, cfg.Listen, id)
               Goose.StartStop.Logf(2,"Stopping bot %s@%s via %s",id,h.Name,url)
               resp, err = svc.appcfg.HttpsStopClient.Get(url)

               if err != nil {
                  Goose.StartStop.Logf(1,"Error stopping bot %s@%s (%s)",id,h.Name,err)
                  return
               }

               if resp.StatusCode != http.StatusNoContent {
                  Goose.StartStop.Logf(1,"Error of status code stopping bot %s@%s (%s)",id,h.Name,resp.Status)
               }
            }

         }(botId,botInstance,botCfg,host)
      }
   }

   wg.Wait()

   wg.Add(1) // wait the stop of the masterbot listener

   go (func () {
      Goose.StartStop.Logf(2,"Stopping listener")
      wg.Wait()
      svc.onStop()
      Goose.StartStop.Logf(2,"Masterbot listener finished")
   })()

   time.Sleep(100 * time.Millisecond)

   Goose.StartStop.Logf(2,"Ending masterbot")

   defer wg.Done()

   return stonelizard.Response{
      Status: http.StatusNoContent,
   }
}

