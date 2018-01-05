
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
	consultas:=[]string{"SELECT * FROM " + "tp_suspendidos limit 2"}
	//consulta := "SELECT * FROM " + "tp_suspendidos limit 2"
	resultado, err := bd.Get_datos_db(consultas)
	if err!=nil{
		entry := make(map[string]interface{})
		entry["ERROR"] = "hubo un error en la consulta"
		JsonResponse(entry, w)
	}
	JsonResponse(resultado, w)

}

func Bts_alarmasHandler(w http.ResponseWriter, r *http.Request) {
	consultas := []string{"SELECT * FROM " + "tp_suspendidos limit 5"}
	resultado, err := bd.Get_datos_db(consultas)
	if err!=nil{
		entry := make(map[string]interface{})
		entry["ERROR"] = "hubo un error en la consulta"
		JsonResponse(entry, w)
	}
	fmt.Println(resultado)
	JsonResponse(resultado, w)
}

func Bts_alarmasHandler_xx(w http.ResponseWriter, r *http.Request) {
	
	consultas := CrearQuery(r)
	resultado, err := bd.Get_datos_db(consultas)
	if err!=nil{
		entry := make(map[string]interface{})
		entry["ERROR"] = "hubo un error en la consulta"
		JsonResponse(entry, w)
	}
	JsonResponse(resultado, w)
}


func Bts_alarmasHandler_yy(w http.ResponseWriter, r *http.Request) {
	
	consulta := CrearQuery(r)
	resultado, err := bd.Inserta_actualiza_registros_db(consulta)
	if err!=nil{
		entry := make(map[string]interface{})
		entry["ERROR"] = "hubo un error en la consulta"
		JsonResponse(entry, w)
	}
	fmt.Println(UserPrivilegios)
	JsonResponse(resultado, w)
}


func CrearQuery(r *http.Request)([]string){
	
	
	var userQuery UserConsulta
	var query_string=""
	var consultas_db []string
	
	err := json.NewDecoder(r.Body).Decode(&userQuery)
	if err!=nil{
		return consultas_db
	}

	switch userQuery.NombreConsulta {
		//extrae alarmas
	case "extraer_alarmas":
		query_string=`select * from bts_alarmas where
		cell_name in (select SiteName from bts_latlon where clave like "%dis:rimac%") limit 20`
		consultas_db = append(consultas_db, query_string)

		//extrae bts
	case "extraer_bts":
		query_string=`select * from bts_latlon
		limit 10`
		consultas_db = append(consultas_db, query_string)

		//extrae network elements
	case "extraer_ne":
		query_string=`select * from bts_ne
		limit 10`
		consultas_db = append(consultas_db, query_string)

		//insertar un registro a la tabla usuario
	case "insertar_usuario":
		query_string=`INSERT INTO users (name,email,password ) VALUES ('aaaaaaa','bbbbbsadfb','ccccccccccc')`
		consultas_db = append(consultas_db, query_string)

	case "insertar_usuarios":
		query_string=`INSERT INTO users (name,email,password ) VALUES ('aaaaaaa','bbbbbsadfb','ccccccccccc')`
		consultas_db = append(consultas_db, query_string)
		query_string=`INSERT INTO users (name,email,password ) VALUES ('dddd','eeee','ffff')`
		consultas_db = append(consultas_db, query_string)

	default:
		query_string=""
		consultas_db = append(consultas_db, query_string)
	}
	
	return consultas_db

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
