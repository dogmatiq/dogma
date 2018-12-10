package messages

// Transfer is a command requesting that funds be transferred from one bank
// account to another.
type Transfer struct {
	TransactionID string
	FromAccountID string
	ToAccountID   string
	Amount        uint64
}

// TransferStarted is an event indicating that the process of transferring funds
// from one account to another has begun.
type TransferStarted struct {
	TransactionID string
	FromAccountID string
	ToAccountID   string
	Amount        uint64
}

// CreditAccountForTransfer is a command that credits a bank account with
// transferred funds.
type CreditAccountForTransfer struct {
	TransactionID string
	AccountID     string
	Amount        uint64
}

// AccountCreditedForTransfer is an event that indicates an account has been
// credited with funds from a transfer.
type AccountCreditedForTransfer struct {
	TransactionID string
	AccountID     string
	Amount        uint64
}

// DebitAccountForTransfer is a command that requests a bank account be debited
// for a transfer.
type DebitAccountForTransfer struct {
	TransactionID string
	AccountID     string
	Amount        uint64
}

// AccountDebitedForTransfer is an event that indicates an account has been
// debited funds for a transfer.
type AccountDebitedForTransfer struct {
	TransactionID string
	AccountID     string
	Amount        uint64
}

// TransferDeclined is an event that indicates a requested transfer has been
// declined due to insufficient funds.
type TransferDeclined struct {
	TransactionID string
	AccountID     string
	Amount        uint64
}
