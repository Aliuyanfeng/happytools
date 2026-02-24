/*
 * @Author: LiuYanFeng
 * @Date: 2026-02-06 10:44:35
 * @LastEditors: LiuYanFeng
 * @LastEditTime: 2026-02-06 11:02:50
 * @FilePath: \happytools\backend\services\appsettings\appsettings.go
 * @Description: 像珍惜礼物一样珍惜今天
 *
 * Copyright (c) 2026 by ${git_name_email}, All Rights Reserved.
 */
// Package appsettings
// @author: liuyanfeng
// @date: 2026/2/6
// @description: 应用设置服务
package appsettings

import (
	"github.com/Aliuyanfeng/happytools/backend/store"
)

type AppSettingsService struct{}

// NewAppSettingsService 创建应用设置服务实例
func NewAppSettingsService() *AppSettingsService {
	return &AppSettingsService{}
}

// GetLastUsedTime 获取上次使用时间
func (s *AppSettingsService) GetLastUsedTime() string {
	// 获取上次使用时间
	lastUsed, err := store.GetLastUsedTime()
	if err != nil || lastUsed == nil {
		return ""
	}

	// 格式化为可读的时间字符串
	formatted := lastUsed.Format("2006-01-02 15:04:05")
	return formatted
}

// UpdateLastUsedTime 更新上次使用时间
func (s *AppSettingsService) UpdateLastUsedTime() error {
	return store.UpdateLastUsedTime()
}

// GetLastUsedTimestamp 获取上次使用时间戳（秒）
func (s *AppSettingsService) GetLastUsedTimestamp() *int64 {
	lastUsed, err := store.GetLastUsedTime()
	if err != nil || lastUsed == nil {
		return nil
	}

	timestamp := lastUsed.Unix()
	return &timestamp
}
