package serializers

type ItemSerializer struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	CategoryID  uint    `copier:"Category_id" json:"category_id"`
}
