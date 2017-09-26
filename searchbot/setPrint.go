package searchbot

import (
   "fmt"
   "strings"
   "golang.org/x/tools/container/intsets"
)

func setPrint(set *intsets.Sparse) string {
   var i int
   var s []string

   for i=0; i < set.Len(); i++ {
      if set.Has(i) {
         s = append(s,fmt.Sprintf("%d",i))
         if len(s) == set.Len() {
            break
         }
      }
   }

   return fmt.Sprintf("[%s]",strings.Join(s,", "))
}

