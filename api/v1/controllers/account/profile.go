package account

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gophergala2016/Go-Serve/api/v1/models"
	"github.com/gophergala2016/Go-Serve/api/v1/controllers"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
)

type profileController struct{}

var  Profile profileController

func (profile profileController) Create(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	flag := 1
	var p models.Profile

	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &p)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", "password=password host=localhost dbname=go_service_development sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if flag == 1 {
		if p.User_id == 0 || p.Image == "" || p.Name == "" || p.Mobile_number == "" || p.Age == "" || p.Gender == "" {

			result, err := govalidator.ValidateStruct(p)
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
			goto create_profile_end
		}
	}
 	if flag == 1 {
	 	user_existance := controllers.Check_for_user(p.User_id)
		if user_existance == false {
			flag = 0
			b, err := json.Marshal(models.SuccessServiceMessage{
			Success: "false",
			Message: "user does not exist"  ,
		})
		if err != nil {
			log.Fatal(err)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(b)
		goto create_profile_end
		}

 	}
	
	if flag == 1 {
		user_existance := controllers.Check_for_user_session(p.User_id)
		if user_existance == false {
			flag = 0
			b, err := json.Marshal(models.SuccessServiceMessage{
			Success: "false",
			Message: "need to login first"  ,
		})
		if err != nil {
			log.Fatal(err)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(b)
		goto create_profile_end
		}

	}

	if flag == 1 {
   	update_profile, err := db.Query("UPDATE users set name = $1, mobile_number = $2, image = $3, age = $4, gender = $5 where id = $6", p.Name, p.Mobile_number, p.Image, p.Age, p.Gender, p.User_id)
		if err != nil || update_profile == nil {
			log.Fatal(err)
			}
		b, err := json.Marshal(models.SuccessServiceMessage{
			Success: "true",
			Message: "profile updated successfully"  ,
		})
		if err != nil {
			log.Fatal(err)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(b)
		goto create_profile_end
	}

	if flag == 0 {
		b, err := json.Marshal(models.SuccessServiceMessage{
			Success: "false",
			Message: "Unable to update profile"  ,
		})
		if err != nil {
			log.Fatal(err)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(b)
		goto create_profile_end
	}


create_profile_end:
	db.Close()
}
