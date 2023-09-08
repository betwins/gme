package i18n_code

import (
	"errors"
	"fmt"
	"github.com/betwins/gme-common/result"
	"github.com/betwins/gtools"
)

type ResultId string

const (
	UrlNotFound         ResultId = "404"
	SystemError         ResultId = "系统异常"
	DbConnectErr        ResultId = "数据库连接失败"
	DbInsertErr         ResultId = "数据库插入失败"
	DbUpdateErr         ResultId = "数据库更新失败"
	DbDeleteErr         ResultId = "数据库删除失败"
	DataNotFound        ResultId = "数据库查无数据"
	ParamLost           ResultId = "参数不可为空"
	ParamError          ResultId = "参数错误{}"
	ParamTypeError      ResultId = "参数{}类型错误"
	ConnectFail         ResultId = "网络连接失败"
	ServiceUnavailable  ResultId = "服务不存在"
	Success             ResultId = "success"
	DbQueryErr          ResultId = "数据库查询失败"
	ConcurrencyConflict ResultId = "并发操作冲突，请重试"
)

func (receiver ResultId) Error() error {
	return errors.New(string(receiver))
}

func (receiver ResultId) ErrorWithMsg(msg string) error {
	return errors.New(fmt.Sprintf("%s: %s", string(receiver), msg))
}

func (receiver ResultId) ErrorWithArgs(args ...any) error {
	return errors.New(gtools.Format(string(receiver), args...))
}

func (receiver ResultId) MGError() result.Result[any] {
	return result.Error(receiver.GetCode(), string(receiver))
}

func (receiver ResultId) MGErrorWithMsg(msg string) result.Result[any] {
	return result.Error(receiver.GetCode(), fmt.Sprintf("%s:%s", string(receiver), msg))
}

func (receiver ResultId) MGErrorWithArgs(args ...any) result.Result[any] {
	return result.Error(receiver.GetCode(), gtools.Format(string(receiver), args))
}

func (receiver ResultId) CheckParametersLost(params map[string]string, paramNames ...string) result.Result[any] {
	for _, param := range paramNames {
		v := params[param]
		if v == "" {
			return result.Error(receiver.GetCode(), gtools.Format(string(receiver)+":{}", param))
		}
	}
	return result.Success[any](nil)
}

func (receiver ResultId) GetCode() int {
	switch receiver {
	case Success:
		return 1
	default:
		return -1
	}
}
