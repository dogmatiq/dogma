package bank

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/examples/bank/messages"
)

// DepositProcessHandler manages the process of depositing funds into an account.
var DepositProcessHandler dogma.ProcessMessageHandler = depositProcessHandler{}

type depositProcessHandler struct{}

func (depositProcessHandler) New() dogma.ProcessRoot {
	return nil
}

func (depositProcessHandler) Configure(c dogma.ProcessConfigurer) {
	c.RouteEventType(messages.DepositStarted{})
	c.RouteEventType(messages.AccountCreditedForDeposit{})
}

func (depositProcessHandler) RouteEventToInstance(_ context.Context, m dogma.Message) (string, bool, error) {
	switch x := m.(type) {
	case messages.DepositStarted:
		return x.TransactionID, true, nil
	case messages.AccountCreditedForDeposit:
		return x.TransactionID, true, nil
	default:
		panic(dogma.UnexpectedMessage)
	}
}

func (depositProcessHandler) HandleEvent(
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

func (depositProcessHandler) HandleTimeout(context.Context, dogma.ProcessScope, dogma.Message) error {
	panic(dogma.UnexpectedMessage)
}
