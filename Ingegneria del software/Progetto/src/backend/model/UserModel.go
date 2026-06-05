package Model
import ("fmt", "backend/domain")

struct UserModel {
	ID int
	Name string
	Email string
	Username string
	Password string
	Passphrase string
}

struct OrderModel {
	ObjectID int
	ObjectName string
	Quantity int
	Status string
	Price float64
}
