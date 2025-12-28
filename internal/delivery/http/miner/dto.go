package miner

type CreateMinerDTO struct {
	Name   string `json:"name" validate:"required"`
	Energy int    `json:"energy" validate:"required"`
	Age    int    `json:"age" validate:"required"`
}

type BuyEquipmentDTO struct {
	Name string `json:"name" validate:"required"`
}
