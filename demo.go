package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
)
import "database/sql"
import "fmt"

func main() {


	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Get("/login", login)
	r.Post("/login",login)
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal(err)
	}

}
func login(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/login?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		log.Println(t.Execute(w, nil))
	} else {
		name := r.FormValue("username")
		pw := r.FormValue("password")
		rows, err := db.Query("select * from auth where name='" + name + "' and pw='" + pw + "'")
		checkErr(err)
		fmt.Println()
		//v := reflect.ValueOf(rows)
		//fmt.Println(v) //number := printResult(query)
		//dates := printResult(rows)
		if rows.Next()==false{
			w.Write([]byte("<b>用户名或密码错误</b>"))
		}else {
			w.Write([]byte("<b>登陆成功</b>"))
		}
		defer r.Body.Close()

	}
}

func checkErr(errMasg error) {
	if errMasg != nil {
		panic(errMasg)
	}
}
