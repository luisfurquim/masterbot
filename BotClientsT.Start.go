package masterbot

import (
   "regexp"
)


func (bc *BotClientsT) Start(config *ConfigT, cmdline string, debugLevel int) {
   var err           error
   var botId         string
   var botInstance   int
   var botCfg       *BotClientT

   Goose.Ping.Logf(2,"Registering ping jobs [%s]",config.BotPingRate)

   for botId, botCfg = range *bc {
      if botCfg.Status == BotStatPaused {
         continue
      }

      if (botCfg.SearchPath != "") && (botCfg.SearchPathRE==nil) {
         botCfg.SearchPathRE = regexp.MustCompile(botCfg.SearchPath)
      }

      botCfg.CronPingId = make([]int,len(botCfg.Host))
      botCfg.CronPingFn = make([]func(),len(botCfg.Host))
      (*bc)[botId]      = botCfg

      for botInstance,_ = range botCfg.Host {
         Goose.Ping.Logf(3,"Agendando instancia %s (%d) de pinger, |pingId|=%d",botCfg.Host[botInstance].Name,botInstance,len(botCfg.CronPingId))
         botCfg.CronPingFn[botInstance] = (func(bot *BotClientT, cmd string, id string, instance int) (func()) {
            return func() {
               var err        error
					//Testing botCfg.Status==BotStatStopped to prevent start bot from cron scheduling.
               if botCfg.Status ==  BotStatPaused || botCfg.Status == BotStatStopped {
                  return
               }

               err = bot.Start(id,instance,cmd,config,debugLevel)
               if err != nil {
                  Goose.StartStop.Logf(1,"Error starting bot %s@%s (%s)",id,bot.Host[instance].Name,err)
               }
            }
         })(botCfg,cmdline,botId,botInstance) // Closure to avoid direct access to bc and having it changing from time to time
         botCfg.CronPingId[botInstance], err = Kairos.AddFunc(config.BotPingRate, botCfg.CronPingFn[botInstance])
         if err != nil {
            Goose.Ping.Logf(1,"Error scheduling bot %s@%s ping job (%s)",botId,botCfg.Host[botInstance],err)
         } else {
            go botCfg.CronPingFn[botInstance]()
         }
      }

      botCfg.Status = BotStatRunning
      if botCfg.OnStatUpdate != nil {
         botCfg.OnStatUpdate(BotStatRunning)
      }

   }
}

