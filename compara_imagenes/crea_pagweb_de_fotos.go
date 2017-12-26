package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"bufio"
	"os"
	)

func main() {
	
	reader := bufio.NewReader(os.Stdin)
    df:=""
	fmt.Print("Ingrese el directorio donde estan las fotos: ")
	fmt.Scanf("%s", &df)
	fmt.Println(df)
	text, _ := reader.ReadString('\n')
	

    fmt.Println("Ingrese el titulo de la pagina web: ")
    titulo := ""
	fmt.Scanf("%s", &titulo)

	

	text, _ = reader.ReadString('\n')
	
	fmt.Println(text)
	
	df=strings.Replace(df, "\\", "/", -1)
	files, err := ioutil.ReadDir(df+"/.")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(titulo)

	tabla:=fmt.Sprintf(`<!DOCTYPE html>
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
	  </tr>`,titulo)

	tabla_fin:=`</table> 
	</div>
	</body>
	</html>`

	nro_id:=""
	nro_id1:=""
	flag:=0
	flag_unico:=0
	str_tmp:=""
	str_tmp1:=""

	s := strings.Split("12345678_xxxx_yyyyy_zzzz", "_")
	nro_id=s[0]+"_"+s[1]
	
	for _, file := range files {
		if strings.Index(file.Name(), "_")>=0{

			str_tmp1=`<tr style='background-color: skyblue;font-weight: bold;'>
			<td style='background-color: skyblue;'>`+file.Name()[:5]+`</td>
			<td style='background-color: skyblue;'></td>
			<td style='background-color: skyblue;'></td>
			<td style='background-color: skyblue;'></td>
			<td style='background-color: skyblue;'></td></tr>`
			
		s := strings.Split(nro_id, "_")

		s1 := strings.Split(file.Name(), "_")
		nro_id1=s1[0]+"_"+s1[1]
		if flag_unico==0{
			tabla=tabla+str_tmp1
			flag_unico=1
			
		}

		if nro_id1!=nro_id&&flag==0{
			//fmt.Println(nro_id1,nro_id)
			//fmt.Println(nro_id[6:],string(nro_id[6:])=="a")
			if s[1]=="d"&&s1[1]=="a"{
				tabla=tabla+str_tmp1
				
			}
			
			tabla=tabla+"<tr>"
			nro_id=nro_id1
			flag=1
		}
		
		if nro_id1!=nro_id&&flag==1{
			if s[1]=="d"&&s1[1]=="a"{
				tabla=tabla+str_tmp1
			}
			tabla=tabla+"</tr>"
			nro_id=nro_id1
			flag=0
		}
		str_tmp="<td style='font-size: 12px;color: blue;font-weight: bold;'>"+file.Name()+"<br> <img src=\""+df+"/"+file.Name()+"\" style=\"width:150px;height:250px;\"> </td>"
		
		tabla=tabla+str_tmp
		
	}
	}
	
	
	pag:=tabla+"</tr>" + tabla_fin
	fmt.Print(pag)
	fo, err := os.Create(titulo+".html")
    if err != nil {
        panic(err)
	}
	if _, err := fo.Write([]byte(pag)); err != nil {
		panic(err)
	}
	// open output file
    
	
	fo.Close()
}