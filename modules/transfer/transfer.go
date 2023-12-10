package transfer

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	mydb "github.com/aryyawijaya/simple-bank/db/sqlc"
	"github.com/aryyawijaya/simple-bank/modules"
	"github.com/gin-gonic/gin"
)

type Store interface {
	TransferTx(ctx context.Context, arg mydb.TransferTxParams) (mydb.TransferTxResult, error)
	GetAccount(ctx context.Context, id int64) (mydb.Account, error)
}

type TransferModule struct {
	store   Store
	wrapper modules.Wrapper
}

func NewTransferModule(store Store, wrapper modules.Wrapper) *TransferModule {
	return &TransferModule{
		store:   store,
		wrapper: wrapper,
	}
}

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (tm *TransferModule) CreateTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, tm.wrapper.ErrResp(err))
		return
	}

	if !tm.validAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}
	if !tm.validAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	arg := mydb.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := tm.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, tm.wrapper.ErrResp(err))
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

func (tm *TransferModule) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := tm.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, tm.wrapper.ErrResp(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, tm.wrapper.ErrResp(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, tm.wrapper.ErrResp(err))
		return false
	}

	return true
}
