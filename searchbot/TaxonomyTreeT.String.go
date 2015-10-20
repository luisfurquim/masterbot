package searchbot

import (
   "fmt"
)


func (t TaxonomyTreeT) String() string {
   var n *TaxonomyTreeT
   var s string

   if t.Id >= 0 {
      s = fmt.Sprintf("%c[%d]",t.Rune,t.Id)
   } else {
      s = fmt.Sprintf("%c",t.Rune)
   }

   if len(t.Next) > 0 {
      s += "{"
      for _, n = range t.Next {
         s += n.String()
      }
      s += "}"
   }

   return s
}

