package masterbot

func New(botId string, appcfg *ConfigT, onStop func()) *ServiceT {
   return &ServiceT{
      botId: botId,
      onStop: onStop,
      appcfg: appcfg,
   }
}

