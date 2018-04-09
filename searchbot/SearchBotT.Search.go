package searchbot

import (
   "fmt"
//   "time"
   "bytes"
   "strings"
   "net/url"
   "net/http"
   "io/ioutil"
   "encoding/json"
   "encoding/xml"
//   "github.com/luisfurquim/masterbot"
   "github.com/luisfurquim/stonelizard"
//   "github.com/luisfurquim/goose"
   "golang.org/x/tools/container/intsets"
)

func (sb *SearchBotT) Search(searchBy map[string]string, searchFor []string, login string, password string, tid string) stonelizard.Response {
   var providers              *intsets.Sparse
   var searchFields           *intsets.Sparse
   var commonFields           *intsets.Sparse
   var oneShotProviders       *intsets.Sparse
   var field                   string
   var isFragmented            bool
   var i                       int
   var p                      *TaxonomyTreeT
   var hasQueryParm            bool
   var rep     chan map[string]ResponseFieldT
   var response     map[string]ResponseFieldT
   var responses  map[string][]ResponseFieldT
   var respCount               int

   //TODO: Readaptar
   searchBy["X-Login"]    = login
   searchBy["X-Password"] = password
   searchBy["X-Trackid"]  = tid

//   Goose.Search = goose.Alert(6)
//   defer func() { Goose.Search = goose.Alert(2) }()


   providers = &intsets.Sparse{}

   Goose.Search.Logf(2,"TID:%s len(sb.Providers): %d",tid,len(sb.Providers))

   // Fill the providers set with all provider currently known
   for i=0; i < len(sb.Providers); i++ {
      providers.Insert(i)
   }

   // Determine if there is at least one bot providing all data needed
   // by repeatedly computing providers ∩= 'providers of a given field'
   for _, field = range searchFor {
      Goose.Search.Logf(0,"TID:%s will constrain by %s", tid, field)
      i, _, p = sb.Taxonomy.Search(field)
      Goose.Search.Logf(0,"TID:%s field %s has id=%d", tid, field, p.Id)
      if ((i+1)!=len(field)) || (p==nil) || (p.Id<0) {
         Goose.Search.Logf(1,"TID:%s %s: %s", tid, ErrUndefinedField, field)
         return stonelizard.Response{
            Status: http.StatusInternalServerError,
            Body: fmt.Sprintf("%s: %s",ErrUndefinedField, field),
         }
      }
      if sb.ByProvision[p.Id] == nil {
         isFragmented = true
         break
      }
      providers.IntersectionWith(sb.ByProvision[p.Id])
      if providers.IsEmpty() {
         isFragmented = true
         break
      }
      Goose.Search.Logf(0,"TID:%s constraining by %s gives %s", tid, field, providers)
   }

   Goose.Search.Logf(4,"TID:%s Determined if there is at least one bot providing all data needed (isFragmented=%#v): %#v", tid, isFragmented, providers.String())

   if !isFragmented {
      // Select in the bots that have all information needed
      // those who require only information we have
      searchFields = &intsets.Sparse{}
      for field, _ = range searchBy {
         Goose.Search.Logf(6,"TID:%s sb.Taxonomy.Search(%s)", tid, field)
         i, _, p = sb.Taxonomy.Search(field)
         Goose.Search.Logf(6,"TID:%s i=%d, p=%#v", tid, i, p)
         if p != nil {
            Goose.Search.Logf(3,"TID:%s Selecting new search field %s with id=%d", tid, field, p.Id)
            searchFields.Insert(p.Id)
         }
      }

      Goose.Search.Logf(4,"TID:%s Bitstring of search created: %#v", tid, searchFields.String())

      oneShotProviders = &intsets.Sparse{}
      commonFields     = &intsets.Sparse{}

      Goose.Search.Logf(4,"TID:%s providers.Max(): %d", tid, providers.Max())
      for i=0; i <= providers.Max(); i++ {
         Goose.Search.Logf(4,"TID:%s sb.Providers[%d].Requires: %s",tid,i,sb.Providers[i].Requires.String())
         commonFields.Intersection(searchFields,sb.Providers[i].Requires)
         if commonFields.Len() == sb.Providers[i].Requires.Len() {
            Goose.Search.Logf(3,"TID:%s Intersection at %d",tid,i)
            oneShotProviders.Insert(i)
         }
      }

      Goose.Search.Logf(4,"TID:%s oneShotProviders: %s",tid,oneShotProviders.String())

      // If there is at least one bot who gives all fields
      // we need and requires just fields we already have...
      if oneShotProviders.Len() > 0 {
         rep = make(chan map[string]ResponseFieldT,oneShotProviders.Len())
         Goose.Search.Logf(4,"TID:%s len(sb.Providers): %d",tid,len(sb.Providers))
         for i=0; i <= oneShotProviders.Max(); i++ {
            Goose.Search.Logf(4,"TID:%s oneShotProvider: %d",tid,i)
            if oneShotProviders.Has(i) {
               go func(instance int, report chan map[string]ResponseFieldT) {
                  var err                    error
                  var req                   *http.Request
                  var host                   string
                  var path                   string
                  var swParm                 stonelizard.SwaggerParameterT
                  var body        map[string]interface{}
                  var b_body               []byte
                  var resp                  *http.Response
                  var qryResponse map[string]ResponseFieldT
                  var nHost                  int
//                  var buf                  []byte
//                  var n                 int

                  Goose.Search.Logf(2,"TID:%s searching instance %d: ",tid,instance)

                  defer func() { rep<- qryResponse }()

                  for nHost=0; nHost<len(sb.Providers[instance].Bot.Host); nHost++ {
                     if sb.Providers[instance].Bot.Listen[0] == ':' {
                        host = sb.Providers[instance].Bot.Host[nHost].Name
                     } else {
                        nHost = len(sb.Providers[instance].Bot.Host)
                     }
                     host = sb.Providers[instance].Operation.Schemes[0] + "://" + host + sb.Providers[instance].Bot.Listen

                     path = sb.Providers[instance].Path
                     body = map[string]interface{}{}

                     Goose.Search.Logf(4,"TID:%s  Will add search path=%s, body=%#v",tid,path,body)
                     Goose.Search.Logf(4,"TID:%s  Swagger reports the operation parameters are %#v",tid,sb.Providers[instance].Operation.Parameters)
                     hasQueryParm = false
                     for _, swParm = range sb.Providers[instance].Operation.Parameters {
                        Goose.Search.Logf(3,"TID:%s  adding search parm: %s",tid,swParm.Name)
                        if swParm.In == "path" {
                           path = strings.Replace(path,"{" + swParm.Name + "}",searchBy[swParm.Name],-1)
                           Goose.Search.Logf(3,"TID:%s path now is: %s",tid,path)
                        } else if swParm.In == "query" {
                           if !hasQueryParm {
                              path += "?"
                              hasQueryParm = true
                           } else {
                              path += "&"
                           }
                           path += swParm.Name + "=" + url.QueryEscape(searchBy[swParm.Name])
                           Goose.Search.Logf(4,"TID:%s path now is: %#v",tid,path)
                        } else if swParm.In == "body" {
                           body[swParm.Name] = searchBy[swParm.Name]
                           Goose.Search.Logf(4,"TID:%s body now is: %#v",tid,body)
                        }
                     }

                     if sb.Providers[instance].Operation.Consumes[0] == "application/json" {
                        b_body, err = json.Marshal(body)
                     } else if sb.Providers[instance].Operation.Consumes[0] == "application/xml" {
                        b_body, err = xml.Marshal(body)
                     }

                     if err != nil {
                        Goose.Search.Logf(1,"TID:%s %s: %s",tid,ErrMarshalingRequestBody,err)
                        return
                     }

                     Goose.Search.Logf(4,"TID:%s Requesting search via %s:%s%s with body=%#v",tid,sb.Providers[instance].HttpMethod, host, path, body)

                     req, err = http.NewRequest(sb.Providers[instance].HttpMethod, host + path, bytes.NewReader(b_body))
                     if err!=nil {
                        Goose.Search.Logf(1,"%s: %s",ErrAssemblyingRequest,err)
                        return
                     }

                     for _, swParm = range sb.Providers[instance].Operation.Parameters {
                        if swParm.In == "header" {
                           if _, ok := req.Header[swParm.Name]; ok {
                              req.Header[swParm.Name] = append(req.Header[swParm.Name],searchBy[swParm.Name])
                           } else {
                              req.Header[swParm.Name] = []string{searchBy[swParm.Name]}
                           }

//                           req.Header.Add(swParm.Name,searchBy[swParm.Name])
                        }
                     }

                     resp, err = sb.HttpsSearchClient.Do(req)
                     if err!=nil {
                        Goose.Search.Logf(1,"TID:%s %s: %s",tid,ErrQueryingSearchBot,err)
                        continue // Let's try querying another instance of the search bots
                     }

/*
                     if resp.ContentLength > 0 {
                        b_body = make([]byte,resp.ContentLength)
                        err = io.ReadFull(resp.Body,b_body)
                     } else {
                        b_body = make([]byte,bufsz)
                        buf    = b_body
                        err    = nil
                        for err == nil {
                           n, err = resp.Body.Read(buf)
                           if (n==0) && (err==io.EOF) {
                              break
                           }
                           if (err!=nil) && (err!=io.EOF) {
                              Goose.Search.Logf(1,"TID:%s %s: %s",tid,ErrReadingResponseBody,err)
                              return
                           }
                           if n < bufsz {
                              b_body = b_body[:len(b_body)-(bufsz-n)]
                              break
                           }
                           b_body = append(b_body,...make([]byte,bufsz))
                           buf    = b_body[len(b_body)-bufsz:]
                        }
                     }
*/
//                     tmpbody, _ :=ioutil.ReadAll(resp.Body)
//                     Goose.Search.Logf(4,"TID:%s Response: %s",tid,tmpbody)

                     b_body, err = ioutil.ReadAll(resp.Body)
                     if err != nil {
                        Goose.Search.Logf(1,"TID:%s %s: %s",tid,err,b_body)
                        return
                     }

                     Goose.Search.Logf(1,"TID:%s body: %s",tid,b_body)

                     if sb.Providers[instance].Operation.Produces[0] == "application/json" {
                        err = json.NewDecoder(bytes.NewReader(b_body)).Decode(&body)
                     } else if sb.Providers[instance].Operation.Produces[0] == "application/xml" {
                        err = xml.NewDecoder(bytes.NewReader(b_body)).Decode(&body)
                     }

                     Goose.Search.Logf(4,"TID:%s Response body: %#v",tid,body)

                     if err != nil {
                        Goose.Search.Logf(1,"TID:%s %s: %s",tid,ErrUnmarshalingRequestBody,err)
                        return
                     }

                     qryResponse = map[string]ResponseFieldT{}
                     for _, fieldName := range searchFor {
                        Goose.Search.Logf(4,"TID:%s fetching Field %s: %#v",tid,fieldName,body[fieldName])
                        Goose.Search.Logf(7,"TID:%s body[dtupdate]: %#v",tid,body["dtupdate"])
                        Goose.Search.Logf(7,"TID:%s sb.Providers[instance].BotId: %#v",tid,sb.Providers[instance].BotId)
                        Goose.Search.Logf(7,"TID:%s sb.Providers[instance].Bot.Host[nHost]: %#v",tid,sb.Providers[instance].Bot.Host[nHost])
                        qryResponse[fieldName] = ResponseFieldT{
                           Value:    body[fieldName],
                           Source:   sb.Providers[instance].BotId + "@" + sb.Providers[instance].Bot.Host[nHost].Name,
                           DtUpd:    body["DtUpdate"].(string),
                        }
                     }

                     Goose.Search.Logf(4,"TID:%s ResponseFieldT: %#v",tid,qryResponse)
                     break
                  }
               }(i,rep)
            }
         }

         responses = map[string][]ResponseFieldT{}
         for respCount < oneShotProviders.Len() {
            Goose.Search.Logf(4,"TID:%s Waiting response %d/%d",tid,respCount, oneShotProviders.Len())
            response = <-rep
            respCount++
            if response != nil {
               Goose.Search.Logf(4,"TID:%s Got response from bot: %#v",tid,response)
               for k,v := range response {
                  if _, ok := responses[k]; ok {
                     responses[k] = append(responses[k],v)
                  } else {
                     responses[k] = []ResponseFieldT{v}
                  }
               }
            } else {
               Goose.Search.Logf(1,"TID:%s Bot instance failed",tid)
            }
         }

         Goose.Search.Logf(4,"TID:%s Final consolidated ResponseFieldT: %#v",tid,responses)
         return stonelizard.Response{
            Status: http.StatusOK,
            Body: responses,
         }
      }
   }

/*
   if !isFragmented {
      for i=0; i < providers.Max(); i++ {
         if providers.Has(i) {
            go func(instance int) {
               sb.HttpsSearchClient.Get(sb.Providers[i])
            }(i)
         }
      }
   }
*/

   return stonelizard.Response{
      Status: http.StatusOK,
      Body: "Unimplemented yet!",
   }
}


