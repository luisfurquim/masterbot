package searchbot

import (
//   "fmt"
//   "time"
   "bytes"
   "strings"
   "net/http"
   "encoding/json"
   "encoding/xml"
//   "github.com/luisfurquim/masterbot"
   "github.com/luisfurquim/stonelizard"
   "golang.org/x/tools/container/intsets"
)

func (sb SearchBotT) Search(searchBy map[string]string, searchFor []string) {
   var providers         *intsets.Sparse
   var searchFields      *intsets.Sparse
   var commonFields      *intsets.Sparse
   var oneShotProviders  *intsets.Sparse
   var field              string
   var isFragmented       bool
   var i                  int
   var p                 *TaxonomyTreeT
//   var prov               int

   providers = &intsets.Sparse{}

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
         return
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
         for i=0; i <= oneShotProviders.Max(); i++ {
            if oneShotProviders.Has(i) {
               go func(instance int) {
                  var err               error
                  var req              *http.Request
                  var host              string
                  var path              string
                  var swParm            stonelizard.SwaggerParameterT
                  var body   map[string]interface{}
                  var b_body          []byte
                  var resp             *http.Response
//                  var buf             []byte
//                  var n                 int
                  var nHost             int

                  Goose.Logf(1," searching instance %d: ",instance)

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
                     for _, swParm = range sb.Providers[instance].Operation.Parameters {
                        Goose.Logf(3," adding search parm: %s",swParm.Name)
                        if swParm.In == "path" {
                           path = strings.Replace(path,"{" + swParm.Name + "}",searchBy[swParm.Name],-1)
                           Goose.Logf(3," path now is: %s",path)
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

                     if sb.Providers[instance].Operation.Produces[0] == "application/json" {
                        err = json.NewDecoder(resp.Body).Decode(&body)
                     } else if sb.Providers[instance].Operation.Produces[0] == "application/xml" {
                        err = xml.NewDecoder(resp.Body).Decode(&body)
                     }

                     if err != nil {
                        Goose.Logf(1,"%s: %s",ErrUnmarshalingRequestBody,err)
                        return
                     }

                     Goose.Logf(4,"Response body: %#v",body)
                     break
                  }
               }(i)
            }
         }
         return
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
}


