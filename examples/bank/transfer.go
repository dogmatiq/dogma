package bank

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/examples/bank/messages"
)

// TransferProcess manages the process of transferring funds between accounts.
var TransferProcess dogma.ProcessMessageHandler = transferProcess{}

type transfer struct {
	ToAccountID string
}

type transferProcess struct{}

func (transferProcess) New() dogma.ProcessRoot {
	return &transfer{}
}

func (transferProcess) RouteEvent(_ context.Context, m dogma.Message, _ bool) (string, bool, error) {
	switch x := m.(type) {
	case messages.TransferStarted:
		return x.TransactionID, true, nil
	case messages.AccountDebitedForTransfer:
		return x.TransactionID, true, nil
	case messages.AccountCreditedForTransfer:
		return x.TransactionID, true, nil
	case messages.TransferDeclined:
		return x.TransactionID, true, nil
	default:
		return "", false, nil
	}
}

func (transferProcess) HandleEvent(
	_ context.Context,
	s dogma.ProcessScope,
	m dogma.Message,
) error {
	switch x := m.(type) {
	case messages.TransferStarted:
		s.Begin()

		xfer := s.Root().(*transfer)
		xfer.ToAccountID = x.ToAccountID

		s.ExecuteCommand(messages.DebitAccountForTransfer{
			TransactionID: x.TransactionID,
			AccountID:     x.FromAccountID,
			Amount:        x.Amount,
		})

	case messages.AccountDebitedForTransfer:
		xfer := s.Root().(*transfer)

		s.ExecuteCommand(messages.CreditAccountForTransfer{
			TransactionID: x.TransactionID,
			AccountID:     xfer.ToAccountID,
			Amount:        x.Amount,
		})

	case messages.AccountCreditedForTransfer, messages.TransferDeclined:
		s.End()

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}

func (transferProcess) HandleTimeout(context.Context, dogma.ProcessScope, dogma.Message) error {
	panic(dogma.UnexpectedMessage)
}
