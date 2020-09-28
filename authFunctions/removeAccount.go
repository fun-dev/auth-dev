package authFunctions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func removeAccount(w http.ResponseWriter, r *http.Request) {
	var user User
	var error Error
	//var jwt JWT

	json.NewDecoder(r.Body).Decode(&user)

	if user.Email == "" {
		error.Message = "Email は必須です．"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	if user.Password == "" {
		error.Message = "パスワードは必須です．"
		errorInResponse(w, http.StatusBadRequest, error)
	}

	password := user.Password
	fmt.Println("password: ", password)

	// データベースに接続
	db, err := sql.Open("mysql", DBLocation())
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// 認証キー(Emal)のユーザー情報をDBから取得
	row := db.QueryRow("SELECT * FROM users WHERE email = ?;", user.Email)

	err = row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			error.Message = "ユーザが存在しません．"
			errorInResponse(w, http.StatusBadRequest, error)
		} else {
			log.Fatal(err)
		}
	}

	hasedPassword := user.Password
	fmt.Println("hasedPassword: ", hasedPassword)

	err = bcrypt.CompareHashAndPassword([]byte(hasedPassword), []byte(password))

	if err != nil {
		error.Message = "無効なパスワードです．"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}
}
