package main

import (
	"fmt"
	"github.com/corona10/goimagehash"
	"image/jpeg"
  "os"
  "io/ioutil"
	"log"
	"strings"
  "bufio"
  "time"
	
)

func calculahash(f1 string, fdupli chan map[uint64][]string) {
  if strings.Index(f1, "_")>=0{
    // fmt.Println(f1)
    file1, _ := os.Open(f1)
    defer file1.Close()
    img1, _ := jpeg.Decode(file1)
    hash1, _ := goimagehash.AverageHash(img1)
    if hash1 !=nil {
     //fmt.Println(hash1.GetHash())
     fduplis[hash1.GetHash()]=append(fduplis[hash1.GetHash()],f1)
    }
   }
   count = count-1
return
}

func main() {
  fdupli := make(chan map[uint64][]string)
  count := make(chan int64)
  
  count=0
  reader := bufio.NewReader(os.Stdin)
  df:=""
  fmt.Print("Ingrese el directorio donde estan las fotos: ")
  fmt.Scanf("%s", &df)
  fmt.Println(df)
  text, _ := reader.ReadString('\n')

  fmt.Println(text)

  df1:=strings.Replace(df, "\\", "/", -1)
  files, err := ioutil.ReadDir(df1+"/.")
  if err != nil {
    log.Fatal(err)
  }

  start := time.Now()

  fduplis=make(map[uint64][]string)
  f1 := ""
  
  for _, file := range files {
    count = count+1
    f1 = df+"\\"+file.Name()
    go calculahash(f1)
    if count > 1 {
      time.Sleep(100 * time.Millisecond)
    }
    
  }
  


  for count > 0 {
		time.Sleep(100 * time.Millisecond)
	}
  
  // fmt.Println(fduplis)
  cc:=0
  for kk, ff := range fduplis {
    if len(ff)>1{
      cc=cc+1
      fmt.Println(cc,kk,ff)
    }
     
  }


  t := time.Now()
  elapsed := t.Sub(start)
  fmt.Println(elapsed)
}