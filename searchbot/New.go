package searchbot

import (
   "fmt"
   "time"
   "strings"
   "net/http"
   "encoding/json"
   "github.com/luisfurquim/masterbot"
   "github.com/luisfurquim/stonelizard"
   "golang.org/x/tools/container/intsets"
)

func New(cfg *masterbot.ConfigT) (*SearchBotT, error) {
   var err           error
   var ok            bool
   var res           SearchBotT
   var host          masterbot.Host
   var bot           string
   var botCfg       *masterbot.BotClientT
   var swagger       stonelizard.SwaggerT
   var resp         *http.Response
   var path          string
   var service       stonelizard.SwaggerPathT
   var provider      ProviderT
   var opParm        stonelizard.SwaggerParameterT
   var method        string
   var operation     stonelizard.SwaggerOperationT
   var httpStatus    string
   var opResp        stonelizard.SwaggerResponseT
   var prop          string
   var newProviderId int
   var fieldId       int

   Goose.Taxonomy.Logf(2,"Starting searchbot indexing")

   res = SearchBotT{
      Providers:            []ProviderT{},
      Taxonomy:               TaxonomyTreeT{Rune: 0, Id: 0, Next: []*TaxonomyTreeT{}},
      ByProvision:   map[int]*intsets.Sparse{},
      ByRequirement: map[int]*intsets.Sparse{},
      Config:                 cfg,
      HttpsSearchClient:      cfg.HttpsClient(time.Duration(cfg.BotCommTimeout) * time.Second),
   }

   time.Sleep(2 * time.Duration(cfg.BotCommTimeout) * time.Second)

   masterbot.Kairos.AddFunc(cfg.BotPingRate, (func(bots *masterbot.BotClientsT) (func()) {
      return func() {
         res.Lock()
         Goose.Taxonomy.Logf(2,"Refreshing searchbot database")
         res.Providers     = []ProviderT{}
         res.ByProvision   = map[int]*intsets.Sparse{}
         res.ByRequirement = map[int]*intsets.Sparse{}

         for bot, botCfg = range cfg.Bot {
            for _, host = range botCfg.Host {
               err = func () error {
                  var url        string
                  var err        error

                  url   = fmt.Sprintf("https://%s%s/swagger.json", host.Name, botCfg.Listen)
                  Goose.Taxonomy.Logf(2,"fetching swagger.json via %s",url)

                  resp, err = res.HttpsSearchClient.Get(url)

                  if resp != nil {
                     defer resp.Body.Close()
                  }

                  if err != nil {
                     Goose.Taxonomy.Logf(1,"%s from %s@%s (%s)",ErrTmoutFetchingSwagger,bot,host.Name,err)
                     return ErrTmoutFetchingSwagger
                  }

                  if resp.StatusCode != http.StatusOK {
                     Goose.Taxonomy.Logf(1,"%s from %s@%s at %s (status code=%d)",ErrHttpStatusFetchingSwagger,bot,host.Name,url,resp.StatusCode)
                     return ErrHttpStatusFetchingSwagger
                  }

                  err = json.NewDecoder(resp.Body).Decode(&swagger)
                  if err != nil {
                     Goose.Taxonomy.Logf(1,"%s of %s@%s (%s)",ErrDecodingSwagger,bot,host.Name,err)
                     return ErrDecodingSwagger
                  }

                  Goose.Taxonomy.Logf(3,"fetched swagger.json")
                  return nil
               }()

               if err != nil {
                  continue
               }

               Goose.Taxonomy.Logf(2,"fetched swagger.json no error")

               for path, service = range swagger.Paths {
                  for method, operation = range service {
                     Goose.Taxonomy.Logf(3,"swagger method: %s, op: %s",method,operation.OperationId)
                     if ((len(path)>=8) && (path[:8]=="/search/")) || ((botCfg.SearchPathRE!=nil) && (botCfg.SearchPathRE.MatchString(path))) || HasSearchTag(operation.Tags) {
                        Goose.Taxonomy.Logf(3,"Found search operation: %s, path=%s",operation.OperationId, path)
                        Goose.Taxonomy.Logf(7,"operation parameters: %#v",operation.Parameters)
                        provider = ProviderT{
                           Bot:            botCfg,
                           BotId:          bot,
                           Path:           swagger.BasePath + path,
                           HttpMethod:     strings.ToUpper(method),
                           Operation:      operation,
                           Requires:      &intsets.Sparse{},
                           Provides:      &intsets.Sparse{},
                        }
                        for httpStatus, opResp = range operation.Responses {
                           if httpStatus[0] == '2' {
                              Goose.Taxonomy.Logf(7,"Testing response: %s",httpStatus)
                              if opResp.Schema != nil {
                                 for prop, _ = range opResp.Schema.Properties {
                                    Goose.Taxonomy.Logf(7,"found prop: %s",prop)
                                    fieldId = res.Taxonomy.Add(prop)
      //                              _, _, pTmp := res.Taxonomy.Search(prop)
      //                              Goose.Taxonomy.Logf(4,"Added Taxonomy: %s as %d, search reports it is %d",prop,fieldId,pTmp.Id)
                                    Goose.Taxonomy.Logf(7,"Taxonomy: %s",res.Taxonomy)
                                    provider.Provides.Insert(fieldId)
                                 }
                              }
                              if !provider.Provides.IsEmpty() {
                                 Goose.Taxonomy.Logf(7,"found response: %s",httpStatus)
                                 break
                              }
                           }
                        }
                        if !provider.Provides.IsEmpty() {
                           newProviderId = len(res.Providers)
                           for _, opParm = range operation.Parameters {
                              Goose.Taxonomy.Logf(7,"Registering taxonomy: %s",opParm.Name)
                              fieldId = res.Taxonomy.Add(opParm.Name)
                              Goose.Taxonomy.Logf(3,"Registered taxonomy: fieldId:%d(%s)",fieldId,opParm.Name)
                              provider.Requires.Insert(fieldId)
                              Goose.Taxonomy.Logf(7,"Inserted provider")
                              if _, ok = res.ByRequirement[fieldId]; !ok {
                                 res.ByRequirement[fieldId] = &intsets.Sparse{}
                              }
                              res.ByRequirement[fieldId].Insert(newProviderId)
                              Goose.Taxonomy.Logf(7,"Indexed provider")
                           }
                           Goose.Taxonomy.Logf(3,"done taxonomy")
                           for fieldId=0; fieldId<provider.Provides.Max(); fieldId++ {
                              if provider.Provides.Has(fieldId) {
                                 Goose.Taxonomy.Logf(7,"Indexing fieldId: %d",fieldId)
                                 if _, ok := res.ByProvision[fieldId]; !ok {
                                    res.ByProvision[fieldId] = &intsets.Sparse{}
                                    Goose.Taxonomy.Logf(7,"created index for fieldId=%d",fieldId)
                                 }
                                 Goose.Taxonomy.Logf(7,"Indexed fieldId=%d, res.ByProvision:%#v",fieldId,res.ByProvision)
                                 res.ByProvision[fieldId].Insert(newProviderId)
                                 Goose.Taxonomy.Logf(7,"Indexed fieldId=%d -> newProviderId:%d",fieldId,newProviderId)
                                 Goose.Taxonomy.Logf(7,"ByProvision['%d']: %s",fieldId,res.ByProvision[fieldId])
                              }
                           }
                           Goose.Taxonomy.Logf(3,"done provider")
                           Goose.Taxonomy.Logf(4,"Adding provider: %#v, requires: %#v, provides: %#v",provider,*provider.Requires,*provider.Provides)
                           res.Providers = append(res.Providers,provider)
                        }
                        Goose.Taxonomy.Logf(3,"End provider")
                     }
                  }
               }

               Goose.Taxonomy.Logf(2,"end registering search host: %s",host.Name)
               Goose.Taxonomy.Logf(4,"Taxonomy: %#v",res.Taxonomy)
               break
            }
         }
         res.Unlock()
      }
   })(&cfg.Bot))

   // Timeout was set just because the bots must answer quickly when we ask for its definitions
   // but to do actual services, they will probably need much more time
   res.HttpsSearchClient = cfg.HttpsClient(0)
   return &res, nil
}

