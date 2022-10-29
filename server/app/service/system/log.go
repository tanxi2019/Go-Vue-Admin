package system

import (
	"fmt"
	"server/app/model/system"
	"server/app/model/system/reqo"
	"server/global"
	"strings"
)

type LogService interface {
	GetOperationLogs(req *reqo.OperationLogListRequest) ([]system.OperationLog, int64, error)
	BatchDeleteOperationLogByIds(ids []uint) error
	SaveOperationLogChannel(olc <-chan *system.OperationLog) //处理OperationLogChan将日志记录到数据库
}

type Log struct{}

func NewLogService() LogService {
	return Log{}
}

func (l Log) GetOperationLogs(req *reqo.OperationLogListRequest) ([]system.OperationLog, int64, error) {
	var list []system.OperationLog
	db := global.DB.Model(&system.OperationLog{}).Order("start_time DESC")

	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	ip := strings.TrimSpace(req.Ip)
	if ip != "" {
		db = db.Where("ip LIKE ?", fmt.Sprintf("%%%s%%", ip))
	}
	path := strings.TrimSpace(req.Path)
	if path != "" {
		db = db.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	status := req.Status
	if status != 0 {
		db = db.Where("status = ?", status)
	}

	// 分页
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := req.PageNum
	pageSize := req.PageSize
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}

	return list, total, err

}

func (l Log) BatchDeleteOperationLogByIds(ids []uint) error {
	err := global.DB.Where("id IN (?)", ids).Unscoped().Delete(&system.OperationLog{}).Error
	return err
}

// SaveOperationLogChannel var Logs []system.OperationLog //全局变量多个线程需要加锁，所以每个线程自己维护一个
//处理OperationLogChan将日志记录到数据库
func (l Log) SaveOperationLogChannel(olc <-chan *system.OperationLog) {
	// 只会在线程开启的时候执行一次
	Logs := make([]system.OperationLog, 0)

	// 一直执行--收到olc就会执行
	for log := range olc {
		Logs = append(Logs, *log)
		// 每10条记录到数据库
		if len(Logs) > 5 {
			global.DB.Create(&Logs)
			Logs = make([]system.OperationLog, 0)
		}
	}
}
