package controller

import (
	"github.com/betwins/gme-common/order"
	"github.com/betwins/gme-common/result"
	"github.com/betwins/gtools"
	"github.com/gin-gonic/gin"
	"gme/i18n_code"
	"gme/service"
)

type orderCtrl struct{}

var Order orderCtrl

// GetOrderss getAlbums	godoc
// @Summary		查询订单列表
// @Description	查询订单列表
// @Tags	订单管理
// @Accept	application/json
// @Produce json
// @Success 200 {obj ect} nil
// @Router /orders [get]
func (ctrl *orderCtrl) GetOrders(c *gin.Context) result.Result[any] {
	var req order.QueryListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		return i18n_code.ParamTypeError.MGErrorWithArgs(err.Error())
	}

	req.PageIndex = gtools.Ternary[int](req.PageIndex == 0, 1, req.PageIndex)
	req.PageSize = gtools.Ternary[int](req.PageSize == 0, 10, req.PageSize)

	list, totalCount, err := service.Order.GetEntityList(&req)
	if err != nil {
		return result.Error(-1, err.Error())
	}

	count := totalCount / req.PageSize
	if totalCount%req.PageSize > 0 {
		count += 1
	}

	return result.SuccessWithPage[[]order.Base](list, count, req.PageIndex, req.PageSize, totalCount).ToAny()
}
