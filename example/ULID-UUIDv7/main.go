package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math/big"
	"strings"
	"time"
)

// =============================
//  ULID
// =============================

// Base32 编码字符表（Crockford Base32）
const base32Chars = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

// 生成 ULID
func GenerateULID() string {
	now := time.Now().UnixMilli()

	// 时间戳部分：48位，编码成10个Base32字符
	timeChars := make([]byte, 10)
	t := now
	for i := 9; i >= 0; i-- {
		timeChars[i] = base32Chars[t&0x1F]
		t >>= 5
	}

	// 随机数部分：80位，编码成16个Base32字符
	randBytes := make([]byte, 10)
	rand.Read(randBytes)
	randNum := new(big.Int).SetBytes(randBytes)

	randChars := make([]byte, 16)
	mask := big.NewInt(0x1F)
	tmp := new(big.Int)
	for i := 15; i >= 0; i-- {
		tmp.And(randNum, mask)
		randChars[i] = base32Chars[tmp.Int64()]
		randNum.Rsh(randNum, 5)
	}

	return string(timeChars) + string(randChars)
}

// 逆运算 ULID → 时间
func ParseULIDTime(ulid string) time.Time {
	ulid = strings.ToUpper(ulid)
	timePart := ulid[:10]

	var ms int64
	for _, c := range timePart {
		ms <<= 5
		ms |= int64(strings.IndexRune(base32Chars, c))
	}

	return time.UnixMilli(ms)
}

// =============================
//  UUID v7
// =============================

// 生成 UUID v7
func GenerateUUIDv7() string {
	now := time.Now().UnixMilli()

	b := make([]byte, 16)
	rand.Read(b)

	// 前48位写入时间戳
	binary.BigEndian.PutUint32(b[0:4], uint32(now>>16))
	binary.BigEndian.PutUint16(b[4:6], uint16(now&0xFFFF))

	// 第7字节高4位设为版本号 7
	b[6] = (b[6] & 0x0F) | 0x70

	// 第9字节高2位设为变体标识 10
	b[8] = (b[8] & 0x3F) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// 逆运算 UUID v7 → 时间
func ParseUUIDv7Time(uuid string) time.Time {
	// 去掉横线，取前12位十六进制（48位时间戳）
	clean := strings.ReplaceAll(uuid, "-", "")
	timePart := clean[:12]

	var ms int64
	fmt.Sscanf(timePart, "%x", &ms)

	return time.UnixMilli(ms)
}

// =============================
//  主函数演示
// =============================

func main() {
	fmt.Println("========== ULID ==========")
	ulid := GenerateULID()
	fmt.Println("生成的 ULID   :", ulid)
	fmt.Println("逆运算时间    :", ParseULIDTime(ulid).Format("2006-01-02 15:04:05.000"))

	fmt.Println()
	fmt.Println("========== UUID v7 ==========")
	uuid := GenerateUUIDv7()
	fmt.Println("生成的 UUID v7:", uuid)
	fmt.Println("逆运算时间    :", ParseUUIDv7Time(uuid).Format("2006-01-02 15:04:05.000"))

	fmt.Println()
	fmt.Println("========== 批量生成对比 ==========")
	fmt.Printf("%-30s  %s\n", "ULID", "UUID v7")
	fmt.Println(strings.Repeat("-", 70))
	for i := 0; i < 5; i++ {
		fmt.Printf("%-30s  %s\n", GenerateULID(), GenerateUUIDv7())
	}
}