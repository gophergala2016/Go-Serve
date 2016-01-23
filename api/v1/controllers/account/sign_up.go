package account

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gophergala2016/Go-Serve/api/v1/controllers"
	"github.com/gophergala2016/Go-Serve/api/v1/models"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type registrationController struct{}

var Registration registrationController

func (r registrationController) Create(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	flag := 1
	var u models.User

	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &u)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", "password=password host=localhost dbname=go_service_development sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	users, err := db.Exec("CREATE TABLE users (id SERIAL, name varchar(100), mobile_number varchar(100), password varchar(100), image varchar(2048), age int, gender varchar(100), device_token varchar(320), created_at timestamptz, PRIMARY KEY(id), UNIQUE (mobile_number))")
	if err != nil || users == nil {
		log.Fatal(err)
	}
	devices, err := db.Exec("CREATE TABLE devices (id int, devise_token varchar(320), user_id int, CONSTRAINT devices_users_key FOREIGN KEY(user_id) REFERENCES users(id), PRIMARY KEY(devise_token))")
	if err != nil || devices == nil {
		log.Fatal(err)
	}
	sessions, err := db.Exec("CREATE TABLE sessions (id int, user_id int, CONSTRAINT sessions_users_key FOREIGN KEY(user_id) REFERENCES users(id), devise_token varchar(320), CONSTRAINT sessions_devices_key FOREIGN KEY(devise_token) REFERENCES devices(devise_token));")
	if err != nil || sessions == nil {
		log.Fatal(err)
	}

	mobile_res, err := db.Query("SELECT mobile_number FROM users ")
	if err != nil {
		log.Fatal(err)
	}

	fetch_id, err := db.Query("SELECT coalesce(max(id), 0) FROM users")
	if err != nil {
		log.Fatal(err)
	}

	if flag == 1 {
		if u.Name == "" || u.Password == "" || u.Password_confirmation == "" || u.Devise_token == "" {

			result, err := govalidator.ValidateStruct(u)
			if err != nil {
				println("error: " + err.Error())
			}
			fmt.Println(result)
			flag = 0
			b, err := json.Marshal(models.ErrorMessage{
				Success: "false",
				Error:   err.Error(),
			})
			if err != nil {
				log.Fatal(err)
			}
			rw.Header().Set("Content-Type", "application/json")
			rw.Write(b)
			goto create_user_end
		}
	}
	if flag == 1 {
		for mobile_res.Next() {
			var email string
			err = mobile_res.Scan(&mobile_number)
			if err != nil {
				log.Fatal(err)
			}

			if mobile_number == u.Mobile_number {
				b, err := json.Marshal(models.ErrorMessage{
					Success: "false",
					Error:   "Mobile number already exist",
				})
				if err != nil {
					log.Fatal(err)
				}
				rw.Header().Set("Content-Type", "application/json")
				rw.Write(b)
				fmt.Println("Mobile_number already exist")
				flag = 0
				goto create_user_end
			}
		}
		defer mobile_res.Close()
		if u.Password != u.Password_confirmation {
			b, err := json.Marshal(models.ErrorMessage{
				Success: "false",
				Error:   "Password and Password_confirmation do not match",
			})
			if err != nil {
				log.Fatal(err)
			}
			rw.Header().Set("Content-Type", "application/json")
			rw.Write(b)
			goto create_user_end
		}
		session_response, err := db.Query("SELECT devise_token,user_id from sessions")
		if err != nil {
			log.Fatal(err)
		}
		for session_response.Next() {
			var devise_token string
			var id int
			err := session_response.Scan(&devise_token, &id)
			if err != nil {
				log.Fatal(err)
			}
			if devise_token == u.Devise_token {
				b, err := json.Marshal(models.ErrorMessage{
					Success: "false",
					Error:   "Session already Exist",
				})

				if err != nil {
					log.Fatal(err)
				}
				rw.Header().Set("Content-Type", "application/json")
				rw.Write(b)
				fmt.Println("Session already Exist")
				goto create_user_end
			}
		}
		for fetch_id.Next() {
			var id int
			err = fetch_id.Scan(&id)

			if err != nil {
				log.Fatal(err)
			}
			id = id + 1

			var sStmt string = "insert into users (id, name, mobile_number, password, devise_token) values ($1,$2,$3,$4,$5)"
			db, err := sql.Open("postgres", "password=password host=localhost dbname=postgres sslmode=disable")
			if err != nil {
				log.Fatal(err)
			}
			stmt, err := db.Prepare(sStmt)
			if err != nil {
				log.Fatal(err)
			}

			key := []byte("traveling is fun")
			password := []byte(u.Password)
			confirm_password := []byte(u.Password_confirmation)
			encrypt_password := controllers.Encrypt(key, password)
			encrypt_password_confirmation := controllers.Encrypt(key, confirm_password)

			user_res, err := stmt.Exec(id, u.Nmae, u.Mobile_number, encrypt_password, encrypt_password_confirmation, u.Devise_token)
			if err != nil || user_res == nil {
				log.Fatal(err)
			}

			var devise string = "insert into devices(devise_token,user_id)values ($1,$2)"
			dev, err := db.Prepare(devise)
			if err != nil {
				log.Fatal(err)
			}
			dev_res, err := dev.Exec(u.Devise_token, id)
			if err != nil || dev_res == nil {
				log.Fatal(err)
			}
			var session string = "insert into sessions (start_time, user_id,devise_token) values ($1,$2,$3)"
			ses, err := db.Prepare(session)
			if err != nil {
				log.Fatal(err)
			}
			start_time := time.Now()
			session_res, err := ses.Exec(start_time, id, u.Devise_token)
			if err != nil || session_res == nil {
				log.Fatal(err)
			}

			fmt.Printf("StartTime: %v\n", time.Now())
			fmt.Println("User created Successfully!")

			user := models.Register{id, u.Firstname, u.Lastname, u.Email, u.Password, u.Password_confirmation, u.City, u.State, u.Country, u.Devise_token}

			b, err := json.Marshal(models.SignIn{
				Success: "true",
				Message: "User created Successfully!",
				User:    user,
				Session: models.Session{id, start_time},
			})

			if err != nil || res == nil {
				log.Fatal(err)
			}

			rw.Header().Set("Content-Type", "application/json")
			rw.Write(b)
			stmt.Close()
		}
		defer fetch_id.Close()
	}
create_user_end:
	db.Close()
}
