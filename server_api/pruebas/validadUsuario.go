package main

import (
	"fmt"

	bd "github.com/ronaldespinoza7560/go_proys/server_api/basedatos"
)

func main() {

	priv := bd.ValidarUsuario("b!", "c", "users")
	fmt.Print(priv.Ingreso,priv.Usuario,priv.Nivel_acceso,priv.Accesos)

}
