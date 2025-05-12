package system

import (
	"time"
)

type UserLoginHistory struct {
	ID       uint64    `gorm:"primarykey"`                // 主键类型修改为 uint64
	UserID   uint64    `gorm:"not null"`                  // 修改为 uint64，确保和数据库的 bigint unsigned 一致！！
	DateTime time.Time `gorm:"column:date_time;not null"` //没有设置默认值，数据库报错
}

func (UserLoginHistory) TableName() string {
	return "sys_user_login_history" //绑定表格
}

// type UserLoginTimeRange struct {
//
// }

// type GetUserIdsByLoginTimeRequest struct {
// 	StartDate string `json:"start_time" binding:"required"` //带required是必填项，not null
// 	EndDate   string `json:"end_time" binding:"required"`
// }
