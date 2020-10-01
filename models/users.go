package models

type Users struct {
	Nama   string `json:"nama"`
	Umur   int    `json:"umur"`
	Alamat string `json:"alamat"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}
