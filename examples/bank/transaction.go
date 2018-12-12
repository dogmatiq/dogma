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

func (transactionHandler) Configure(c dogma.AggregateConfigurer) {
	c.RouteCommandType(messages.Deposit{})
	c.RouteCommandType(messages.Withdraw{})
	c.RouteCommandType(messages.Transfer{})
}

func (transactionHandler) RouteCommandToInstance(m dogma.Message) string {
	switch x := m.(type) {
	case messages.Deposit:
		return x.TransactionID
	case messages.Withdraw:
		return x.TransactionID
	case messages.Transfer:
		return x.TransactionID
	default:
		panic(dogma.UnexpectedMessage)
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
