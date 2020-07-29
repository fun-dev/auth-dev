package main

import (
	"./authFunctions"

	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func signup(w http.ResponseWriter, r *http.Request) {
	authFunctions.Signup(w, r)
}

func login(w http.ResponseWriter, r *http.Request) {
	authFunctions.Login(w, r)
}

var db *sql.DB

func main() {
	db, err := sql.Open("mysql", authFunctions.DBLocation())
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	db.Close()


	// urls.py
	router := mux.NewRouter()

	// endpoint
	router.HandleFunc("/signup", signup).Methods("POST")
	router.HandleFunc("/login", login).Methods("POST")

	// console に出力する
	log.Println("サーバー起動 : 8000 port で受信")

	// log.Fatal は、異常を検知すると処理の実行を止めてくれる
	log.Fatal(http.ListenAndServe(":8000", router))
}
