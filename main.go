package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	//"./docs"

	//_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	//"github.com/swaggo/swag/example/celler/httputil"
	//"github.com/swaggo/swag/example/celler/model"
)

type List struct {
	Id        int
	Full_name string
	Birthday  string
	Address   string
}

var database *sql.DB

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := database.Exec("delete from Lists where id = $1", id)
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, "/", 301)
}

func EditPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	row := database.QueryRow("select * from Lists where id = $1", id)
	prod := List{}

	err := row.Scan(&prod.Id, &prod.Full_name, &prod.Birthday, &prod.Address)

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(404), http.StatusNotFound)
	} else {
		tmpl, _ := template.ParseFiles("templates/edit.html")
		tmpl.Execute(w, prod)
	}
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	id := r.FormValue("id")
	full_name := r.FormValue("full_name")
	birthday := r.FormValue("birthday")
	address := r.FormValue("address")

	_, err = database.Exec("update Lists set full_name=$1, birthday=$2, address = $3 where id = $4",
		full_name, birthday, address, id)

	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/", 301)
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "templates/create.html")
	} else {

		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		full_name := r.FormValue("full_name")
		birthday := r.FormValue("birthday")
		address := r.FormValue("address")

		_, err = database.Exec("insert into Lists (full_name, birthday, address) values ($1, $2, $3)",
			full_name, birthday, address)

		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/", 301)
	}
}

// @Summary Retrieves user based on given ID
// @Produce json
// @Param id path integer true "User ID"
// @Success 200 {object} models.User
// @Router /users/{id} [get]
func IndexHandler(w http.ResponseWriter, r *http.Request) {

	rows, err := database.Query("select * from Lists")

	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	lists := []List{}

	for rows.Next() {

		p := List{
			Id:        0,
			Full_name: "",
			Birthday:  "",
			Address:   "",
		}

		err := rows.Scan(&p.Id, &p.Full_name, &p.Birthday, &p.Address)
		p.Birthday = p.Birthday[0:10]

		if err != nil {
			fmt.Println(err)
			continue
		}

		lists = append(lists, p)
	}

	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, lists)
}

// @title Blueprint Swagger API
// @version 1.0
// @description Swagger API for Golang Project Blueprint.

// @BasePath /api/v1

func main() {

	db, err := sql.Open("postgres", "user=user password=pass dbname=betadb sslmode=disable")

	if err != nil {
		log.Println(err)
	}
	database = db
	defer db.Close()
	//router := mux.NewRouter()
	//router.HandleFunc("/", IndexHandler)
	//router.HandleFunc("/create", CreateHandler)
	//router.HandleFunc("/edit/{id:[0-9]+}", EditPage).Methods("GET")
	//router.HandleFunc("/edit/{id:[0-9]+}", EditHandler).Methods("POST")
	//router.HandleFunc("/delete/{id:[0-9]+}", DeleteHandler)

	//http.Handle("/", router)

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//r.PUT("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.IndexHandler))
	r.Run(":8080")
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//fmt.Println("Server is listening...")
	//http.ListenAndServe(":8181", nil)
}
