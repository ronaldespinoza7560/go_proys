package main

import (
	"bufio"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/corona10/goimagehash"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var fduplis map[uint64][]string
var mutex sync.Mutex

func calculahash(f1 string) {
	if strings.Index(f1, "_") >= 0 {
		// fmt.Println(f1)
		file1, _ := os.Open(f1)
		defer file1.Close()
		img1, _ := jpeg.Decode(file1)
		hash1, _ := goimagehash.AverageHash(img1)
		if hash1 != nil {
			//fmt.Println(hash1.GetHash())
			mutex.Lock()
			fduplis[hash1.GetHash()] = append(fduplis[hash1.GetHash()], f1)
			mutex.Unlock()
		}
	}

	wg.Done()
}

var wg sync.WaitGroup

func main() {

	reader := bufio.NewReader(os.Stdin)
	df := ""
	fmt.Print("Ingrese el directorio donde estan las fotos: ")
	fmt.Scanf("%s", &df)
	fmt.Println(df)
	text, _ := reader.ReadString('\n')

	fmt.Println(text)

	df1 := strings.Replace(df, "\\", "/", -1)
	files, err := ioutil.ReadDir(df1 + "/.")
	if err != nil {
		log.Fatal(err)
	}
	wg.Add(len(files))
	start := time.Now()

	fduplis = make(map[uint64][]string)
	f1 := ""

	for _, file := range files {

		f1 = df + "\\" + file.Name()
		go calculahash(f1)

	}
	wg.Wait()

	// fmt.Println(fduplis)
	cc := 0
	for kk, ff := range fduplis {
		if len(ff) > 1 {
			cc = cc + 1
			fmt.Println(cc, kk, ff)
		}

	}

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(elapsed)
}
