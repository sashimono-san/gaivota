package data

type User struct {
	ID        int
	SKU       string
	Email     string
	Name      string
	CreatedOn string
	UpdatedOn string
	DeletedOn string
}

var users = []*User{}
