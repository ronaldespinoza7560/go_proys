package main

import (
	"fmt"

	bd "github.com/ronaldespinoza7560/go_proys/server_api/basedatos"
)

func main() {

	tot := bd.ValidarUsuario("b!", "c", "users")
	fmt.Print(tot)

}
