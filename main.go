package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/davecgh/go-spew/spew"
    "github.com/dgrijalva/jwt-go"
    "golang.org/x/crypto/bcrypt"
    _ "github.com/go-sql-driver/mysql"
)


type User struct {
    // 大文字だと Public 扱い
    ID       int    `json:"id"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

type JWT struct {
    Token string `json:"token"`
}

type Error struct {
    Message string `json:"message"`
}

func errorInResponse(w http.ResponseWriter, status int, error Error) {
    w.WriteHeader(status) //HTTP status コードが入る
    json.NewEncoder(w).Encode(error)
    return
}

func responseByJSON(w http.ResponseWriter, data interface{}) {
    json.NewEncoder(w).Encode(data)
    return
}


func signup(w http.ResponseWriter, r *http.Request) {
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
        log.Fatal(err)
    }

    fmt.Println("パスワード: ", user.Password)
    fmt.Println("ハッシュ化されたパスワード", hash)

    user.Password = string(hash)
    fmt.Println("コンバート後のパスワード: ", user.Password)


    // データベースに接続
    db, err := sql.Open("mysql", "root:hoge@/authdb")
    if err != nil {
        log.Fatal(err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal(err)
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


func createToken(user User) (string, error) {
    var err error

    // 鍵となる文字列(多分なんでもいい)
    secret := "auth-using-go"

    // Token を作成
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "email": user.Email,
        "iss":   "__init__", // JWT の発行者が入る(文字列(__init__)は任意)
    })

   //Dumpを吐く
    spew.Dump(token)

    tokenString, err := token.SignedString([]byte(secret))

    fmt.Println("-----------------------------")
    fmt.Println("tokenString:", tokenString)

    if err != nil {
        log.Fatal(err)
    }

    return tokenString, nil
}


func login(w http.ResponseWriter, r *http.Request) {
    var user User
    var error Error
    var jwt JWT

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
     db, err := sql.Open("mysql", "root:hoge@/authdb")
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
            error.Message = "ユーザが存在しません。"
            errorInResponse(w, http.StatusBadRequest, error)
        } else {
            log.Fatal(err)
        }
    }

    hasedPassword := user.Password
    fmt.Println("hasedPassword: ", hasedPassword)

    err = bcrypt.CompareHashAndPassword([]byte(hasedPassword), []byte(password))

    if err != nil {
        error.Message = "無効なパスワードです。"
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



var db *sql.DB

func main() {
    db, err := sql.Open("mysql", "root:hoge@/authdb")
    if err != nil {
        log.Fatal(err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }

    defer db.Close()

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
