package model

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var rxEmail = regexp.MustCompile("")
var rxUname = regexp.MustCompile("")
var rxPwd = regexp.MustCompile("")

// Validate ...
func (u *Users) Validate() bool {
	var matchL, matchF bool = true, true
	u.Errors = make(map[string]string)
	// Validity of Email
	matchEmail := rxEmail.Match([]byte(u.Email))
	if !matchEmail {
		u.Errors["Email"] = "Please enter a valid e-mail address, e.g. yourmail@example.com"
	}
	// Username
	matchUname := rxUname.Match([]byte(u.Username))
	if !matchUname {
		u.Errors["Username"] = "Please use latin alphabet, lowercase and uppercase characters with numbers or special characters"
	}
	// Password
	matchPwd := rxPwd.Match([]byte(u.Password))
	if !matchPwd {
		u.Errors["Password"] = "The password must include one lowercase, uppercase, special characters and number"
	}
	// First and Last name
	if u.Firstname != "" && u.Lastname != "" {
		matchF = rxUname.Match([]byte(u.Firstname))
		matchL = rxUname.Match([]byte(u.Lastname))
		if !matchF || !matchL {
			u.Errors["FLName"] = "Please use latin alphabet"
		}
	}
	if u.Password != u.ConfirmPwd {
		u.Errors["Confirm"] = "The password confirmation does not match"
	}

	return len(u.Errors) == 0
}

// HashPassword ...
func HashPassword(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	return string(hash)
}

// UnameOrEmail -->
func UnameOrEmail(query string) bool {
	for _, v := range query {
		if v == '@' {
			return true
		}
	}
	return false
}

// ComparePassword --> Decrypt Password
func ComparePassword(hashpwd, pwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashpwd), []byte(pwd)); err != nil {
		return false
	}
	return true
}
