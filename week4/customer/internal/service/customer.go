package service

import (
	"context"
	v12 "customer/api/v1"
	error2 "customer/api/v1/error"
	"customer/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type CustomerService struct {
	v12.UnimplementedCustomerServer

	uc  *biz.CustomerUsercase
	log *log.Helper
}

func NewCustomerService(uc *biz.CustomerUsercase, logger log.Logger) *CustomerService {
	return &CustomerService{uc: uc, log: log.NewHelper(logger)}
}

func (s *CustomerService) CreateCustomer(ctx context.Context, req *v12.CreateCustomerRequest) (*v12.CreateCustomerReply, error) {
	return &v12.CreateCustomerReply{}, nil
}
func (s *CustomerService) UpdateCustomer(ctx context.Context, req *v12.UpdateCustomerRequest) (*v12.UpdateCustomerReply, error) {
	return &v12.UpdateCustomerReply{}, nil
}
func (s *CustomerService) DeleteCustomer(ctx context.Context, req *v12.DeleteCustomerRequest) (*v12.DeleteCustomerReply, error) {
	return &v12.DeleteCustomerReply{}, nil
}
func (s *CustomerService) GetCustomer(ctx context.Context, req *v12.GetCustomerRequest) (*v12.GetCustomerReply, error) {

	s.log.WithContext(ctx).Infof("GetCustomer Received: %v", req.GetId())

	if req.GetId() == "error" {
		return nil, error2.ErrorUserNotFound("Customer not found: %s", req.GetId())
	}
	return &v12.GetCustomerReply{Message: "Hello " + req.GetId()}, nil

	return &v12.GetCustomerReply{}, nil
}
func (s *CustomerService) ListCustomer(ctx context.Context, req *v12.ListCustomerRequest) (*v12.ListCustomerReply, error) {
	return &v12.ListCustomerReply{}, nil
}
