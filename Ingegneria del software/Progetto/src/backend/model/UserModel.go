package model

type UserModel struct {
	ID         int
	Name       string
	Email      string
	Username   string
	Password   string
	Passphrase string
}

type OrderModel struct {
	ObjectID   int
	ObjectName string
	Quantity   int
	Status     string
	Price      float64
}
