package user_repository_impl_pb

import (
	"context"

	"github.com/gofrs/uuid"
	pb_user_internal "github.com/krissukoco/go-food-order-microservices/api/proto/user_service/user_internal"
	"github.com/krissukoco/go-food-order-microservices/internal/services/auth/entity"
	user_repository "github.com/krissukoco/go-food-order-microservices/internal/services/auth/repository/user"
	"github.com/krissukoco/go-food-order-microservices/internal/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type repository struct {
	userClient   pb_user_internal.UserServiceInternalClient
	grpcCallOpts []grpc.CallOption
}

func New(
	userClient pb_user_internal.UserServiceInternalClient,
	grpcCallOpts ...grpc.CallOption,
) user_repository.Repository {
	return &repository{userClient, grpcCallOpts}
}

func (r *repository) GetPosUser(ctx context.Context, id uuid.UUID) (*entity.PosUser, error) {
	res, err := r.userClient.GetPosUser(ctx, &pb_user_internal.Id{Id: id.String()}, r.grpcCallOpts...)
	if err != nil {
		errResp, code := transport.ParseGRPCError(err)
		if code == codes.NotFound {
			return nil, entity.ErrNotFound
		}
		return nil, errResp
	}
	return posUserFromPb(res), nil
}

func (r *repository) GetPosUserByEmail(ctx context.Context, email string) (*entity.PosUser, error) {
	res, err := r.userClient.GetPosUserByEmail(ctx, &pb_user_internal.EmailReq{Email: email}, r.grpcCallOpts...)
	if err != nil {
		errResp, code := transport.ParseGRPCError(err)
		if code == codes.NotFound {
			return nil, entity.ErrNotFound
		}
		return nil, errResp
	}
	return posUserFromPb(res), nil
}

func (r *repository) GetBackofficeUser(ctx context.Context, id uuid.UUID) (*entity.BackofficeUser, error) {
	res, err := r.userClient.GetBackofficeUser(ctx, &pb_user_internal.Id{Id: id.String()}, r.grpcCallOpts...)
	if err != nil {
		errResp, code := transport.ParseGRPCError(err)
		if code == codes.NotFound {
			return nil, entity.ErrNotFound
		}
		return nil, errResp
	}
	return backofficeUserFromPb(res), nil
}

func (r *repository) GetBackofficeUserByEmail(ctx context.Context, email string) (*entity.BackofficeUser, error) {
	res, err := r.userClient.GetBackofficeUserByEmail(ctx, &pb_user_internal.EmailReq{Email: email}, r.grpcCallOpts...)
	if err != nil {
		errResp, code := transport.ParseGRPCError(err)
		if code == codes.NotFound {
			return nil, entity.ErrNotFound
		}
		return nil, errResp
	}
	return backofficeUserFromPb(res), nil
}

func (r *repository) UpsertCustomer(ctx context.Context, data entity.Customer) error {
	panic("unimplemented")
}

func posUserFromPb(in *pb_user_internal.PosUser) *entity.PosUser {
	if in == nil {
		return &entity.PosUser{}
	}
	return &entity.PosUser{
		Id:            uuid.FromStringOrNil(in.Id),
		RestaurantId:  uuid.FromStringOrNil(in.RestaurantId),
		Name:          in.Name,
		Email:         in.Email,
		Password:      in.Password,
		FirstPassword: in.FirstPassword,
	}
}

func backofficeUserFromPb(in *pb_user_internal.BackofficeUser) *entity.BackofficeUser {
	if in == nil {
		return &entity.BackofficeUser{}
	}
	return &entity.BackofficeUser{
		Id:            uuid.FromStringOrNil(in.Id),
		RestaurantId:  uuid.FromStringOrNil(in.RestaurantId),
		Name:          in.Name,
		Email:         in.Email,
		Password:      in.Password,
		FirstPassword: in.FirstPassword,
	}
}
