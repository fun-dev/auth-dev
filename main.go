package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/fun-dev/auth-dev/auth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func signup(w http.ResponseWriter, r *http.Request) {
	auth.Signup(w, r)
}

func login(w http.ResponseWriter, r *http.Request) {
	auth.Login(w, r)
}

func removeAccount(w http.ResponseWriter, r *http.Request) {
	auth.RemoveAccount(w, r)
}

var db *sql.DB

func main() {
	db, err := sql.Open("mysql", auth.DBLocation())
	if err != nil {
		log.Println(err)
		return
	}

	err = db.Ping()
	if err != nil {
		log.Println(err)
		return
	}

	db.Close()

	// urls.py
	router := mux.NewRouter()

	// endpoint
	router.HandleFunc("/signup", signup).Methods("POST")
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/removeAccount", removeAccount).Methods("POST")

	// console に出力する
	log.Println("サーバー起動 : 8000 port で受信")

	// log.Fatal は、異常を検知すると処理の実行を止めてくれる
	log.Fatal(http.ListenAndServe(":8000", router))
}
