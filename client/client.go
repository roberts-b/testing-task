package client

import (
	"errors"
	"github.com/imdario/mergo"
	"net/mail"
)

type Client struct {
	Id     		*int `json:"id"`
	Name   		*string `json:"name"`
	Surname     *string `json:"surname"`
	Email 		*string `json:"email"`
}

func (client *Client)Validate() (e error) {

	if client.Id == nil {
		e = errors.New("Missing id")
		return
	}
	if client.Name == nil {
		e = errors.New("Missing name")
		return
	}
	if client.Surname == nil {
		e = errors.New("Missing surname")
		return
	}
	if client.Email == nil {
		e = errors.New("Missing email")
		return
	}

	_, err := mail.ParseAddress(*client.Email)
	if err != nil {
		e = errors.New("Provided email has wrong format")
		return
	}

	return
}

func (client *Client)Merge(source Client) (err error) {
	err = mergo.Merge(client, source, mergo.WithOverride)
	return
}
