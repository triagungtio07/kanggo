package schema

type Order struct {
	Base
	UserId    int64   `gorm:"not null"`
	ProductId int64   `gorm:"not null"`
	Amount    float64 `gorm:"not null"`
	Status    string  `gorm:"not null;type:varchar(10);default:'pending'"`
}

func (Order) TableName() string {
	return "orders"
}
