package main
import(
	"fmt"
	bd "github.com/ronaldespinoza7560/go_proys/server_api/basedatos"
)
func main(){
	
	querySelecs := []string{"select ne_name from network_elements limit 11"}
	tab, err := bd.Get_datos_db(querySelecs)
	if err !=nil{
		fmt.Println("err")
	}
	tab_n:=agrupa_nelems(tab, 1)
	fmt.Println(tab)
	fmt.Println(tab_n)
}
func agrupa_nelems(tab []map[string]interface{}, k int)([][]string){
	var rsets =[][]string{}
	var rset= []string{}
	
	cad:=""
	i:=1
	for _,v := range tab{
		cad=v["ne_name"].(string)
		rset=append(rset,cad)
		if !(i<k){
			rsets=append(rsets,rset)
			rset=nil
			i=0
		}
		i++
	}
	if len(rset)>0{
		rsets=append(rsets,rset)
	}
	return rsets
}