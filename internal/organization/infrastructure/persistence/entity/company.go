package entity

type CompanyModel struct {
	ID   int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Name string `gorm:"column:name;not null"`
}

func (CompanyModel) TableName() string {
	return "companies"
}
