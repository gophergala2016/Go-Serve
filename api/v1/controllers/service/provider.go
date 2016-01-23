package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gophergala2016/Go-Serve/api/v1/models"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
)

type serviceController struct{}

var Service serviceController

func (s serviceController) Create(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	flag := 1
	var u models.Service

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

	if flag == 1 {
		if u.Type == 0 || u.Description == "" || u.Address == "" || u.City == "" || u.State == "" || u.Country == "" {

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
			goto create_service_end
		}
	}

	if flag == 1 {
		fetch_id, err := db.Query("SELECT coalesce(max(id), 0) FROM service_provider")
		if err != nil {
			log.Fatal(err)
		}

		for fetch_id.Next() {
			var id int
			err = fetch_id.Scan(&id)

			if err != nil {
				log.Fatal(err)
			}
			id = id + 1

			var insert_service string = "insert into service_provider(id, type, description, experience, certificate, address, city, state,country) values ($1,$2,$3,$4,$5,$6,$7,$8,$9)"
			prepare_insert_service, err := db.Prepare(insert_service)
			if err != nil {
				log.Fatal(err)
			}
			user_res, err := prepare_insert_service.Exec(id, u.Type, u.Description, u.Experience, u.Certificate, u.Address, u.City, u.State, u.Country)
			if err != nil || user_res == nil {
				log.Fatal(err)
			}
		}

		b, err := json.Marshal(models.SuccessServiceMessage{
			Success: "false",
			Message: "Service created successfully",
		})
		if err != nil {
			log.Fatal(err)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(b)
		flag = 0
		goto create_service_end
	} else {
		b, err := json.Marshal(models.ErrorMessage{
			Success: "false",
			Error:   "Service cannot be created",
		})

		if err != nil {
			log.Fatal(err)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(b)
		fmt.Println("Session already Exist")
		goto create_service_end
	}

create_service_end:
	db.Close()
}
