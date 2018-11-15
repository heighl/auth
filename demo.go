package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"reflect"
)
import "database/sql"
import "fmt"

func main() {
	//db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/login?charset=utf8")
	//if err != nil {
	//	panic(err.Error())
	//}
	//defer db.Close()
	//
	//err = db.Ping()
	//if err != nil {
	//	panic(err.Error())
	//}
	//
	//rows, err := db.Query("select name,pw from auth")
	//defer rows.Close()
	//for rows.Next() {
	//	//var id int
	//	var name string
	//	var pw string
	//	err = rows.Scan(&name, &pw)
	//	fmt.Printf(name, pw)
	//}
	//err = rows.Err()
	//if err != nil {
	//	panic(err.Error())
	//}

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
		v := reflect.ValueOf(rows)
		fmt.Println(v) //number := printResult(query) 
		dates := printResult(rows)
		w.Write([]byte(dates))
		defer r.Body.Close()

	}
}

func checkErr(errMasg error) {
	if errMasg != nil {
		panic(errMasg)
	}
}
func printResult(query *sql.Rows) string {
	column, _ := query.Columns()
	values := make([][]byte, len(column))
	scans := make([]interface{}, len(column))
	for i := range values { //让每一行数据都填充到[][]byte里面
		scans[i] = &values[i]
	}
	results := make(map[int]map[string]string) //最后得到的map
	i := 0
	for query.Next() { //循环，让游标往下移动
		if err := query.Scan(scans...); err != nil {
			fmt.Println(err)
			return "3"
		}
		row := make(map[string]string)
		for k, v := range values {
			key := column[k]
			row[key] = string(v)
		}
		results[i] = row //装入结果集中
		i++
	}
	for k, v := range results {
		fmt.Println(k)
		fmt.Println(v)
		record := "<b>登陆成功</b>\nid:" + results[k]["id"] + "\n" + "用户名:" + results[k]["name"] + "\n" + "密码:" + results[k]["pw"]
		return record
	}
	return "<b>用户名或密码错误</b>"

}
