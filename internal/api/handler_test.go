package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ziflex/rm-rf-production/internal/api"
	"github.com/ziflex/rm-rf-production/internal/server"
	"github.com/ziflex/rm-rf-production/pkg/accounts"
	"github.com/ziflex/rm-rf-production/pkg/common"
	"github.com/ziflex/rm-rf-production/pkg/transactions"
	"github.com/ziflex/rm-rf-production/spec"
)

type mockAccountsService struct {
	mock.Mock
}

func (m *mockAccountsService) CreateAccount(ctx context.Context, creation accounts.AccountCreation) (accounts.Account, error) {
	args := m.Mock.Called(ctx, creation)

	return args.Get(0).(accounts.Account), args.Error(1)
}

func (m *mockAccountsService) GetAccountByID(ctx context.Context, id int64) (accounts.Account, error) {
	args := m.Mock.Called(ctx, id)

	return args.Get(0).(accounts.Account), args.Error(1)
}

type mockTransactionsService struct {
	mock.Mock
}

func (m *mockTransactionsService) CreateTransaction(ctx context.Context, creation transactions.TransactionCreation) (transactions.Transaction, error) {
	args := m.Mock.Called(ctx, creation)

	return args.Get(0).(transactions.Transaction), args.Error(1)
}

func createServer(accSvc accounts.Service, txSvc transactions.Service) (*server.Server, error) {
	logger := zerolog.New(io.Discard).With().Timestamp().Logger()

	return server.NewServer(api.NewHandler(
		accSvc,
		txSvc,
	), server.Options{
		Logger: logger,
		Spec:   spec.File,
	})
}

func toJSON(t *testing.T, v any) io.Reader {
	b, err := json.Marshal(v)
	assert.NoError(t, err)

	return bytes.NewReader(b)
}

func TestCreateAccount_Success(t *testing.T) {
	mockAccSvc := new(mockAccountsService)
	svr, err := createServer(mockAccSvc, &mockTransactionsService{})
	assert.NoError(t, err)

	go func() {
		if err := svr.Run(8080); err != nil && err != http.ErrServerClosed {
			t.Errorf("server error: %v", err)
		}
	}()

	time.Sleep(1 * time.Second)

	defer func() {
		if err := svr.Shutdown(); err != nil {
			t.Errorf("shutdown error: %v", err)
		}
	}()

	creation := accounts.AccountCreation{
		DocumentNumber: "12345678900",
	}
	mockAccSvc.On("CreateAccount", mock.Anything, creation).Return(accounts.Account{
		ID:             1,
		DocumentNumber: creation.DocumentNumber,
	}, nil)

	payload := toJSON(t, api.AccountCreateRequest{
		DocumentNumber: creation.DocumentNumber,
	})
	resp, err := http.Post("http://localhost:8080/accounts", "application/json", payload)

	assert.NoError(t, err)

	body, err := io.ReadAll(resp.Body)

	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result api.CreateAccount201JSONResponse

	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.AccountId)
	assert.Equal(t, creation.DocumentNumber, result.DocumentNumber)
	mockAccSvc.AssertExpectations(t)
}

func TestCreateAccount_Error_Validation(t *testing.T) {
	mockAccSvc := new(mockAccountsService)
	svr, err := createServer(mockAccSvc, &mockTransactionsService{})
	assert.NoError(t, err)

	go func() {
		if err := svr.Run(8080); err != nil && err != http.ErrServerClosed {
			t.Errorf("server error: %v", err)
		}
	}()

	time.Sleep(1 * time.Second)

	defer func() {
		if err := svr.Shutdown(); err != nil {
			t.Errorf("shutdown error: %v", err)
		}
	}()

	creation := accounts.AccountCreation{
		DocumentNumber: "",
	}

	payload := toJSON(t, creation)
	resp, err := http.Post("http://localhost:8080/accounts", "application/json", payload)

	assert.NoError(t, err)

	body, err := io.ReadAll(resp.Body)

	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var result api.Error

	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)
	assert.Equal(t, "badRequest", result.Code)
	mockAccSvc.AssertExpectations(t)
}

