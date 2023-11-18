package main


import (
   "fmt"
   "sort"
   "encoding/hex"
   "errors"
   "crypto/rand"
)


var stackOfBlobs []string
var blobsToExecute []string


func addToStack(rollupBlobs []string) error {
   for _, element := range rollupBlobs {
       if float64(len(element))/2 <= 128 {
           stackOfBlobs = append(stackOfBlobs, element)
       } else {
           return errors.New("Blob too big.")
       }
   }
   return nil
}


func createBlob() {
   for len(stackOfBlobs) > 0 {
       var blobSize float64
       var blobData string


       // Sort stackOfBlobs by length in descending order
       sort.Slice(stackOfBlobs, func(i, j int) bool {
           return len(stackOfBlobs[i]) > len(stackOfBlobs[j])
       })


       fmt.Printf("Sorted L2 blobs: %v\n", stackOfBlobs)
       var removedBlobs []string


       for _, element := range stackOfBlobs {
           if (blobSize + float64(len(element))/2) <= 128 {
               fmt.Printf("Element: %s\n", element)
               fmt.Printf("Element length: %v\n", float64(len(element))/2)
               blobSize += float64(len(element)) / 2
               removedBlobs = append(removedBlobs, element)
               fmt.Printf("Blob size: %v\n", blobSize)
               blobData += element
           }
       }


       // Remove processed blobs from stackOfBlobs
       stackOfBlobs = subtractSets(stackOfBlobs, removedBlobs)


       fmt.Printf("Finished blob size: %v\n", blobSize)
       finishBlob(blobData)
   }
}


func finishBlob(blob string) {
   blobsToExecute = append(blobsToExecute, blob)
   fmt.Printf("Blob data: %s\n", blob)
}


// subtractSets subtracts elements in toRemove from slice
func subtractSets(slice, toRemove []string) []string {
   var result []string


OuterLoop:
   for _, s := range slice {
       for _, r := range toRemove {
           if s == r {
               continue OuterLoop
           }
       }
       result = append(result, s)
   }


   return result
}


func main() {
   // Example Usage:
   rollupBlobs := make([]string, 0)
   for n := 10; n <= 20; n++ {
       randomBytes := make([]byte, n)
       rand.Read(randomBytes)
       rollupBlobs = append(rollupBlobs, hex.EncodeToString(randomBytes))
   }
   rollupBlobs = append(rollupBlobs, rollupBlobs...) // Repeat 5 times
   fmt.Println(rollupBlobs)


   if err := addToStack(rollupBlobs); err != nil {
       fmt.Println(err)
       return
   }


   createBlob()
   fmt.Printf("Blobs to execute: %v\n", blobsToExecute)
}