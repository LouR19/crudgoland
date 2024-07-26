package main

/* listo calisto */
import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Usuarios struct {
	ID     int
	Nombre string
	Correo string
}

var vistas = template.Must(template.ParseGlob("vistas/*"))

func cbd() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contra := ""
	Nombre := "Registro"

	conexion, err := sql.Open(Driver, Usuario+":"+Contra+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion

}

func main() {
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/crear", Create)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/edit", Editar)
	http.HandleFunc("/ed", Edit)
	log.Println("Servidor conectado...")
	http.ListenAndServe(":8080", nil)

}
func Edit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		cn_establecida := cbd()

		ed, err := cn_establecida.Prepare("UPDATE usuarios SET Nombre=?, Correo=? WHERE ID=?")
		if err != nil {
			panic(err.Error())
		}
		ed.Exec(nombre, correo, id)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)

	}

}

func Editar(w http.ResponseWriter, r *http.Request) {
	idusers := r.URL.Query().Get("id")
	fmt.Println(idusers)

	cn_establecida := cbd()
	edit, err := cn_establecida.Query("SELECT * FROM usuarios WHERE ID = ?", idusers)
	if err != nil {
		panic(err.Error())
	}

	usuario := Usuarios{}

	for edit.Next() {
		var ID int
		var nombre, correo string
		err = edit.Scan(&ID, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}

		usuario.ID = ID
		usuario.Nombre = nombre
		usuario.Correo = correo
	}
	fmt.Println(usuario)
	vistas.ExecuteTemplate(w, "edit", usuario)
}

func Borrar(w http.ResponseWriter, r *http.Request) {
	idusers := r.URL.Query().Get("id")
	cn_establecida := cbd()
	borrarus, err := cn_establecida.Prepare("DELETE FROM usuarios WHERE ID = ?")
	if err != nil {
		panic(err.Error())
	}
	borrarus.Exec(idusers)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

func Inicio(w http.ResponseWriter, r *http.Request) {
	cn_establecida := cbd()
	registro, err := cn_establecida.Query("SELECT * FROM usuarios")
	if err != nil {
		panic(err.Error())
	}

	usuario := Usuarios{}
	us := []Usuarios{}

	for registro.Next() {
		var ID int
		var nombre, correo string
		err = registro.Scan(&ID, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}

		usuario.ID = ID
		usuario.Nombre = nombre
		usuario.Correo = correo

		us = append(us, usuario)
	}

	vistas.ExecuteTemplate(w, "inicio", us)
}
func Create(w http.ResponseWriter, r *http.Request) {
	vistas.ExecuteTemplate(w, "crear", nil)
}
func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")
		cn_establecida := cbd()
		createusers, err := cn_establecida.Prepare("INSERT INTO usuarios(Nombre, Correo) VALUES (?,?)")
		if err != nil {
			panic(err.Error())
		}
		createusers.Exec(nombre, correo)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)

	}

}
