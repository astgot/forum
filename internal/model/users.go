package model

//Users ...
type Users struct {
	ID           int    `db:"id"`
	Firstname    string `db:"firstname"`
	Lastname     string `db:"lastname"`
	Username     string `db:"username"`
	Email        string `db:"email"`
	Password     string `db:"password"`
	ConfirmPwd   string
	EncryptedPwd string
	Errors       map[string]string
}
