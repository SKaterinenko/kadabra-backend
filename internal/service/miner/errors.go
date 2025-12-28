package service

import "errors"

var (
	NotBool      = errors.New("Принимается только true или false")
	NotEquipment = errors.New("Такого оборудования нет")
)
