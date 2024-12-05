package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	mockdb "github.com/ilhamgepe/simplebank/db/mock"
	db "github.com/ilhamgepe/simplebank/db/sqlc"
	"github.com/ilhamgepe/simplebank/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "Invalid ID",
			accountID: -100,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Contains(t, recorder.Body.String(), "Invalid ID")
			},
		},
		{
			name:      "Not Found",
			accountID: 9999,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				require.Contains(t, recorder.Body.String(), "No data found")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			// defer ctrl.Finish() // no need to call finish, itu kata package nya kalo di klik finish wkwkw

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreateAccountAPI(t *testing.T) {
	pgErr := &pgconn.PgError{
		Code:    "23505",
		Message: "duplicate key value violates unique constraint",
	}
	testCases := []struct {
		name          string
		requestBody   createAccountRequest
		buildStub     func(store *mockdb.MockStore, req createAccountRequest)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, req createAccountRequest)
	}{
		{
			name: "OK",
			requestBody: createAccountRequest{
				Owner:    "ilham",
				Currency: "USD",
			},
			buildStub: func(store *mockdb.MockStore, req createAccountRequest) {
				store.EXPECT().
					CreateAccount(gomock.Any(), db.CreateAccountParams{
						Owner:    req.Owner,
						Balance:  0,
						Currency: req.Currency,
					}).
					Times(1).
					Return(db.Account{
						ID:       1,
						Owner:    req.Owner,
						Balance:  0,
						Currency: req.Currency,
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, req createAccountRequest) {
				requireBodyMatchAccount(t, recorder.Body, db.Account{
					ID:       1,
					Owner:    req.Owner,
					Balance:  0,
					Currency: req.Currency,
				})
			},
		},
		{
			name: "bad request",
			requestBody: createAccountRequest{
				Owner:    "",
				Currency: "USDT",
			},
			buildStub: func(store *mockdb.MockStore, req createAccountRequest) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, req createAccountRequest) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				require.Contains(t, recorder.Body.String(), "required")
			},
		},
		{
			name: "duplicate entry",
			requestBody: createAccountRequest{
				Owner:    "owner",
				Currency: "USD",
			},
			buildStub: func(store *mockdb.MockStore, req createAccountRequest) {
				store.EXPECT().
					CreateAccount(gomock.Any(), db.CreateAccountParams{
						Owner:    req.Owner,
						Currency: req.Currency,
						Balance:  0,
					}).
					Times(1).
					Return(db.Account{}, pgErr)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, req createAccountRequest) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
				require.Contains(t, recorder.Body.String(), "Duplicate entry detected")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			tc.buildStub(store, tc.requestBody)

			server := newTestServer(t, store)

			recoder := httptest.NewRecorder()

			url := "/accounts"

			body, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			server.router.ServeHTTP(recoder, request)

			tc.checkResponse(t, recoder, tc.requestBody)
		})
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	b, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotResponse Response
	err = json.Unmarshal(b, &gotResponse)
	require.NoError(t, err)

	var gotAccount db.Account
	switch v := gotResponse.Data.(type) {
	case map[string]interface{}:
		databytes, err := json.Marshal(v)
		require.NoError(t, err)

		err = json.Unmarshal(databytes, &gotAccount)
		require.NoError(t, err)
	}

	require.Equal(t, account, gotAccount)
}

func randomAccount() db.Account {
	return db.Account{
		ID:       int64(utils.RandomInt(1, 1000)),
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
}

func randomAccounts(count int) []db.Account {
	var accounts []db.Account
	for i := 0; i < count; i++ {
		account := db.Account{
			ID:       int64(utils.RandomInt(1, 1000)),
			Owner:    utils.RandomOwner(),
			Balance:  utils.RandomMoney(),
			Currency: utils.RandomCurrency(),
		}
		accounts = append(accounts, account)
	}

	return accounts
}
