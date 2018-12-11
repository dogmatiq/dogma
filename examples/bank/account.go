package bank

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/examples/bank/messages"
)

// AccountHandler implements the domain logic for a bank account.
//
// It centralizes all transactions that are applied to an account in order to
// enforce a strict no-overdraw policy.
var AccountHandler dogma.AggregateMessageHandler = accountHandler{}

type account struct {
	Balance uint64
}

func (a *account) ApplyEvent(m dogma.Message) {
	switch x := m.(type) {
	case messages.AccountCreditedForDeposit:
		a.Balance += x.Amount
	case messages.AccountCreditedForTransfer:
		a.Balance += x.Amount
	case messages.AccountDebitedForWithdrawal:
		a.Balance -= x.Amount
	case messages.AccountDebitedForTransfer:
		a.Balance -= x.Amount
	}
}

type accountHandler struct{}

func (accountHandler) New() dogma.AggregateRoot {
	return &account{}
}

// RouteCommand returns the ID of the aggregate that should receive m.
func (accountHandler) RouteCommand(m dogma.Message, _ bool) (string, bool) {
	switch x := m.(type) {
	case messages.OpenAccount:
		return x.AccountID, true
	case messages.CreditAccountForDeposit:
		return x.AccountID, true
	case messages.CreditAccountForTransfer:
		return x.AccountID, true
	case messages.DebitAccountForWithdrawal:
		return x.AccountID, true
	case messages.DebitAccountForTransfer:
		return x.AccountID, true
	default:
		return "", false
	}
}

// HandleCommand handles a domain command that has been routed to this aggregate.
func (accountHandler) HandleCommand(s dogma.AggregateScope, m dogma.Message) {
	switch x := m.(type) {
	case messages.OpenAccount:
		openAccount(s, x)
	case messages.CreditAccountForDeposit:
		creditForDeposit(s, x)
	case messages.CreditAccountForTransfer:
		creditForTransfer(s, x)
	case messages.DebitAccountForWithdrawal:
		debitForWithdrawal(s, x)
	case messages.DebitAccountForTransfer:
		debitForTransfer(s, x)
	default:
		panic(dogma.UnexpectedMessage)
	}
}

func openAccount(s dogma.AggregateScope, m messages.OpenAccount) {
	if !s.Create() {
		s.Log("account has already been opened")
		return
	}

	s.RecordEvent(messages.AccountOpened{
		AccountID: m.AccountID,
		Name:      m.Name,
	})
}

func creditForDeposit(s dogma.AggregateScope, m messages.CreditAccountForDeposit) {
	s.RecordEvent(messages.AccountCreditedForDeposit{
		TransactionID: m.TransactionID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
	})
}

func creditForTransfer(s dogma.AggregateScope, m messages.CreditAccountForTransfer) {
	s.RecordEvent(messages.AccountCreditedForTransfer{
		TransactionID: m.TransactionID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
	})
}

func debitForWithdrawal(s dogma.AggregateScope, m messages.DebitAccountForWithdrawal) {
	a := s.Root().(*account)

	if a.Balance >= m.Amount {
		s.RecordEvent(messages.AccountDebitedForWithdrawal{
			TransactionID: m.TransactionID,
			AccountID:     m.AccountID,
			Amount:        m.Amount,
		})
	} else {
		s.RecordEvent(messages.WithdrawalDeclined{
			TransactionID: m.TransactionID,
			AccountID:     m.AccountID,
			Amount:        m.Amount,
		})
	}
}

func debitForTransfer(s dogma.AggregateScope, m messages.DebitAccountForTransfer) {
	a := s.Root().(*account)

	if a.Balance >= m.Amount {
		s.RecordEvent(messages.AccountDebitedForTransfer{
			TransactionID: m.TransactionID,
			AccountID:     m.AccountID,
			Amount:        m.Amount,
		})
	} else {
		s.RecordEvent(messages.TransferDeclined{
			TransactionID: m.TransactionID,
			AccountID:     m.AccountID,
			Amount:        m.Amount,
		})
	}
}
