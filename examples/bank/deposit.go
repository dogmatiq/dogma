package bank

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/examples/bank/messages"
)

// DepositProcess manages the process of depositing funds into an account.
var DepositProcess dogma.ProcessMessageHandler = depositProcess{}

type depositProcess struct{}

func (depositProcess) New() dogma.ProcessRoot {
	return nil
}

func (depositProcess) RouteEvent(_ context.Context, m dogma.Message, _ bool) (string, bool, error) {
	switch x := m.(type) {
	case messages.DepositStarted:
		return x.TransactionID, true, nil
	case messages.AccountCreditedForDeposit:
		return x.TransactionID, true, nil
	default:
		return "", false, nil
	}
}

func (depositProcess) HandleEvent(
	_ context.Context,
	s dogma.ProcessScope,
	m dogma.Message,
) error {
	switch x := m.(type) {
	case messages.DepositStarted:
		s.Begin()
		s.ExecuteCommand(messages.CreditAccountForDeposit{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
		})

	case messages.AccountCreditedForDeposit:
		s.End()

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}

func (depositProcess) HandleTimeout(context.Context, dogma.ProcessScope, dogma.Message) error {
	panic(dogma.UnexpectedMessage)
}
