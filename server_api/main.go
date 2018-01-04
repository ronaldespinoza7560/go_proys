package main

//  Generate RSA signing files via shell (adjust as needed):
//
//  $ openssl genrsa -out app.rsa 1024
//  $ openssl rsa -in app.rsa -pubout > app.rsa.pub
//
// Code borrowed and modified from the following sources:
// https://www.youtube.com/watch?v=dgJFeqeXVKw
// https://goo.gl/ofVjK4
// https://github.com/dgrijalva/jwt-go
//

import (
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
	"log"
	"net/http"
//	"strings"
//	"time"

//	"crypto/rsa"
	"github.com/codegangsta/negroni"
//	"github.com/dgrijalva/jwt-go"
//	"github.com/dgrijalva/jwt-go/request"
	"github.com/rs/cors"
//	bd "github.com/ronaldespinoza7560/go_proys/server_api/basedatos"
	mj "github.com/ronaldespinoza7560/go_proys/server_api/manejadores"
	
)



func StartServer() {
	host:="localhost"
	port:="8080"

	mux := http.NewServeMux()
	
	mux.HandleFunc("/login", mj.LoginHandler)
	
	mux.HandleFunc("/xx", mj.Bts_alarmasHandler1)
	// Protected Endpoints

	mux.Handle("/resource", negroni.New(
		negroni.HandlerFunc(mj.ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(mj.ProtectedHandler)),
	))

	mux.Handle("/bts_alarmas", negroni.New(
		negroni.HandlerFunc(mj.ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(mj.Bts_alarmasHandler)),
	))

	log.Println("Now listening..."+host+":"+port)


	// cors.Default() setup the middleware with default options being
    // all origins accepted with simple methods (GET, POST). See
    // documentation below for more options.
    handler := cors.Default().Handler(mux)
	http.ListenAndServe(host+":"+port, handler)
}

// var claim_user string
// var claim_nivel_acceso string
// var claim_accesos string
var tabla_usuarios string = "users"


func main() {

	mj.InitKeys()
	StartServer()
}

