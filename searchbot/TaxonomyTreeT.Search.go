package searchbot

import (
   "fmt"
)

func (t *TaxonomyTreeT) Search(item string) (int, *TaxonomyTreeT, *TaxonomyTreeT) {
   var i      int
   var r      rune
   var p, p2 *TaxonomyTreeT
   var sp2    string

   p = t
   Goose.Taxonomy.Logf(7,"taxonomy.Search: p=%#v",p)

itemRunes:
   for i, r = range item {
      Goose.Taxonomy.Logf(7,"taxonomy.Search: i=%d, r=%c",i,r)
      if p.Next != nil {
         for _, p2 = range p.Next {
            if p2 != nil {
               sp2 = fmt.Sprintf("%#v",*p2)
            } else {
               sp2 = "nil"
            }
            Goose.Taxonomy.Logf(7,"taxonomy.Search.2: p2=%s",sp2)
            if p2.Rune == r {
               p = p2
               continue itemRunes
            }
         }
      }
      if p2 != nil {
         sp2 = fmt.Sprintf("%#v",*p2)
      } else {
         sp2 = "nil"
      }
      Goose.Taxonomy.Logf(7,"taxonomy.Search.3: p2=%s",sp2)
      p2 = nil
      break
   }

   Goose.Taxonomy.Logf(7,"taxonomy.Search end: i=%d, p=%#v",i,*p)

   return i, p, p2
}

