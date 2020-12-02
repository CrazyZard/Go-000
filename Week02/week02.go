package main

import (
	"database/sql"
	 xerrors "errors"
	"log"
	"net/http"
)

type user struct {
	id int
	name string
	age int
}

var RowNotFound = xerrors.New("row not found")

//dao层
func getUserById(id int)(*user,err) {
	_,err := dbGetUserById()
	if err != nil {
		return nil ,err
	}
	return &user{id, name, age}, nil
}


func getErrorCode(err error) int {
	switch errors.Cause(err) {
	case RowNotFound:
		return 404
	default:
		return 500
	}
}


//service
func service(id int) (*user,err) {
	user , err :=   getUserById(id)
	if err != nil {
		if errors.Is(err,RowNotFound) {
			return nil, errors.Wrapf(RowNotFound, "user [id=%d] not found", id)
		}
		return nil,err
	}
	return user,nil
}


func main()  {
	user ,err := service(12)
	if err != nil {
		return http.Response{Status:getErrorCode(err),Body: {"message":"系统错误，请稍后再试"}}
	}
	return http.Response{Body: user,Status: 200}
}

func dbGetUserById() ([]map[string]interface{}, error) {
	return nil, sql.ErrNoRows
}