package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gophergala2016/Go-Serve/api/v1/controllers"
	"github.com/gophergala2016/Go-Serve/api/v1/models"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
)

type listController struct{}

var List listController

func (s listController) Index(rw http.ResponseWriter, req *http.Request) {
	var l models.ServiceList
	flag := 1

	db, err := sql.Open("postgres", "password=password host=localhost dbname=go_service_development sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	get_service_list, err := db.Query("SELECT user_id, type, description, experience, certificate, address, city, state, country from service_provider")
	if err != nil || get_service_list == nil {
		log.Fatal(err)
	}
	no_of_issues := 0
	for get_service_list.Next() {
		var User_id int
		var Type int
		var Description string
		var Experience int
		var Certificate bool
		var Address string
		var City string
		var State string
		var Country string
		var Name string
		var Image string
		var Mobile_number string
		err := get_service_list.Scan(&User_id, &Type, &Description, &Experience, &Certificate, &Address, &City, &State, &Country)
		if err != nil {
			log.Fatal(err)
		}
		get_user_details, err := db.Query("SELECT name, image, mobile_number from users where id = $1", User_id)
		if err != nil || get_user_details == nil {
			log.Fatal(err)
		}
		for get_user_details.Next() {
			err := get_user_details.Scan(&Name, &Image, &Mobile_number)
			log.Fatal(err)
		}
	}
	l.Service_Details = append(l.Service_Details, models.User_Service{User_id, Name, Mobile_number, Image, Type, Description, Experience, Certificate, Address, City, State, Country})
	no_of_issues++
	if flag == 1 {
		b, err := json.Marshal(models.ServiceList{
			Success:         "true",
			No_Of_Service:   no_of_issues,
			Service_Details: l.Service_Details,
		})
		if err != nil {
			log.Fatal(err)
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.Write(b)
		goto index_end
	}
index_end:
}

func (s listController) Show(rw http.ResponseWriter, req *http.Request) {
	var l models.UserServeList
	vars := mux.Vars(req)
	id := vars["id"]
	user_id, err := strconv.Atoi(id)

	flag := 1

	db, err := sql.Open("postgres", "password=password host=localhost dbname=go_service_development sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	user_existance := controllers.Check_for_user(user_id)
	if user_existance == false {
		flag = 0
		b, err := json.Marshal(models.SuccessServiceMessage{
			Success: "false",
			Message: "user does not exist",
		})
		if err != nil {
			log.Fatal(err)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(b)
		goto user_service_list_end
	}

	if flag == 1 {
		user_existance := controllers.Check_for_user_session(user_id)
		if user_existance == false {
			flag = 0
			b, err := json.Marshal(models.SuccessServiceMessage{
				Success: "false",
				Message: "need to login first",
			})
			if err != nil {
				log.Fatal(err)
			}
			rw.Header().Set("Content-Type", "application/json")
			rw.Write(b)
			goto user_service_list_end
		}
	}

	if flag == 1 {
		get_service_list, err := db.Query("SELECT user_id, type, description, experience, certificate, address, city, state, country from service_provider where user_id = $1", user_id)
		if err != nil || get_service_list == nil {
			log.Fatal(err)
		}
		no_of_issues := 0
		for get_service_list.Next() {
			var User_id int
			var Type int
			var Description string
			var Experience int
			var Certificate bool
			var Address string
			var City string
			var State string
			var Country string

			err := get_service_list.Scan(&User_id, &Type, &Description, &Experience, &Certificate, &Address, &City, &State, &Country)
			if err != nil {
				log.Fatal(err)
			}
			l.Service_Details = append(l.Service_Details, models.Service{User_id, Type, Description, Experience, Certificate, Address, City, State, Country})
			no_of_issues++
		}

		if flag == 1 {
			b, err := json.Marshal(models.UserServeList{
				Success:         "true",
				No_Of_Service:   no_of_issues,
				Service_Details: l.Service_Details,
			})
			if err != nil {
				log.Fatal(err)
			}

			rw.Header().Set("Content-Type", "application/json")
			rw.Write(b)
			goto user_service_list_end
		}
	}
user_service_list_end:
}
