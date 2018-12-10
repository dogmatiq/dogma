package messages

// OpenAccount is a command requesting that a new bank account be opened.
type OpenAccount struct {
	AccountID string
	Name      string
}

// AccountOpened is an event indicating that a new bank account has been opened.
type AccountOpened struct {
	AccountID string
	Name      string
}
