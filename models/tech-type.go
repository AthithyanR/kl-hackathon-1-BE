package models

type TechType struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (TechType) TableName() string {
	return "tech_type"
}
