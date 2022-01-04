package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var DB *sqlx.DB
var DBErr error

func init() {
	DB, DBErr = sqlx.Connect("postgres", "user=postgres password=momin1234 dbname=glogin sslmode=disable")
	if DBErr != nil {
		log.Fatalln("error while connecting to database", DBErr.Error())
	}
}

type LoginForm struct {
	Email    string `db:"email"`
	Password string `db:"password"`
}
type Admin struct {
	ID       int32  `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

func (l Admin) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Email,
			validation.Required.Error("Email is required"),
		),
		validation.Field(&l.Password,
			validation.Required.Error("Password is required"),
			validation.Length(6, 16).Error("Password must be 6 to 16 characters length"),
		),
	)
}

func (l LoginForm) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Email,
			validation.Required.Error("Email is required"),
		),
		validation.Field(&l.Password,
			validation.Required.Error("Password is required"),
			validation.Length(6, 16).Error("Password must be 6 to 16 characters length"),
		),
	)
}

var mySigninKey = []byte("secretkey")

func GenerateJWT(userID int32) (string, error) {
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["id"] = userID
	claims["exp"] = time.Now().Add(time.Minute * 100).Unix()

	tokenString, err := token.SignedString(mySigninKey)
	if err != nil {
		fmt.Printf("something went wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil

}

func SuperAdminCreate(w http.ResponseWriter, req *http.Request) {
	var admin Admin

	err := json.NewDecoder(req.Body).Decode(&admin)
	if err != nil {
		log.Println("error occur while login req body data  : ", err.Error())
	}

	if vErr := admin.Validate(); vErr != nil {
		errorjson, err := json.Marshal(vErr)
		if err != nil {
			log.Println("error occur when convert in json ", err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write((errorjson))
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("error while encrypted password: ", err.Error())
	}
	admin.Password = string(hash)

	query := `
		INSERT INTO admin(
			email,
			password
		)
		VALUES(
			:email,
			:password
		)
		RETURNING id
	`
	var id int32
	stmt, err := DB.PrepareNamed(query)
	if err != nil {
		log.Println("db error: failed prepare ", err.Error())
		return
	}

	if err := stmt.Get(&id, admin); err != nil {
		log.Println("db error: failed to insert data ", err.Error())
		return
	}
}

func SuperAdminLogin(w http.ResponseWriter, req *http.Request) {
	var loginForm LoginForm

	err := json.NewDecoder(req.Body).Decode(&loginForm)
	if err != nil {
		log.Println("error occur while login req body data  : ", err.Error())
	}
	if vErr := loginForm.Validate(); vErr != nil {
		errorjson, err := json.Marshal(vErr)
		if err != nil {
			log.Println("error occur when convert in json ", err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write((errorjson))
		return
	}

	var admin Admin

	query := `SELECT  id, email, password FROM admin where email = $1`
	if err := DB.Get(&admin, query, loginForm.Email); err != nil {
		if err != nil {
			log.Println("error while getting user from db: ", err.Error())
			return
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(loginForm.Password)); err != nil {
		log.Println("password do not match: ", err.Error())
		return
	}

	token, err := GenerateJWT(admin.ID)
	if err != nil {
		fmt.Printf("error token")
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}
