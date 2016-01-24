package models

type Service struct {
	User_id     int    `valid:"numeric"`
	Type        int    `valid:"numeric"`
	Description string `valid:"alphanum,required"`
	Experience  int    `valid:"numeric"`
	Certificate bool
	Address     string `valid:"alphanum,required"`
	City        string `valid:"alphanum,required"`
	State       string `valid:"alphanum,required"`
	Country     string `valid:"alphanum,required"`
}

type SuccessServiceMessage struct {
	Success string
	Message string
}

type ServiceList struct {
	Success         string
	No_Of_Service   int
	Service_Details []Service
}
