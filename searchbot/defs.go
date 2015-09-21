package searchbot


import (
   "golang.org/x/tools/container/intsets"
   "github.com/luisfurquim/masterbot"
)


type ProviderT struct {  // Provider is a search operation. Bots may offer more than 1 provider
   Bot            *masterbot.BotClientT
   Path            string    // Restful path of the service
   HttpMethod      string    // GET, POST, PUT, DELET, etc.
   Operation      *SwaggerOperationT  // Service specification details
   Requires []string         // Input parameters
   Provides []string         // Return values
}

type SearchBotT struct {
   Providers              []ProviderT
   ByProvision   map[string]intsets.Sparse
   ByRequirement map[string]intsets.Sparse
}

type BotClientsT masterbot.BotClientsT

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


