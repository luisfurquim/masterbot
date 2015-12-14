package searchbot

import (
   "fmt"
   "net/http"
   "github.com/luisfurquim/stonelizard"
)

func deepTaxonomy(p *TaxonomyTreeT, field string) []string {
   var ret   []string
   var p2     *TaxonomyTreeT

   Goose.Logf(3,"taxonomy.Dump field=%s, p=%#v",field,*p)

   if p.Rune != 0 {
      field += fmt.Sprintf("%c",p.Rune)
   }

   if p.Id >= 0 {
      ret = []string{field}
   } else {
      ret = []string{}
   }

   for _, p2 = range p.Next {
      ret = append(ret,deepTaxonomy(p2,field)...)
   }

   return ret
}

func (sb *SearchBotT) TaxonomyDump() stonelizard.Response {
   var tx []string
   var p   *TaxonomyTreeT

   tx = []string{}

   Goose.Logf(4,"taxonomy p=%#v",sb.Taxonomy)

   for _, p = range sb.Taxonomy.Next {
      Goose.Logf(3,"taxonomy.root fields=%#v, p=%#v",tx,*p)
      tx = append(tx,deepTaxonomy(p,"")...)
   }

   return stonelizard.Response{
      Status: http.StatusOK,
      Body: tx,
   }
}

