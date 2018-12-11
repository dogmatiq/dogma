package bank

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/examples/bank/messages"
)

// TransactionHandler implements the domain logic for a transaction of any
// kind against an account.
//
// It's sole purpose is to ensure the global uniqueness of transaction IDs.
var TransactionHandler dogma.AggregateMessageHandler = transactionHandler{}

type transaction struct{}

func (t *transaction) ApplyEvent(dogma.Message) {
}

type transactionHandler struct{}

func (transactionHandler) New() dogma.AggregateRoot {
	return &transaction{}
}

func (transactionHandler) RouteCommand(m dogma.Message, _ bool) (string, bool) {
	switch x := m.(type) {
	case messages.Deposit:
		return x.TransactionID, true
	case messages.Withdraw:
		return x.TransactionID, true
	case messages.Transfer:
		return x.TransactionID, true
	default:
		return "", false
	}
}

func (transactionHandler) HandleCommand(s dogma.AggregateScope, m dogma.Message) {
	if !s.Create() {
		s.Log("transaction already exists")
		return
	}

	switch x := m.(type) {
	case messages.Deposit:
		s.RecordEvent(messages.DepositStarted{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
		})

	case messages.Withdraw:
		s.RecordEvent(messages.WithdrawalStarted{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
		})

	case messages.Transfer:
		s.RecordEvent(messages.TransferStarted{
			TransactionID: x.TransactionID,
			FromAccountID: x.FromAccountID,
			ToAccountID:   x.ToAccountID,
			Amount:        x.Amount,
		})

	default:
		panic(dogma.UnexpectedMessage)
	}
}
