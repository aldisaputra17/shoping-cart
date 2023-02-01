package entities

type Products struct {
	ID           int    `gorm:"primaryKey" json:"id"`
	Name         string `gorm:"type:varchar(255)" json:"name"`
	Product_Code int    `json:"product_code"`
	Quantity     int    `json:"quantity"`
}
