package authFunctions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func createToken(user User) (string, error) {
	var err error

	// 鍵となる文字列(多分なんでもいい)
	secret := "auth-using-go"

	// Token を作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"iss":   "FunCloudAPI", // JWT の発行者が入る(文字列は任意)
	})

	//Dumpを吐く
	spew.Dump(token)

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("-----------------------------")
	fmt.Println("tokenString:", tokenString)

	return tokenString, nil
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user User
	var error Error
	var jwt JWT

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
	    log.Println(err)
	    return //関数を終了する
	}

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

	token, err := createToken(user)

	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	jwt.Token = token

	responseByJSON(w, jwt)
}
