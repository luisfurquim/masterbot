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
   BotId           string
   Path            string    // Restful path of the service
   HttpMethod      string    // GET, POST, PUT, DELETE, etc.
   Operation       stonelizard.SwaggerOperationT  // Service specification details
   Requires       *intsets.Sparse                 // Input parameters (set of Taxonomy entries)
   Provides       *intsets.Sparse                 // Return values (set of Taxonomy entries)
}

type ResponseFieldT struct {
   Value    interface{}
   Source   string
   DtUpd    string
}

type SearchBotT struct {
   // defines the root of this service, and its meta data.
   root stonelizard.Void `root:"/searchbot/" consumes:"application/json" produces:"application/json" allowGzip:"true" enableCORS:"true"`

   // defines global information about the service
   info    stonelizard.Void `title:"SearchBot" description:"Wrapper robot for aggregation of search bots" tos:"Free to use under the terms of MPL2.0" version:"0.1"`
   contact stonelizard.Void `contact:"Luis Otávio de Colla Furquim" url:"http://www.prrs.mpf.mp.br" email:"vuco@mpf.mp.br"`
   license stonelizard.Void `license:"MPL2.0" url:"https://www.mozilla.org/en-US/MPL/2.0/"`
//   extDoc  stonelizard.Void `url:""`

   // end point for querying the wrapper
   search map[string][]ResponseFieldT `method:"GET" path:"/search" header:"X-Login,X-Password" query:"searchBy,searchFor" ok:"Query succesful" X-Login:"Login name" X-Password:"User's password" searchBy:"Key-Value pairs of input parameters, key names must follow the taxonomy" searchFor:"List of fields to retrieve, field names must follow the taxonomy"`

   // end point for dumping the taxonomy
   taxonomyDump []string `method:"GET" path:"/taxonomy" ok:"Query succesful"`

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



