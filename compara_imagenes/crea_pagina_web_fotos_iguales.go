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
		//hash1, _ := goimagehash.PerceptionHash(img1)
		if hash1 != nil {
			//fmt.Println(hash1.GetHash())
			mutex.Lock()
			fduplis[hash1.GetHash()] = append(fduplis[hash1.GetHash()], f1)
			mutex.Unlock()
		}
	}

	wg.Done()
}

func crea_pagina_web(fdup map[uint64][]string) {
	titulo := "FOTOS_DUPLICADAS"

	tabla := fmt.Sprintf(`<!DOCTYPE html>
	<html>
	<body>
	<div>
	<h1 style='color: blue;'>%s</h1>
	 <table style='border-spacing: 0px;'>
	  <tr>
		<th>Trabajo1</th>
		<th>Trabajo2</th>
		<th>Trabajo3</th>
		<th>Trabajo4</th>
	  </tr>`, titulo)

	tabla_fin := `</table> 
	</div>
	</body>
	</html>`

	str_tmp := ""
	str_tmp1 := ""

	for _, file := range fdup {
		str_tmp = ""
		str_tmp1 = ""
		
		nom:=strings.Split(file[0],"\\")
		nom1:=nom[len(nom)-1]
		nombre:=strings.Split(nom1,"_")[0]
		if len(file) > 1 {
			str_tmp1 = `<tr style='background-color: skyblue;font-weight: bold;'>
			<td style='background-color: skyblue;'>` + nombre + `</td>
			<td style='background-color: skyblue;'></td>
			<td style='background-color: skyblue;'></td>
			<td style='background-color: skyblue;'></td>
			<td style='background-color: skyblue;'></td></tr>`

			for _, x := range file {
				nom2:=strings.Split(x,"\\")
				nom3:=nom2[len(nom2)-1]
				str_tmp = str_tmp + "<td style='font-size: 12px;color: blue;font-weight: bold;'>" + nom3 + "<br> <img src=\"" + x + "\" style=\"width:150px;height:250px;\"> </td>"
			}

			tabla = tabla + str_tmp1 + "<tr>" + str_tmp + "</tr>"
		}

	}
	pag := tabla + tabla_fin
	fmt.Print(pag)
	fo, err := os.Create(titulo + ".html")
	if err != nil {
		panic(err)
	}
	if _, err := fo.Write([]byte(pag)); err != nil {
		panic(err)
	}
	// open output file

	fo.Close()
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
	crea_pagina_web(fduplis)

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(elapsed)
}
