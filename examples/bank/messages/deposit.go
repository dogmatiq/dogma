package messages

// Deposit is a command requesting that funds be deposited into a bank account.
type Deposit struct {
	TransactionID string
	AccountID     string
	Amount        uint64
}

// DepositStarted is an event indicating that the process of depositing funds
// into an account has begun.
type DepositStarted struct {
	TransactionID string
	AccountID     string
	Amount        uint64
}

// CreditAccountForDeposit is a command that credits a bank account with
// deposited funds.
type CreditAccountForDeposit struct {
	TransactionID string
	AccountID     string
	Amount        uint64
}

// AccountCreditedForDeposit is an event that indicates an account has been
// credited with funds from a deposit.
type AccountCreditedForDeposit struct {
	TransactionID string
	AccountID     string
	Amount        uint64
}
