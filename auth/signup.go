package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	var user User
	var error Error

	fmt.Println(r.Body)

	json.NewDecoder(r.Body).Decode(&user)

	// ログイン情報が未入力の場合
	if user.Email == "" {
		error.Message = "Email は必須です．"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}
	if user.Password == "" {
		error.Message = "パスワードは必須です．"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	fmt.Println("---------------------")

	// パスワードのハッシュを生成
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("パスワード: ", user.Password)
	fmt.Println("ハッシュ化されたパスワード", hash)

	user.Password = string(hash)
	fmt.Println("コンバート後のパスワード: ", user.Password)

	// データベースに接続
	db, err := sql.Open("mysql", DBLocation())
	if err != nil {
		log.Println(err)
		return
	}

	err = db.Ping()
	if err != nil {
		log.Println(err)
		return
	}

	defer db.Close()

	// クエリの発行
	ins, err := db.Prepare("INSERT INTO users(email, password) VALUES(?, ?);")
	if err != nil {
		error.Message = "データベース処理に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	ins.Exec(user.Email, user.Password)

	defer ins.Close()

	// DB に登録できたらパスワードを空にする
	user.Password = ""
	w.Header().Set("Content-Type", "application/json")

	// JSON 形式で結果を返却
	responseByJSON(w, user)
}
