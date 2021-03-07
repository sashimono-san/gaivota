package data

type Portfolio struct {
	ID        int
	SKU       string
	UserId    int
	Name      string
	CreatedOn string
	UpdatedOn string
	DeletedOn string
}

var potfolios = []*Portfolio{}
