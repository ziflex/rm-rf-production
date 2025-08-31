package accounts

type (
	AccountCreation struct {
		DocumentNumber string `json:"document_number" db:"document_number"`
	}

	Account struct {
		ID             int64  `json:"id" db:"id"`
		DocumentNumber string `json:"document_number" db:"document_number"`
	}
)
