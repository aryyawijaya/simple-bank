package account

import (
	"context"
	"database/sql"
	"net/http"

	mydb "github.com/aryyawijaya/simple-bank/db/sqlc"
	"github.com/aryyawijaya/simple-bank/modules"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type Store interface {
	CreateAccount(ctx context.Context, arg mydb.CreateAccountParams) (mydb.Account, error)
	GetAccount(ctx context.Context, id int64) (mydb.Account, error)
	ListAccounts(ctx context.Context, arg mydb.ListAccountsParams) ([]mydb.Account, error)
}

type AccountModule struct {
	store   Store
	wrapper modules.Wrapper
}

func NewAccountModule(store Store, wrapper modules.Wrapper) *AccountModule {
	return &AccountModule{
		store:   store,
		wrapper: wrapper,
	}
}

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (am *AccountModule) Create(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, am.wrapper.ErrResp(err))
		return
	}

	arg := mydb.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := am.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, am.wrapper.ErrResp(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, am.wrapper.ErrResp(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (am *AccountModule) Get(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, am.wrapper.ErrResp(err))
		return
	}

	account, err := am.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, am.wrapper.ErrResp(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, am.wrapper.ErrResp(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountsRequest struct {
	PageID   int32 `form:"pageId" binding:"required,min=1"`
	PageSize int32 `form:"pageSize" binding:"required,min=5,max=10"`
}

func (am *AccountModule) GetAll(ctx *gin.Context) {
	var req getAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, am.wrapper.ErrResp(err))
		return
	}

	arg := mydb.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := am.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, am.wrapper.ErrResp(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
