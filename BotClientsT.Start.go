package masterbot

import (
   "regexp"
   "encoding/json"
)


func (bc *BotClientsT) Start(config *ConfigT, debugLevel int) {
   var err           error
   var botId         string
   var botInstance   int
   var botCfg        BotClientT
   var botCfgFile  []byte

   Goose.Logf(2,"Registering ping jobs [%s]",config.BotPingRate)

   for botId, botCfg = range *bc {
      if (botCfg.SearchPath != "") && (botCfg.SearchPathRE==nil) {
         botCfg.SearchPathRE = regexp.MustCompile(botCfg.SearchPath)
      }

      botCfg.PageNotFound = config.PageNotFound
      botCfg.Pem          = config.Pem
      botCfg.BinDir       = config.BinDir

      botCfgFile, err = json.Marshal(botCfg)
      if err != nil {
         Goose.Logf(1,"Error marshaling botconfig of %s@%s (%s)",botId,botCfg.Host[botInstance],err)
         continue
      }

      for botInstance,_ = range botCfg.Host {
         Kairos.AddFunc(config.BotPingRate, (func(bot *BotClientT, configFile string, id string, instance int) (func()) {
            return func() {
               var err        error
//               Goose.Logf(2,"Pinging slave bots")

               err = bot.Start(id,instance,[]byte(configFile),config,debugLevel)
               if err != nil {
                  Goose.Logf(1,"Error starting bot %s@%s (%s)",id,bot.Host[instance],err)
               }
            }
         })(&botCfg,string(botCfgFile),botId,botInstance)) // Closure to avoid direct access to bc and having it changing from time to time
      }
   }
}

