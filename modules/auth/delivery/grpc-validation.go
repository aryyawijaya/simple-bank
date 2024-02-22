package authdelivery

import (
	"github.com/aryyawijaya/simple-bank/pb"
	utilgrpc "github.com/aryyawijaya/simple-bank/util/grpc"
	"github.com/aryyawijaya/simple-bank/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func validateLoginRequest(req *pb.LoginRequest) []*errdetails.BadRequest_FieldViolation {
	violations := []*errdetails.BadRequest_FieldViolation{}

	// Username
	if err := validator.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, utilgrpc.NewFieldViolation("username", err))
	}

	// Password
	if err := validator.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, utilgrpc.NewFieldViolation("password", err))
	}

	return violations
}
