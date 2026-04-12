package models

type User struct {
	ID       int64
	Login    string
	PassHash []byte
}

type App struct {
	ID   int
	Name string
	// TODO убрать
	Secret string
}