func TestGetAccountByID_Success(t *testing.T) {
	mockAccSvc := new(mockAccountsService)
	svr, err := createServer(mockAccSvc, &mockTransactionsService{})
	assert.NoError(t, err)

	go func() {
		if err := svr.Run(8080); err != nil && err != http.ErrServerClosed {
			t.Errorf("server error: %v", err)
		}
	}()

	time.Sleep(1 * time.Second)

	defer func() {
		if err := svr.Shutdown(); err != nil {
			t.Errorf("shutdown error: %v", err)
		}
	}()

	expected := accounts.Account{
		ID:             1,
		DocumentNumber: "12345678900",
	}

	mockAccSvc.On("GetAccountByID", mock.Anything, expected.ID).Return(expected, nil)

	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/accounts/%d", expected.ID))

	assert.NoError(t, err)

	body, err := io.ReadAll(resp.Body)

	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result api.GetAccount200JSONResponse

	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, result.AccountId)
	assert.Equal(t, expected.DocumentNumber, result.DocumentNumber)
	mockAccSvc.AssertExpectations(t)
}

func TestGetAccountByID_Error_NotFound(t *testing.T) {
	mockAccSvc := new(mockAccountsService)
	svr, err := createServer(mockAccSvc, &mockTransactionsService{})
	assert.NoError(t, err)

	go func() {
		if err := svr.Run(8080); err != nil && err != http.ErrServerClosed {
			t.Errorf("server error: %v", err)
		}
	}()

	time.Sleep(1 * time.Second)

	defer func() {
		if err := svr.Shutdown(); err != nil {
			t.Errorf("shutdown error: %v", err)
		}
	}()

	expected := accounts.Account{
		ID:             1,
		DocumentNumber: "12345678900",
	}

	mockAccSvc.On("GetAccountByID", mock.Anything, expected.ID).Return(accounts.Account{}, common.ErrNotFound)

	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/accounts/%d", expected.ID))

	assert.NoError(t, err)

	body, err := io.ReadAll(resp.Body)

	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var result api.Error

	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)
	assert.Equal(t, "notFound", result.Code)
	mockAccSvc.AssertExpectations(t)
}

func TestGetAccountByID_Error_Validation(t *testing.T) {
	mockAccSvc := new(mockAccountsService)
	svr, err := createServer(mockAccSvc, &mockTransactionsService{})
	assert.NoError(t, err)

	go func() {
		if err := svr.Run(8080); err != nil && err != http.ErrServerClosed {
			t.Errorf("server error: %v", err)
		}
	}()

	time.Sleep(1 * time.Second)

	defer func() {
		if err := svr.Shutdown(); err != nil {
			t.Errorf("shutdown error: %v", err)
		}
	}()

	resp, err := http.Get("http://localhost:8080/accounts/foobar")

	assert.NoError(t, err)

	body, err := io.ReadAll(resp.Body)

	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var result api.Error

	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)
	assert.Equal(t, "badRequest", result.Code)
	mockAccSvc.AssertExpectations(t)
}

func TestCreateTransaction_Success(t *testing.T) {
	mockTxSvc := new(mockTransactionsService)
	svr, err := createServer(&mockAccountsService{}, mockTxSvc)
	assert.NoError(t, err)

	go func() {
		if err := svr.Run(8080); err != nil && err != http.ErrServerClosed {
			t.Errorf("server error: %v", err)
		}
	}()

	time.Sleep(1 * time.Second)

	defer func() {
		if err := svr.Shutdown(); err != nil {
			t.Errorf("shutdown error: %v", err)
		}
	}()

	expected := transactions.Transaction{
		ID:            1,
		AccountID:     1,
		OperationType: transactions.OperationTypePurchase,
		Amount:        100.0,
		EventDate:     time.Now(),
	}

	input := transactions.TransactionCreation{
		AccountID:     expected.AccountID,
		OperationType: expected.OperationType,
		Amount:        expected.Amount,
	}

	mockTxSvc.On("CreateTransaction", mock.Anything, input).Return(expected, nil)

	payload := toJSON(t, api.TransactionCreateRequest{
		AccountId:       input.AccountID,
		OperationTypeId: api.OperationType(input.OperationType),
		Amount:          input.Amount,
	})
	resp, err := http.Post("http://localhost:8080/transactions", "application/json", payload)
	assert.NoError(t, err)

	body, err := io.ReadAll(resp.Body)

	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var creationResult api.CreateTransaction201JSONResponse
	err = json.Unmarshal(body, &creationResult)
	assert.NoError(t, err)

	assert.Equal(t, expected.ID, creationResult.TransactionId)
	assert.Equal(t, expected.AccountID, creationResult.AccountId)
	assert.Equal(t, int(expected.OperationType), int(creationResult.OperationTypeId))
	assert.Equal(t, -expected.Amount, -creationResult.Amount)
	mockTxSvc.AssertExpectations(t)
}

