package transactions

import (
	"strings"
	"time"
)

type (
	OperationType int

	TransactionCreation struct {
		AccountID     int64         `json:"account_id" db:"account_id"`
		OperationType OperationType `json:"operation_type" db:"operation_type"`
		Amount        float64       `json:"amount" db:"amount"`
	}

	Transaction struct {
		ID            int64         `json:"id" db:"id"`
		AccountID     int64         `json:"account_id" db:"account_id"`
		OperationType OperationType `json:"operation_type" db:"operation_type"`
		Amount        float64       `json:"amount" db:"amount"`
		EventDate     time.Time     `json:"event_date" db:"event_date"`
	}
)

const (
	OperationTypeUnknown OperationType = iota
	OperationTypePurchase
	OperationTypeInstallmentPurchase
	OperationTypeWithdrawal
	OperationTypePayment
)

func NewOperationType(op int) OperationType {
	switch op {
	case 1:
		return OperationTypePurchase
	case 2:
		return OperationTypeInstallmentPurchase
	case 3:
		return OperationTypeWithdrawal
	case 4:
		return OperationTypePayment
	default:
		return OperationTypeUnknown
	}
}

func NewOperationTypeFromString(s string) OperationType {
	switch strings.ToLower(s) {
	case "purchase":
		return OperationTypePurchase
	case "installment_purchase":
		return OperationTypeInstallmentPurchase
	case "withdrawal":
		return OperationTypeWithdrawal
	case "payment":
		return OperationTypePayment
	default:
		return OperationTypeUnknown
	}
}

func (o OperationType) String() string {
	switch o {
	case OperationTypePurchase:
		return "purchase"
	case OperationTypeInstallmentPurchase:
		return "installment_purchase"
	case OperationTypeWithdrawal:
		return "withdrawal"
	case OperationTypePayment:
		return "payment"
	default:
		return ""
	}
}
