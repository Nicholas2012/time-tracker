package models

type User struct {
	ID             int
	PassportSerie  int
	PassportNumber int
	Name           string
	Surname        string
	Patronymic     string
	Address        string
}
