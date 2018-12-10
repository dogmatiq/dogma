package bank

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/examples/bank/messages"
)

// WithdrawalProcess manages the process of withdrawing funds from an account.
var WithdrawalProcess dogma.ProcessMessageHandler = withdrawalProcess{}

type withdrawalProcess struct{}

func (withdrawalProcess) New() dogma.ProcessRoot {
	return nil
}

func (withdrawalProcess) RouteEvent(_ context.Context, m dogma.Message, _ bool) (string, bool, error) {
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

func (withdrawalProcess) HandleEvent(
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

func (withdrawalProcess) HandleTimeout(context.Context, dogma.ProcessScope, dogma.Message) error {
	panic(dogma.UnexpectedMessage)
}
