package schema

type User struct {
	Base
	Name     string `gorm:"type:varchar(255);null"`
	Email    string `gorm:"type:varchar(255);unique;not null"`
	Password string `gorm:"type:varchar(255);not null"`
}

func (User) TableName() string {
	return "users"
}
