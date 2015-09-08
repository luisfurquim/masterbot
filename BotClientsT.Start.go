package masterbot

import (
   "os"
   "encoding/json"
)


func (bc *BotClientsT) Start(config *ConfigT, debugLevel int) {
   Goose.Logf(2,"Registering ping job [%s]",config.BotPingRate)
   Kairos.AddFunc(config.BotPingRate, (func(bots *BotClientsT) (func()) {
      return func() {
         var botId       string
         var botCfg      BotClientT
         var botInstance int

         Goose.Logf(2,"Pinging slave bots")

         for botId, botCfg = range *bots {
            for botInstance,_ = range botCfg.Host {
               go func(id string, instance int, cfg BotClientT) {
                  var err        error
                  var botCfgFile []byte

                  cfg.PageNotFound = config.PageNotFound
                  cfg.Pem          = config.Pem
                  cfg.BinDir       = config.BinDir

                  botCfgFile, err = json.Marshal(cfg)
                  if err != nil {
                     Goose.Logf(1,"Error marshaling botconfig (%s)",err)
                     os.Exit(1)
                  }

                  err = cfg.Start(id,instance,botCfgFile,config,debugLevel)
                  if err != nil {
                     Goose.Logf(1,"Error starting bot %s (%s)",id,err)
                  }
               }(botId,botInstance,botCfg)
            }
         }
      }
   })(bc)) // Closure to avoid direct access to bc and having it changing from time to time
}

