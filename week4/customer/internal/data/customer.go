package data

import (
	"context"
	"customer/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type customerRepo struct {
	data *Data
	log  *log.Helper
}

// NewCustomerRepo .
func NewCustomerRepo(data *Data, logger log.Logger) biz.CustomerRepo {
	return &customerRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *customerRepo) CreateCustomer(ctx context.Context, g *biz.Customer) error {
	return nil
}

func (r *customerRepo) UpdateCustomer(ctx context.Context, g *biz.Customer) error {
	return nil
}
