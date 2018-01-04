
package manejadores
/**
* Aqui se definen todos los manejadores de las rutas.
*/
import (
	"encoding/json"
	"fmt"
	"net/http"
 	bd "github.com/ronaldespinoza7560/go_proys/server_api/basedatos"
)

type UserConsulta struct {
	NombreConsulta string `json:"nombre_consulta"`  //update,delete,insert,query
	// Parametro1 string `json:"parametro1"`		   //
	// Parametro2 string `json:"parametro2"`
	// Parametro3 string `json:"parametro3"`
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {

	// response := Response{"Gained access to protected resource"}
	// fmt.Println(claim_user)
	// fmt.Println(claim_nivel_acceso)
	// fmt.Println(claim_accesos)
	// JsonResponse(response, w)
	fmt.Println(UserPrivilegios)
	fmt.Println(UserPrivilegios.user)
	consulta := "SELECT * FROM " + "tp_suspendidos limit 2"
	resultado, err := bd.GetJSON(consulta)
	if err!=nil{
		entry := make(map[string]interface{})
		entry["ERROR"] = "hubo un error en la consulta"
		JsonResponse(entry, w)
	}
	JsonResponse(resultado, w)

}

func Bts_alarmasHandler(w http.ResponseWriter, r *http.Request) {
	consulta := "SELECT * FROM " + "tp_suspendidos limit 20"
	resultado, err := bd.GetJSON(consulta)
	if err!=nil{
		entry := make(map[string]interface{})
		entry["ERROR"] = "hubo un error en la consulta"
		JsonResponse(entry, w)
	}
	fmt.Println(resultado)
	JsonResponse(resultado, w)
}

func Bts_alarmasHandler1(w http.ResponseWriter, r *http.Request) {
	
	consulta := CrearQuery(r)
	resultado, err := bd.GetJSON(consulta)
	if err!=nil{
		entry := make(map[string]interface{})
		entry["ERROR"] = "hubo un error en la consulta"
		JsonResponse(entry, w)
	}
	JsonResponse(resultado, w)
}

func CrearQuery(r *http.Request)(string){
	
	var userQuery UserConsulta
	var query_string=""
	err := json.NewDecoder(r.Body).Decode(&userQuery)
	if err!=nil{
		return ""
	}

	if userQuery.NombreConsulta == "extraer_alarmas"{
		query_string=`select * from bts_alarmas where
		cell_name in (select SiteName from bts_latlon where clave like "%dis:rimac%") limit 20`
	}

	return query_string

}

func JsonResponse(response interface{}, w http.ResponseWriter) {

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
