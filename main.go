package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexionBD() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contrasenia := ""
	Nombre := "sistema"

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

var plantillas = template.Must(template.ParseGlob("Template/*"))

func main() {
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/Create", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/Delete", Borrar)
	http.HandleFunc("/Update", Editar)
	http.HandleFunc("/actualizar", Actualizar)

	log.Print("Servidor On...")
	http.ListenAndServe(":8080", nil)
}
func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("txtid")
		nombre := r.FormValue("txtnombre")
		correo := r.FormValue("txtcorreo")

		conexionEstablecida := conexionBD()
		Modificar, err := conexionEstablecida.Prepare("update empleados set nombre=?,correo=? where id=?")
		if err != nil {
			panic(err.Error())
		}
		Modificar.Exec(nombre, correo, id)
		http.Redirect(w, r, "/", 301)
	}
}

func Editar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	conexionEstablecida := conexionBD()
	registro, err := conexionEstablecida.Query("select * from empleados where id=?", idEmpleado)
	if err != nil {
		panic(err.Error())
	}
	empleado := Empleado{}
	for registro.Next() {
		var id int
		var nombre, correo string
		err = registro.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

	}
	plantillas.ExecuteTemplate(w, "Update", empleado)

}

func Borrar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")

	conexionEstablecida := conexionBD()
	BorrarRegistro, err := conexionEstablecida.Prepare("delete from empleados where id=?")
	if err != nil {
		panic(err.Error())
	}
	BorrarRegistro.Exec(idEmpleado)
	http.Redirect(w, r, "/", 301)

}

type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

func Inicio(w http.ResponseWriter, r *http.Request) {
	conexionEstablecida := conexionBD()
	registros, err := conexionEstablecida.Query("select * from empleados")
	if err != nil {
		panic(err.Error())
	}

	empleado := Empleado{}
	arregloEmpleado := []Empleado{}
	for registros.Next() {
		var id int
		var nombre, correo string
		err = registros.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo
		arregloEmpleado = append(arregloEmpleado, empleado)
	}
	//fmt.Println(arregloEmpleado)

	plantillas.ExecuteTemplate(w, "Body", arregloEmpleado)
}
func Crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "Create", nil)
}
func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("txtnombre")
		correo := r.FormValue("txtcorreo")

		conexionEstablecida := conexionBD()
		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO empleados(nombre,correo)VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insertarRegistros.Exec(nombre, correo)
		http.Redirect(w, r, "/", 301)
	}
}
