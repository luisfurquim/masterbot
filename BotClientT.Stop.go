package masterbot

import (
   "fmt"
   "time"
   "net/http"
)

func (svc *BotClientT) Stop(botId string, cfg *ConfigT, botInstance int) error {
   var err     error
   var errHist error
   var url     string
   var htcli  *http.Client
   var resp   *http.Response
   var i       int

   Goose.StartStop.Logf(2,"Stopping slave bot %s",botId)

   htcli = cfg.HttpsClient(time.Duration(0))
	if botInstance < 0 {
		svc.Status = BotStatStopped
		if svc.OnStatUpdate != nil {
			svc.OnStatUpdate(BotStatStopped)
			Goose.StartStop.Logf(2, "Mudou para S /sophia/bot/%s/status\n", botId)
		}
	}
   for i, _ = range svc.Host{
      if (botInstance<0) || (botInstance==i) {
         if svc.Host[i].Status == BotStatRunning {
            
            svc.Host[i].Status = BotStatStopped
            
            if svc.Host[i].OnStatUpdate != nil {
               svc.Host[i].OnStatUpdate(BotStatStopped)
               Goose.StartStop.Logf(2, "Alterou status para S no bot %s@%s \n", botId, svc.Host[i].Name)
            }
            
            url   = fmt.Sprintf("https://%s%s/%s/stop", svc.Host[i].Name, svc.Listen, botId)
            Goose.StartStop.Logf(2,"Stopping bot %s via %s",botId,url)
            resp, err = htcli.Get(url)

            if err != nil {
               Goose.StartStop.Logf(1,"%s %s (%s)",ErrStoppingBot,botId,err)
               errHist = ErrStoppingBot
               continue
            }
            if resp.StatusCode != http.StatusNoContent {
               Goose.StartStop.Logf(1,"%s %s (%s)",ErrStatusStoppingBot,botId,resp.Status)
               errHist = ErrStatusStoppingBot
               continue
            }
			
            
         }
      }
   }
   return errHist
}
