package searchbot

import (
   "fmt"
)

func (t *TaxonomyTreeT) Add(item string) int {
   var i      int
   var r      rune
   var p, p2 *TaxonomyTreeT
   var sp2    string

   Goose.Logf(3,"Will search taxonomy: %s", item)

   i, p, p2 = t.Search(item)

   if p2 != nil {
      sp2 = fmt.Sprintf("%#v",*p2)
   } else {
      sp2 = "nil"
   }
   Goose.Logf(3,"tree search response: i=%d, p=%#v, p2=%s",i,*p,sp2)

   if (i+1) == len(item) {
      if p2 != nil {
         if p2.Id < 0 { // Item not found, add it
            Goose.Logf(3,"Item not found, add it")
            p2.Id = t.Id
            t.Id++
         }
         return p2.Id
      }
   }

   Goose.Logf(3,"Item not found, continuing: i=%d, item[i:]=%s",i,item[i:])

   for _, r = range item[i:] {
      p2 = &TaxonomyTreeT{
         Rune: r,
         Id: -1,
         Next: []*TaxonomyTreeT{},
      }
      p.Next = append(p.Next,p2)
      p = p2
   }

   Goose.Logf(3,"Item added: p=%#v",*p)
   if p2 != nil {
      sp2 = fmt.Sprintf("%#v",*p2)
   } else {
      sp2 = "nil"
   }
   Goose.Logf(3,"Item added: p2=%s",sp2)

   p2.Id = t.Id
   t.Id++
   return p2.Id
}
