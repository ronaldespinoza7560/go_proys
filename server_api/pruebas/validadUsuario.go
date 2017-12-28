package main
import (
	"fmt"
	bd "github.com/ronaldespinoza7560/go_proys/server_api/basedatos"
)

func main(){
	fmt.Print(bd.ValidarUsuario("b!","c","users"))
}