package main

import (
	"fmt"
	controllers "glogin/controllers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3000"
)

// var DB *sqlx.DB
// var DBError error

// func init() {
// 	DB, DBError := sqlx.Connect("postgres", "user=postgres password=momin1234 dbname=glogin sslmode=disable")
// 	if DBError != nil {
// 		log.Fatalln("error occur when database conneting", DBError)
// 	}
// 	fmt.Println(DB)
// }

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/super/admin", controllers.SuperAdminCreate).Methods("POST")
	router.HandleFunc("/super/admin/login", controllers.SuperAdminLogin).Methods("POST")
	router.HandleFunc("/upload", controllers.FileSave).Methods("POST")
	router.HandleFunc("/show", controllers.FileShow).Methods("POST")

	if err := http.ListenAndServe(fmt.Sprintf(":%s", CONN_PORT), router); err != nil {
		log.Fatal("error starting server: ", err)
	}
}
