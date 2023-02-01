package request

type ProductRequest struct {
	ID          int    `json:"-" form:"id,omitempty" `
	Name        string `json:"name" validate:"required" form:"name" bind:"required"`
	ProductCode int    `json:"prod_code" validate:"required" form:"prod_code" bind:"required"`
	Quantity    int    `json:"quantity" validate:"required" form:"quantity" bind:"required"`
}
