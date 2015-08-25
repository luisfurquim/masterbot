package masterbot

import (
   "fmt"
   "sync"
   "time"
   "net/http"
   "github.com/luisfurquim/stonelizard"
)

func (svc ServiceT) Stop() stonelizard.Response {
   var botId      string
   var botCfg     BotClientT
   var wg         sync.WaitGroup

   Goose.Logf(2,"Stopping slave bots")

   wg.Add(len(svc.appcfg.Bot) + 1)

   for botId, botCfg = range svc.appcfg.Bot {
      go func(id string, cfg BotClientT) {
         var err    error
         var url    string
         var htcli *http.Client
         var resp  *http.Response

         defer wg.Done()

         htcli = svc.appcfg.HttpsClient(5 * time.Second)
         url   = fmt.Sprintf("https://%s%s/%s/stop", botCfg.Host, botCfg.Listen, id)
         Goose.Logf(2,"Stopping bot %s via %s",id,url)
         resp, err = htcli.Get(url)

         if err != nil {
            Goose.Logf(1,"Error stopping bot %s (%s)",id,err)
            return
         }

         if resp.StatusCode != http.StatusOK {
            Goose.Logf(1,"Error status stopping bot %s (%s)",id,resp.Status)
         }

      }(botId,botCfg)
   }

   Goose.Logf(2,"Stopping sherlock bot")

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

