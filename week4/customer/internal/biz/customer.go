package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type Customer struct {
	Id string
}

type CustomerRepo interface {
	CreateCustomer(context.Context, *Customer) error
	UpdateCustomer(context.Context, *Customer) error
}

type CustomerUsercase struct {
	repo CustomerRepo
	log  *log.Helper
}

func NewCustomerUsercase(repo CustomerRepo, logger log.Logger) *CustomerUsercase {
	return &CustomerUsercase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *CustomerUsercase) Create(ctx context.Context, g *Customer) error {
	return uc.repo.CreateCustomer(ctx, g)
}

func (uc *CustomerUsercase) Update(ctx context.Context, g *Customer) error {
	return uc.repo.UpdateCustomer(ctx, g)
}
