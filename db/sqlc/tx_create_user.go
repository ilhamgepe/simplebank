package db

import "context"

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

type CreateUserTxResult struct {
	User User
}

// ini membuat catatn transfer, menambahkan account entries, update balance dalam satu kali transaction
func (s *store) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		err = arg.AfterCreate(result.User)

		return err
	})

	return result, err
}
