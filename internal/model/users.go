package model

//Users ...
type Users struct {
	ID           int    `db:"id"`
	Firstname    string `db:"firstname"`
	Lastname     string `db:"lastname"`
	Username     string `db:"username"`
	Email        string `db:"email"`
	Password     string
	ConfirmPwd   string
	EncryptedPwd string `db:"password"`
	Errors       map[string]string
}

// Sessions -->
type Sessions struct {
	ID      int
	Session string
}
