package medicine

type Medicine struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Manufacturer string  `json:"manufacturer"`
	Price        float64 `json:"price"`
	Prescription bool    `json:"prescription"`
	Quantity     int     `json:"quantity"`
}
