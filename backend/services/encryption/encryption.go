/*
 * @Author: LiuYanFeng
 * @Date: 2026-02-24 16:05:00
 * @LastEditors: LiuYanFeng
 * @LastEditTime: 2026-02-24 16:05:00
 * @FilePath: \happytools\backend\services\encryption\encryption.go
 * @Description: 加密工具服务
 *
 * Copyright (c) 2026 by ${git_name_email}, All Rights Reserved.
 */
package encryption

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

// EncryptionService 加密工具服务
type EncryptionService struct{}

// NewEncryptionService 创建加密工具服务实例
func NewEncryptionService() *EncryptionService {
	return &EncryptionService{}
}

// ========== MD5 编解码 ==========

// MD5Encode MD5 编码
func (s *EncryptionService) MD5Encode(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

// MD5EncodeBytes MD5 编码（字节数组）
func (s *EncryptionService) MD5EncodeBytes(input []byte) string {
	hash := md5.Sum(input)
	return hex.EncodeToString(hash[:])
}

// ========== SHA1 编解码 ==========

// SHA1Encode SHA1 编码
func (s *EncryptionService) SHA1Encode(input string) string {
	hash := sha1.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

// SHA1EncodeBytes SHA1 编码（字节数组）
func (s *EncryptionService) SHA1EncodeBytes(input []byte) string {
	hash := sha1.Sum(input)
	return hex.EncodeToString(hash[:])
}

// ========== SHA256 编解码 ==========

// SHA256Encode SHA256 编码
func (s *EncryptionService) SHA256Encode(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

// SHA256EncodeBytes SHA256 编码（字节数组）
func (s *EncryptionService) SHA256EncodeBytes(input []byte) string {
	hash := sha256.Sum256(input)
	return hex.EncodeToString(hash[:])
}

// ========== SHA512 编解码 ==========

// SHA512Encode SHA512 编码
func (s *EncryptionService) SHA512Encode(input string) string {
	hash := sha512.Sum512([]byte(input))
	return hex.EncodeToString(hash[:])
}

// SHA512EncodeBytes SHA512 编码（字节数组）
func (s *EncryptionService) SHA512EncodeBytes(input []byte) string {
	hash := sha512.Sum512(input)
	return hex.EncodeToString(hash[:])
}

// ========== Base64 编解码 ==========

// Base64Encode Base64 编码
func (s *EncryptionService) Base64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

// Base64Decode Base64 解码
func (s *EncryptionService) Base64Decode(input string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", errors.New("invalid base64 string: " + err.Error())
	}
	return string(decoded), nil
}

// Base64URLEncode Base64 URL 编码
func (s *EncryptionService) Base64URLEncode(input string) string {
	return base64.URLEncoding.EncodeToString([]byte(input))
}

// Base64URLDecode Base64 URL 解码
func (s *EncryptionService) Base64URLDecode(input string) (string, error) {
	decoded, err := base64.URLEncoding.DecodeString(input)
	if err != nil {
		return "", errors.New("invalid base64 URL string: " + err.Error())
	}
	return string(decoded), nil
}

// ========== UTF-8 编解码 ==========

// UTF8Encode UTF-8 编码（将字符串转换为 UTF-8 字节表示）
func (s *EncryptionService) UTF8Encode(input string) []byte {
	return []byte(input)
}

// UTF8Decode UTF-8 解码（将 UTF-8 字节转换为字符串）
func (s *EncryptionService) UTF8Decode(input []byte) (string, error) {
	return string(input), nil
}

// UTF8ToHex UTF-8 字符串转十六进制
func (s *EncryptionService) UTF8ToHex(input string) string {
	return hex.EncodeToString([]byte(input))
}

// HexToUTF8 十六进制转 UTF-8 字符串
func (s *EncryptionService) HexToUTF8(input string) (string, error) {
	decoded, err := hex.DecodeString(input)
	if err != nil {
		return "", errors.New("invalid hex string: " + err.Error())
	}
	return string(decoded), nil
}

// ========== 批量编解码 ==========

// BatchEncodeResult 批量编码结果
type BatchEncodeResult struct {
	MD5    string `json:"md5"`
	SHA1   string `json:"sha1"`
	SHA256 string `json:"sha256"`
	SHA512 string `json:"sha512"`
	Base64 string `json:"base64"`
}

// BatchEncode 批量编码（一次性生成所有哈希值）
func (s *EncryptionService) BatchEncode(input string) *BatchEncodeResult {
	return &BatchEncodeResult{
		MD5:    s.MD5Encode(input),
		SHA1:   s.SHA1Encode(input),
		SHA256: s.SHA256Encode(input),
		SHA512: s.SHA512Encode(input),
		Base64: s.Base64Encode(input),
	}
}
