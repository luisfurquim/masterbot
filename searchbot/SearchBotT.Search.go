package searchbot

import (
   "fmt"
//   "time"
   "bytes"
   "strings"
   "net/url"
   "net/http"
//   "io/ioutil"
   "encoding/json"
   "encoding/xml"
//   "github.com/luisfurquim/masterbot"
   "github.com/luisfurquim/stonelizard"
   "github.com/luisfurquim/goose"
   "golang.org/x/tools/container/intsets"
)

func (sb *SearchBotT) Search(searchBy map[string]string, searchFor []string, login string, password string) stonelizard.Response {
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

   Goose = goose.Alert(6)

   defer func() { Goose = goose.Alert(2) }()


   providers = &intsets.Sparse{}

   Goose.Logf(2,"len(sb.Providers): %d",len(sb.Providers))

   // Fill the providers set with all provider currently known
   for i=0; i < len(sb.Providers); i++ {
      providers.Insert(i)
   }

   // Determine if there is at least one bot providing all data needed
   // by repeatedly computing providers ∩= 'providers of a given field'
   for _, field = range searchFor {
      i, _, p = sb.Taxonomy.Search(field)
      if ((i+1)!=len(field)) || (p==nil) || (p.Id<0) {
         Goose.Logf(1,"%s: %s",ErrUndefinedField, field)
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
   }

   Goose.Logf(4,"Determined if there is at least one bot providing all data needed (isFragmented=%#v): %#v",isFragmented,providers)

   if !isFragmented {
      // Select in the bots that have all information needed
      // those who require only information we have
      searchFields = &intsets.Sparse{}
      for field, _ = range searchBy {
         i, _, p = sb.Taxonomy.Search(field)
         searchFields.Insert(p.Id)
      }

      Goose.Logf(4,"Bitstring of search created: %#v",searchFields)

      oneShotProviders = &intsets.Sparse{}
      commonFields     = &intsets.Sparse{}

      Goose.Logf(4,"providers.Max(): %d",providers.Max())
      for i=0; i <= providers.Max(); i++ {
         Goose.Logf(4,"Bitstring of sb.Providers[%d].Requires: %#v",i,sb.Providers[i].Requires)
         commonFields.Intersection(searchFields,sb.Providers[i].Requires)
         if commonFields.Len() == sb.Providers[i].Requires.Len() {
            Goose.Logf(4,"Intersection at %d",i)
            oneShotProviders.Insert(i)
         }
      }

      Goose.Logf(4,"Bitstring of oneShotProviders: %#v",oneShotProviders)

      // If there is at least one bot who gives all fields
      // we need and requires just fields we already have...
      if oneShotProviders.Len() > 0 {
         rep = make(chan map[string]ResponseFieldT,oneShotProviders.Len())
         Goose.Logf(4,"len(sb.Providers): %d",len(sb.Providers))
         for i=0; i <= oneShotProviders.Max(); i++ {
            Goose.Logf(4,"oneShotProvider: %d",i)
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
//                  var buf             []byte
//                  var n                 int

                  Goose.Logf(1," searching instance %d: ",instance)

                  defer func() { rep<- qryResponse }()

                  for nHost=0; nHost<len(sb.Providers[instance].Bot.Host); nHost++ {
                     if sb.Providers[instance].Bot.Listen[0] == ':' {
                        host = sb.Providers[instance].Bot.Host[nHost]
                     } else {
                        nHost = len(sb.Providers[instance].Bot.Host)
                     }
                     host = sb.Providers[instance].Operation.Schemes[0] + "://" + host + sb.Providers[instance].Bot.Listen

                     path = sb.Providers[instance].Path
                     body = map[string]interface{}{}

                     Goose.Logf(4," Will add search path=%s, body=%#v",path,body)
                     Goose.Logf(4," Swagger reports the operation parameters are %#v",sb.Providers[instance].Operation.Parameters)
                     hasQueryParm = false
                     for _, swParm = range sb.Providers[instance].Operation.Parameters {
                        Goose.Logf(3," adding search parm: %s",swParm.Name)
                        if swParm.In == "path" {
                           path = strings.Replace(path,"{" + swParm.Name + "}",searchBy[swParm.Name],-1)
                           Goose.Logf(3," path now is: %s",path)
                        } else if swParm.In == "query" {
                           if !hasQueryParm {
                              path += "?"
                              hasQueryParm = true
                           } else {
                              path += "&"
                           }
                           path += swParm.Name + "=" + url.QueryEscape(searchBy[swParm.Name])
                           Goose.Logf(4," path now is: %#v",path)
                        } else if swParm.In == "body" {
                           body[swParm.Name] = searchBy[swParm.Name]
                           Goose.Logf(4," body now is: %#v",body)
                        }
                     }

                     if sb.Providers[instance].Operation.Consumes[0] == "application/json" {
                        b_body, err = json.Marshal(body)
                     } else if sb.Providers[instance].Operation.Consumes[0] == "application/xml" {
                        b_body, err = xml.Marshal(body)
                     }

                     if err != nil {
                        Goose.Logf(1,"%s: %s",ErrMarshalingRequestBody,err)
                        return
                     }

                     Goose.Logf(4,"Requesting search via %s:%s%s with body=%#v",sb.Providers[instance].HttpMethod, host, path, body)

                     req, err = http.NewRequest(sb.Providers[instance].HttpMethod, host + path, bytes.NewReader(b_body))
                     if err!=nil {
                        Goose.Logf(1,"%s: %s",ErrAssemblyingRequest,err)
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
                        Goose.Logf(1,"%s: %s",ErrQueryingSearchBot,err)
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
                              Goose.Logf(1,"%s: %s",ErrReadingResponseBody,err)
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
//                     Goose.Logf(4,"Response: %s",tmpbody)

                     if sb.Providers[instance].Operation.Produces[0] == "application/json" {
                        err = json.NewDecoder(resp.Body).Decode(&body)
                     } else if sb.Providers[instance].Operation.Produces[0] == "application/xml" {
                        err = xml.NewDecoder(resp.Body).Decode(&body)
                     }

                     Goose.Logf(4,"Response body: %#v",body)

                     if err != nil {
                        Goose.Logf(1,"%s: %s",ErrUnmarshalingRequestBody,err)
                        return
                     }

                     qryResponse = map[string]ResponseFieldT{}
                     for _, fieldName := range searchFor {
                        Goose.Logf(4,"fetching Field %s: %#v",fieldName,body[fieldName])
                        Goose.Logf(7,"body[dtupdate]: %#v",body["dtupdate"])
                        Goose.Logf(7,"sb.Providers[instance].BotId: %#v",sb.Providers[instance].BotId)
                        Goose.Logf(7,"sb.Providers[instance].Bot.Host[nHost]: %#v",sb.Providers[instance].Bot.Host[nHost])
                        qryResponse[fieldName] = ResponseFieldT{
                           Value:    body[fieldName],
                           Source:   sb.Providers[instance].BotId + "@" + sb.Providers[instance].Bot.Host[nHost],
                           DtUpd:    body["DtUpdate"].(string),
                        }
                     }

                     Goose.Logf(4,"ResponseFieldT: %#v",qryResponse)
                     break
                  }
               }(i,rep)
            }
         }

         responses = map[string][]ResponseFieldT{}
         for respCount < oneShotProviders.Len() {
            Goose.Logf(4,"Waiting response %d/%d",respCount, oneShotProviders.Len())
            response = <-rep
            respCount++
            if response != nil {
               Goose.Logf(4,"Got response from bot: %#v",response)
               for k,v := range response {
                  if _, ok := responses[k]; ok {
                     responses[k] = append(responses[k],v)
                  } else {
                     responses[k] = []ResponseFieldT{v}
                  }
               }
            } else {
               Goose.Logf(1,"Bot instance failed")
            }
         }

         Goose.Logf(4,"Final consolidated ResponseFieldT: %#v",responses)
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


