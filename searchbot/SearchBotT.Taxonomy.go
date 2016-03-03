package searchbot

import (
   "fmt"
   "net/http"
   "github.com/luisfurquim/stonelizard"
)

func deepTaxonomy(p *TaxonomyTreeT, field string, tid string) []string {
   var ret   []string
   var p2     *TaxonomyTreeT

   Goose.Search.Logf(6,"TID:%s taxonomy.Dump field=%s, p=%#v",tid,field,*p)

   if p.Rune != 0 {
      field += fmt.Sprintf("%c",p.Rune)
   }

   if p.Id >= 0 {
      ret = []string{field}
   } else {
      ret = []string{}
   }

   for _, p2 = range p.Next {
      ret = append(ret,deepTaxonomy(p2,field,tid)...)
   }

   return ret
}

func (sb *SearchBotT) TaxonomyDump(tid string) stonelizard.Response {
   var tx []string
   var p   *TaxonomyTreeT

   tx = []string{}

   Goose.Search.Logf(4,"TID:%s taxonomy dump requested",tid)
   Goose.Search.Logf(6,"TID:%s taxonomy p=%#v",tid,sb.Taxonomy)

   for _, p = range sb.Taxonomy.Next {
      Goose.Search.Logf(6,"TID:%s taxonomy.root fields=%#v, p=%#v",tid,tx,*p)
      tx = append(tx,deepTaxonomy(p,"",tid)...)
   }

   return stonelizard.Response{
      Status: http.StatusOK,
      Body: tx,
   }
}

