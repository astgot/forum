package model

// Sessions ...
type Sessions struct {
	UserID    int64
	SessionID string
}

/* SessionID assigned to Users.ID, if Sessiontime has not expired, it will keep that session
Otherwise, it will delete existing session from DB, and will assign him new session      */
