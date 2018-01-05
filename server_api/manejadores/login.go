package manejadores
/**
* Aqui se definen todos los manejadores del login.
*/
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
//	"strings"
	"time"

	"crypto/rsa"
	
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	
	bd "github.com/ronaldespinoza7560/go_proys/server_api/basedatos"
	
)

var tabla_usuarios string = "users"


type UserPriv struct {
	user string
	nivel_acceso string
	accesos string
}
var UserPrivilegios UserPriv

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
	Permiso string `json:"permiso"`
}

const (
	// For simplicity these files are in the same folder as the app binary.
	// You shouldn't do this in production.
	privKeyPath = "app.rsa"
	pubKeyPath  = "app.rsa.pub"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func InitKeys() {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	fmt.Print("1")
	fatal(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fmt.Print("2")
	fatal(err)

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	fmt.Print("3")
	fatal(err)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fmt.Print("4")
	fatal(err)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var user UserCredentials

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}

	priv := bd.ValidarUsuario(user.Username, user.Password, tabla_usuarios)
	fmt.Print(priv)
	if !priv.Ingreso{
	
		w.WriteHeader(http.StatusForbidden)
		fmt.Println("Error logging in")
		fmt.Fprint(w, "Invalid credentials")
		return
	
	}

	token := jwt.New(jwt.SigningMethodRS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	//claims["exp"] = time.Now().Add(time.Second * time.Duration(30)).Unix()
	claims["iat"] = time.Now().Unix()

	claims["usuario"]=user.Username
	claims["nivel_acceso"]=priv.Nivel_acceso
	claims["accesos"]=priv.Accesos
	UserPrivilegios=UserPriv{user.Username,priv.Nivel_acceso,priv.Accesos}
	token.Claims = claims

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error extracting the key")
		fatal(err)
	}

	tokenString, err := token.SignedString(signKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		fatal(err)
	}

	response := Token{tokenString,"si"}
	JsonResponse(response, w)

}

func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	//clax := make(jwt.MapClaims)
	Token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		})
	
	if err == nil {
		if Token.Valid {
			Cx := Token.Claims.(jwt.MapClaims)
			UserPrivilegios=UserPriv{fmt.Sprintf("%s", Cx["usuario"]),fmt.Sprintf("%s", Cx["nivel_acceso"]),fmt.Sprintf("%s", Cx["accesos"])}
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized access to this resource")
	}
	
}

