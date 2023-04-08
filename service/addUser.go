package service

import (
	"11pointer/database"
	"11pointer/logger"
	"context"
	"database/sql"
	namedParameterQuery "github.com/Knetic/go-namedParameterQuery"
	"go.uber.org/zap"
)

func addUser(request *addUserRequest) error {

	logger.LOG.Info("Got addUserRequest", zap.Any("request", request))

	err := addUserDatabase(context.Background(), request.Email)
	if err != nil {
		logger.LOG.Error("error in inserting user to database", zap.Any("request", request))
		return err
	}
	logger.LOG.Info("Successfully executed addUserRequest", zap.Any("request", request))
	return nil
}

type AddUserRequestVO struct {
	Email sql.NullString
}

func MakeAddUserRequestVO(email string) *AddUserRequestVO {

	return &AddUserRequestVO{
		Email: GetNullableString(email),
	}
}

func AddUserArgs(model *AddUserRequestVO) map[string]interface{} {

	args := make(map[string]interface{})
	args["email"] = model.Email

	return args
}

func addUserDatabase(ctx context.Context, email string) error {

	model := MakeAddUserRequestVO(email)
	args := AddUserArgs(model)
	query := "INSERT INTO users (`email`, `user_type`, `engagement_id`) VALUES (:email, 'FREE_USER', '1')"
	var rows sql.Result

	namedQuery := namedParameterQuery.NewNamedParameterQuery(query)
	namedQuery.SetValuesFromMap(args)
	err := database.Driver.GetDriver().Exec(ctx, namedQuery.GetParsedQuery(), namedQuery.GetParsedParameters(), &rows)
	if err != nil {
		logger.LOG.Error("Error could not addUserDatabase request", zap.Error(err))
		return err
	}

	_, err = rows.LastInsertId()
	if err != nil {
		logger.LOG.Error("Error could not get lastInsertedId for addUserDatabase", zap.Error(err))
		return err
	}
	return nil
}

func GetNullableString(nullableString string) sql.NullString {

	var result sql.NullString
	if len(nullableString) > 0 {
		result = sql.NullString{String: nullableString, Valid: true}
	} else {
		result = sql.NullString{}
	}
	return result
}
