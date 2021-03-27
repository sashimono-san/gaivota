package data

type Portfolio struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user"`
	Name      string `json:"name"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
	DeletedAt string `json:"-"`
}

var potfolios = []*Portfolio{}
