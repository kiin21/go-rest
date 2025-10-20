package entity

type CompanyEntity struct {
	ID   int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Name string `gorm:"column:name;not null"`
}

func (CompanyEntity) TableName() string {
	return "companies"
}
