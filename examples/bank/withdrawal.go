package bank

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/examples/bank/messages"
)

// WithdrawalProcessHandler manages the process of withdrawing funds from an account.
var WithdrawalProcessHandler dogma.ProcessMessageHandler = withdrawalProcessHandler{}

type withdrawalProcessHandler struct{}

func (withdrawalProcessHandler) New() dogma.ProcessRoot {
	return nil
}

func (withdrawalProcessHandler) RouteEvent(_ context.Context, m dogma.Message, _ bool) (string, bool, error) {
	switch x := m.(type) {
	case messages.WithdrawalStarted:
		return x.TransactionID, true, nil
	case messages.AccountDebitedForWithdrawal:
		return x.TransactionID, true, nil
	case messages.WithdrawalDeclined:
		return x.TransactionID, true, nil
	default:
		return "", false, nil
	}
}

func (withdrawalProcessHandler) HandleEvent(
	_ context.Context,
	s dogma.ProcessScope,
	m dogma.Message,
) error {
	switch x := m.(type) {
	case messages.WithdrawalStarted:
		s.Begin()
		s.ExecuteCommand(messages.DebitAccountForWithdrawal{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
		})

	case messages.AccountDebitedForWithdrawal, messages.WithdrawalDeclined:
		s.End()

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}

func (withdrawalProcessHandler) HandleTimeout(context.Context, dogma.ProcessScope, dogma.Message) error {
	panic(dogma.UnexpectedMessage)
}
