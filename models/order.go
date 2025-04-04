// order.go
package models

type Order struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Quantity int    `json:"quantity"`
	BedName  string `json:"bed_name"`
}
