package schema

type Product struct {
	Base
	Name  string  `gorm:"type:varchar(255);not null"`
	Price float64 `gorm:"not null"`
	Qty   int64   `gorm:"not null"`
}

func (Product) TableName() string {
	return "products"
}
