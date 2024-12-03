package api

// import (
// 	"bytes"
// 	"encoding/json"
// 	"io"
// 	"log"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	mockdb "github.com/ilhamgepe/simplebank/db/mock"
// 	db "github.com/ilhamgepe/simplebank/db/sqlc"
// 	"github.com/ilhamgepe/simplebank/utils"
// 	"github.com/jackc/pgx/v5/pgconn"
// 	"github.com/stretchr/testify/require"
// 	"go.uber.org/mock/gomock"
// )

// func TestCreateUserAPI(t *testing.T) {
// 	user, password := randomUser(t)

// 	testCases := []struct {
// 		name          string
// 		body          createUserRequest
// 		buildStubs    func(store *mockdb.MockStore)
// 		checkResponse func(recoder *httptest.ResponseRecorder)
// 	}{
// 		// {
// 		// 	name: "OK",
// 		// 	body: createUserRequest{
// 		// 		Username: user.Username,
// 		// 		Password: password,
// 		// 		FullName: user.FullName,
// 		// 		Email:    user.Email,
// 		// 	},
// 		// 	buildStubs: func(store *mockdb.MockStore) {
// 		// 		store.EXPECT().
// 		// 			CreateUser(gomock.Any(), gomock.Any()).
// 		// 			Times(1).
// 		// 			Return(user, nil)
// 		// 	},
// 		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
// 		// 		require.Equal(t, http.StatusOK, recorder.Code)
// 		// 		requireBodyMatchUser(t, recorder.Body, user)
// 		// 	},
// 		// },
// 		// {
// 		// 	name: "InternalError",
// 		// 	body: createUserRequest{
// 		// 		Username: user.Username,
// 		// 		Password: password,
// 		// 		FullName: user.FullName,
// 		// 		Email:    user.Email,
// 		// 	},
// 		// 	buildStubs: func(store *mockdb.MockStore) {
// 		// 		store.EXPECT().
// 		// 			CreateUser(gomock.Any(), gomock.Any()).
// 		// 			Times(1).
// 		// 			Return(db.User{}, sql.ErrConnDone)
// 		// 	},
// 		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
// 		// 		require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 		// 	},
// 		// },
// 		{
// 			name: "DuplicateUsername",
// 			body: createUserRequest{
// 				Username: user.Username,
// 				Password: password,
// 				FullName: user.FullName,
// 				Email:    user.Email,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					CreateUser(gomock.Any(), gomock.Any()).
// 					Times(1).
// 					Return(db.User{}, &pgconn.PgError{Code: "23505"})
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				log.Printf("recorder.Code: %+v", recorder)
// 				require.Equal(t, http.StatusForbidden, recorder.Code)
// 			},
// 		},
// 		// {
// 		// 	name: "InvalidUsername",
// 		// 	body: createUserRequest{
// 		// 		Username: "invalid-user#1",
// 		// 		Password: password,
// 		// 		FullName: user.FullName,
// 		// 		Email:    user.Email,
// 		// 	},
// 		// 	buildStubs: func(store *mockdb.MockStore) {
// 		// 		store.EXPECT().
// 		// 			CreateUser(gomock.Any(), gomock.Any()).
// 		// 			Times(0)
// 		// 	},
// 		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
// 		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
// 		// 	},
// 		// },
// 		// {
// 		// 	name: "InvalidEmail",
// 		// 	body: createUserRequest{
// 		// 		Username: user.Username,
// 		// 		Password: password,
// 		// 		FullName: user.FullName,
// 		// 		Email:    "invalid-email",
// 		// 	},
// 		// 	buildStubs: func(store *mockdb.MockStore) {
// 		// 		store.EXPECT().
// 		// 			CreateUser(gomock.Any(), gomock.Any()).
// 		// 			Times(0)
// 		// 	},
// 		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
// 		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
// 		// 	},
// 		// },
// 		// {
// 		// 	name: "TooShortPassword",
// 		// 	body: createUserRequest{
// 		// 		Username: user.Username,
// 		// 		Password: "2",
// 		// 		FullName: user.FullName,
// 		// 		Email:    user.Email,
// 		// 	},
// 		// 	buildStubs: func(store *mockdb.MockStore) {
// 		// 		store.EXPECT().
// 		// 			CreateUser(gomock.Any(), gomock.Any()).
// 		// 			Times(0)
// 		// 	},
// 		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
// 		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
// 		// 	},
// 		// },
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(store)

// 			server := NewServer(store)
// 			recorder := httptest.NewRecorder()

// 			// Marshal body data to JSON
// 			data, err := json.Marshal(tc.body)
// 			require.NoError(t, err)

// 			url := "/users"
// 			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
// 			require.NoError(t, err)

// 			server.router.ServeHTTP(recorder, request)
// 			tc.checkResponse(recorder)
// 		})
// 	}
// }

// func randomUser(t *testing.T) (user db.User, password string) {
// 	password = utils.RandomString(6)
// 	hashedPassword, err := utils.HashPassword(password)
// 	require.NoError(t, err)

// 	user = db.User{
// 		Username: utils.RandomOwner(),
// 		Password: hashedPassword,
// 		FullName: utils.RandomOwner(),
// 		Email:    utils.RandomEmail(),
// 	}
// 	return
// }

// func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
// 	data, err := io.ReadAll(body)
// 	require.NoError(t, err)

// 	var gotUser db.User
// 	err = json.Unmarshal(data, &gotUser)

// 	require.NoError(t, err)
// 	require.Equal(t, user.Username, gotUser.Username)
// 	require.Equal(t, user.FullName, gotUser.FullName)
// 	require.Equal(t, user.Email, gotUser.Email)
// 	require.Empty(t, gotUser.Password)
// }
