/*
 * @Author: LiuYanFeng
 * @Date: 2026-02-24 16:04:00
 * @LastEditors: LiuYanFeng
 * @LastEditTime: 2026-02-24 16:04:00
 * @FilePath: \happytools\backend\services\unitconverter\unitconverter.go
 * @Description: 单位转换服务
 *
 * Copyright (c) 2026 by ${git_name_email}, All Rights Reserved.
 */
package unitconverter

import (
	"errors"
	"fmt"
)

// UnitConverterService 单位转换服务
type UnitConverterService struct{}

// NewUnitConverterService 创建单位转换服务实例
func NewUnitConverterService() *UnitConverterService {
	return &UnitConverterService{}
}

// ========== 字节转换 ==========

// ByteConversionResult 字节转换结果
type ByteConversionResult struct {
	Bytes     float64 `json:"bytes"`
	Kilobytes float64 `json:"kilobytes"`
	Megabytes float64 `json:"megabytes"`
	Gigabytes float64 `json:"gigabytes"`
	Terabytes float64 `json:"terabytes"`
}

// ConvertBytes 字节转换
// unit: "B", "KB", "MB", "GB", "TB"
func (s *UnitConverterService) ConvertBytes(value float64, unit string) (*ByteConversionResult, error) {
	if value < 0 {
		return nil, errors.New("value cannot be negative")
	}

	var bytes float64

	switch unit {
	case "B", "b":
		bytes = value
	case "KB", "kb":
		bytes = value * 1024
	case "MB", "mb":
		bytes = value * 1024 * 1024
	case "GB", "gb":
		bytes = value * 1024 * 1024 * 1024
	case "TB", "tb":
		bytes = value * 1024 * 1024 * 1024 * 1024
	default:
		return nil, fmt.Errorf("unsupported unit: %s", unit)
	}

	return &ByteConversionResult{
		Bytes:     bytes,
		Kilobytes: bytes / 1024,
		Megabytes: bytes / (1024 * 1024),
		Gigabytes: bytes / (1024 * 1024 * 1024),
		Terabytes: bytes / (1024 * 1024 * 1024 * 1024),
	}, nil
}

// ========== 长度转换 ==========

// LengthConversionResult 长度转换结果
type LengthConversionResult struct {
	Millimeters float64 `json:"millimeters"`
	Centimeters float64 `json:"centimeters"`
	Meters      float64 `json:"meters"`
	Kilometers  float64 `json:"kilometers"`
	Inches      float64 `json:"inches"`
	Feet        float64 `json:"feet"`
	Yards       float64 `json:"yards"`
	Miles       float64 `json:"miles"`
}

// ConvertLength 长度转换
// unit: "mm", "cm", "m", "km", "in", "ft", "yd", "mi"
func (s *UnitConverterService) ConvertLength(value float64, unit string) (*LengthConversionResult, error) {
	if value < 0 {
		return nil, errors.New("value cannot be negative")
	}

	var meters float64

	switch unit {
	case "mm":
		meters = value / 1000
	case "cm":
		meters = value / 100
	case "m":
		meters = value
	case "km":
		meters = value * 1000
	case "in":
		meters = value * 0.0254
	case "ft":
		meters = value * 0.3048
	case "yd":
		meters = value * 0.9144
	case "mi":
		meters = value * 1609.344
	default:
		return nil, fmt.Errorf("unsupported unit: %s", unit)
	}

	return &LengthConversionResult{
		Millimeters: meters * 1000,
		Centimeters: meters * 100,
		Meters:      meters,
		Kilometers:  meters / 1000,
		Inches:      meters / 0.0254,
		Feet:        meters / 0.3048,
		Yards:       meters / 0.9144,
		Miles:       meters / 1609.344,
	}, nil
}

// ========== 时间转换 ==========

// TimeConversionResult 时间转换结果
type TimeConversionResult struct {
	Milliseconds float64 `json:"milliseconds"`
	Seconds      float64 `json:"seconds"`
	Minutes      float64 `json:"minutes"`
	Hours        float64 `json:"hours"`
	Days         float64 `json:"days"`
	Weeks        float64 `json:"weeks"`
	Months       float64 `json:"months"`
	Years        float64 `json:"years"`
}

// ConvertTime 时间转换
// unit: "ms", "s", "min", "h", "d", "w", "mon", "y"
func (s *UnitConverterService) ConvertTime(value float64, unit string) (*TimeConversionResult, error) {
	if value < 0 {
		return nil, errors.New("value cannot be negative")
	}

	var seconds float64

	switch unit {
	case "ms":
		seconds = value / 1000
	case "s":
		seconds = value
	case "min":
		seconds = value * 60
	case "h":
		seconds = value * 3600
	case "d":
		seconds = value * 86400
	case "w":
		seconds = value * 604800
	case "mon":
		seconds = value * 2592000 // 30天
	case "y":
		seconds = value * 31536000 // 365天
	default:
		return nil, fmt.Errorf("unsupported unit: %s", unit)
	}

	return &TimeConversionResult{
		Milliseconds: seconds * 1000,
		Seconds:      seconds,
		Minutes:      seconds / 60,
		Hours:        seconds / 3600,
		Days:         seconds / 86400,
		Weeks:        seconds / 604800,
		Months:       seconds / 2592000,
		Years:        seconds / 31536000,
	}, nil
}
