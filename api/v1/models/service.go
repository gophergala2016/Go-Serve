package models

type Service struct {
	User_id     int    `valid:"numeric"`
	Type        int    `valid:"numeric",required`
	Description string `valid:"alphanum,required"`
	Experience  int    `valid:"numeric"`
	Certificate bool
	Address     string `valid:"alphanum,required"`
	City        string `valid:"alphanum,required"`
	State       string `valid:"alphanum,required"`
	Country     string `valid:"alphanum"`
}

type User_Service struct {
	User_id      int
	Name         string
	Image        string
	MobileNumber string
	Type         int
	Description  string
	Experience   int
	Certificate  bool
	Address      string
	City         string
	State        string
	Country      string
}

type SuccessServiceMessage struct {
	Success string
	Message string
}

type ServiceList struct {
	Success         string
	No_Of_Service   int
	Service_Details []User_Service
}

type UserServeList struct {
	Success         string
	No_Of_Service   int
	Service_Details []Service
}
