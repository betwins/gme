package service

import (
	"github.com/betwins/gme-common/order"
	"gme/dao"
	"gme/i18n_code"
)

type orderService struct{}

var Order orderService

func (s *orderService) GetEntityList(req *order.QueryListReq) ([]order.Base, int, error) {
	list, totalCount, err := dao.Order.GetEntityList(req)
	if err != nil {
		return nil, 0, i18n_code.DbQueryErr.Error()
	}

	return list, totalCount, nil
}