func TestCreateTransaction_Error_AccountIDNotFound(t *testing.T) {
	mockTxSvc := new(mockTransactionsService)
	svr, err := createServer(&mockAccountsService{}, mockTxSvc)
	assert.NoError(t, err)

	go func() {
		if err := svr.Run(8080); err != nil && err != http.ErrServerClosed {
			t.Errorf("server error: %v", err)
		}
	}()

	time.Sleep(1 * time.Second)

	defer func() {
		if err := svr.Shutdown(); err != nil {
			t.Errorf("shutdown error: %v", err)
		}
	}()

	input := transactions.TransactionCreation{
		AccountID:     1,
		OperationType: transactions.OperationTypePurchase,
		Amount:        100.0,
	}

	mockTxSvc.On("CreateTransaction", mock.Anything, input).Return(transactions.Transaction{}, common.ErrNotFound)

	payload := toJSON(t, api.TransactionCreateRequest{
		AccountId:       input.AccountID,
		OperationTypeId: api.OperationType(input.OperationType),
		Amount:          input.Amount,
	})
	resp, err := http.Post("http://localhost:8080/transactions", "application/json", payload)
	assert.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var creationResult api.Error
	err = json.Unmarshal(body, &creationResult)
	assert.NoError(t, err)
	assert.Equal(t, "notFound", creationResult.Code)
}

func TestCreateTransaction_Error_Validation(t *testing.T) {
	mockTxSvc := new(mockTransactionsService)
	svr, err := createServer(&mockAccountsService{}, mockTxSvc)
	assert.NoError(t, err)

	type testCase struct {
		name    string
		payload api.TransactionCreateRequest
	}

	tsdata := []testCase{
		{
			name: "Missing account ID",
			payload: api.TransactionCreateRequest{
				OperationTypeId: api.OperationType(transactions.OperationTypePurchase),
				Amount:          100.0,
			},
		},
		{
			name: "Missing operation type ID",
			payload: api.TransactionCreateRequest{
				AccountId: 1,
				Amount:    100.0,
			},
		},
		{
			name: "Missing amount",
			payload: api.TransactionCreateRequest{
				AccountId:       1,
				OperationTypeId: api.OperationType(transactions.OperationTypePurchase),
			},
		},
		{
			name: "Invalid operation type ID",
			payload: api.TransactionCreateRequest{
				AccountId:       1,
				OperationTypeId: 99,
				Amount:          100.0,
			},
		},
	}

	go func() {
		if err := svr.Run(8080); err != nil && err != http.ErrServerClosed {
			t.Errorf("server error: %v", err)
		}
	}()

	time.Sleep(1 * time.Second)

	defer func() {
		if err := svr.Shutdown(); err != nil {
			t.Errorf("shutdown error: %v", err)
		}
	}()

	for _, tc := range tsdata {
		t.Run(tc.name, func(t *testing.T) {
			payload := toJSON(t, tc.payload)
			resp, err := http.Post("http://localhost:8080/transactions", "application/json", payload)
			assert.NoError(t, err)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			var creationResult api.Error
			err = json.Unmarshal(body, &creationResult)
			assert.NoError(t, err)
			assert.Equal(t, "badRequest", creationResult.Code)
		})
	}
}
