package models

type TechType struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	ImgUrl string `json:"imgUrl"`
}

func (TechType) TableName() string {
	return "tech_type"
}
