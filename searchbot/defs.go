package searchbot


import (
   "sync"
   "net/http"
   "github.com/luisfurquim/masterbot"
   "github.com/luisfurquim/stonelizard"
   "golang.org/x/tools/container/intsets"
)


type TaxonomyTreeT struct {
   Rune    rune
   Id      int             // In the root, this is actually the next available id
   Next []*TaxonomyTreeT
}

// Provider is a search operation. Bots may offer more than 1 provider
type ProviderT struct {
   Bot            *masterbot.BotClientT
   Path            string    // Restful path of the service
   HttpMethod      string    // GET, POST, PUT, DELETE, etc.
   Operation       stonelizard.SwaggerOperationT  // Service specification details
   Requires       *intsets.Sparse                 // Input parameters (set of Taxonomy entries)
   Provides       *intsets.Sparse                 // Return values (set of Taxonomy entries)
}

type SearchBotT struct {
   sync.RWMutex
   Providers            []ProviderT
   Taxonomy               TaxonomyTreeT     // List of all known data provided/required by the bots
   ByProvision   map[int]*intsets.Sparse    // Index of bots by which data they provide
   ByRequirement map[int]*intsets.Sparse    // Index of bots by which data they require
   Config                *masterbot.ConfigT
   HttpsSearchClient     *http.Client
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


