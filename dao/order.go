package dao

import (
	"github.com/betwins/gme-common/order"
	"gme/gmysql"
	"gme/i18n_code"
)

type orderDao struct{}

var Order orderDao

func (d *orderDao) GetEntityList(req *order.QueryListReq) ([]order.Base, int, error) {
	var list = make([]order.Base, 0)

	conn, err := gmysql.GetConnection()
	if err != nil {
		return list, 0, i18n_code.DbConnectErr.Error()
	}

	var totalCount int64
	conn = conn.Debug().Model(&order.Base{})
	if req.OrderId != "" {
		conn = conn.Where("order_id = ?", req.OrderId)
	}

	err = conn.Debug().Model(&order.Base{}).Count(&totalCount).Error
	if err != nil {
		return list, 0, err
	}
	if totalCount == 0 {
		return list, int(totalCount), nil
	}

	offsetPos := (req.PageIndex - 1) * req.PageSize
	err = conn.Debug().Model(&order.Base{}).Order("create_at DESC").Limit(req.PageSize).Offset(offsetPos).Find(&list).Error
	if err != nil {
		return list, 0, err
	}

	return list, int(totalCount), nil
}
