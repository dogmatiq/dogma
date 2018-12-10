package messages

// Withdraw is a command requesting that funds be withdrawn from a bank account.
type Withdraw struct {
	TransactionID string
	AccountID     string
	Amount        uint64
}

// WithdrawalStarted is an event indicating that the process of withdrawing
// funds from an account has begun.
type WithdrawalStarted struct {
	TransactionID string
	AccountID     string
	Amount        uint64
}

// DebitAccountForWithdrawal is a command that requests a bank account be
// debited for a withdrawal.
type DebitAccountForWithdrawal struct {
	TransactionID string
	AccountID     string
	Amount        uint64
}

// AccountDebitedForWithdrawal is an event that indicates an account has been
// debited funds for a withdrawal.
type AccountDebitedForWithdrawal struct {
	TransactionID string
	AccountID     string
	Amount        uint64
}

// WithdrawalDeclined is an event that indicates a requested withdrawal has been
// declined due to insufficient funds.
type WithdrawalDeclined struct {
	TransactionID string
	AccountID     string
	Amount        uint64
}
