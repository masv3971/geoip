package model

import "errors"

var (
	//ErrTravelerCantTravel return if no prior document is found for user
	ErrTravelerCantTravel = errors.New("Error_cant travel since no prior possition is found")
)
