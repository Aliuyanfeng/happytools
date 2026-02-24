// Package store
// @author: liuyanfeng
// @date: 2026/2/3
// @description: DailyReport 日报增删改查操作
package store

import (
	"encoding/json"
	"errors"
	"sort"
	"time"

	"go.etcd.io/bbolt"
)

// DailyReport 日报
type DailyReport struct {
	ID        int       `json:"id"`
	Date      string    `json:"date"` // YYYY-MM-DD
	Content   string    `json:"content"`
	Summary   string    `json:"summary"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DailyReportStats 日报统计
type DailyReportStats struct {
	TotalReports  int `json:"total_reports"`
	CurrentStreak int `json:"current_streak"`
	LongestStreak int `json:"longest_streak"`
}

// GetDailyReport 获取指定日期的日报
func GetDailyReport(date string) (*DailyReport, error) {
	var report *DailyReport
	err := DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(dailyReportBucket)
		if b == nil {
			return nil
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var r DailyReport
			if err := json.Unmarshal(v, &r); err != nil {
				continue
			}
			if r.Date == date {
				report = &r
				return nil
			}
		}
		return nil
	})
	return report, err
}

// GetDailyReports 获取日期范围内的日报
func GetDailyReports(startDate, endDate string) ([]DailyReport, error) {
	reports := make([]DailyReport, 0)
	err := DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(dailyReportBucket)
		if b == nil {
			return nil
		}
		return b.ForEach(func(k, v []byte) error {
			var r DailyReport
			if err := json.Unmarshal(v, &r); err != nil {
				return err
			}
			if r.Date >= startDate && r.Date <= endDate {
				reports = append(reports, r)
			}
			return nil
		})
	})

	// 按日期倒序排列
	sort.Slice(reports, func(i, j int) bool {
		return reports[i].Date > reports[j].Date
	})

	return reports, err
}

// GetAllDailyReports 获取所有日报 (分页)
func GetAllDailyReports(limit, offset int) ([]DailyReport, int, error) {
	allReports := make([]DailyReport, 0)

	err := DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(dailyReportBucket)
		if b == nil {
			return nil
		}
		return b.ForEach(func(k, v []byte) error {
			var r DailyReport
			if err := json.Unmarshal(v, &r); err != nil {
				return err
			}
			allReports = append(allReports, r)
			return nil
		})
	})
	if err != nil {
		return nil, 0, err
	}

	// 按日期倒序排列
	sort.Slice(allReports, func(i, j int) bool {
		return allReports[i].Date > allReports[j].Date
	})

	total := len(allReports)

	// 分页
	if offset >= total {
		return []DailyReport{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}

	return allReports[offset:end], total, nil
}

// CreateOrUpdateDailyReport 创建或更新日报
func CreateOrUpdateDailyReport(date, content, summary string, tags []string) (*DailyReport, error) {
	if date == "" || content == "" {
		return nil, errors.New("date and content cannot be empty")
	}

	var report DailyReport
	err := DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(dailyReportBucket)

		// 检查是否已存在该日期的日报
		var existingID *int
		b.ForEach(func(k, v []byte) error {
			var r DailyReport
			if json.Unmarshal(v, &r) == nil && r.Date == date {
				existingID = &r.ID
				return errors.New("break")
			}
			return nil
		})

		now := time.Now()

		if existingID != nil {
			// 更新现有日报
			v := b.Get(itob(*existingID))
			if err := json.Unmarshal(v, &report); err != nil {
				return err
			}
			report.Content = content
			report.Summary = summary
			report.Tags = tags
			report.UpdatedAt = now
		} else {
			// 创建新日报
			id, _ := b.NextSequence()
			report = DailyReport{
				ID:        int(id),
				Date:      date,
				Content:   content,
				Summary:   summary,
				Tags:      tags,
				CreatedAt: now,
				UpdatedAt: now,
			}
		}

		data, err := json.Marshal(report)
		if err != nil {
			return err
		}
		return b.Put(itob(report.ID), data)
	})
	if err != nil {
		return nil, err
	}
	return &report, nil
}

// DeleteDailyReport 删除日报
func DeleteDailyReport(id int) error {
	return DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(dailyReportBucket)
		return b.Delete(itob(id))
	})
}

// GetDailyReportStats 获取日报统计
func GetDailyReportStats() (*DailyReportStats, error) {
	stats := &DailyReportStats{}
	dates := make([]string, 0)

	err := DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(dailyReportBucket)
		if b == nil {
			return nil
		}
		return b.ForEach(func(k, v []byte) error {
			var r DailyReport
			if err := json.Unmarshal(v, &r); err != nil {
				return err
			}
			dates = append(dates, r.Date)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	stats.TotalReports = len(dates)
	if stats.TotalReports == 0 {
		return stats, nil
	}

	// 排序日期
	sort.Strings(dates)

	// 计算连续打卡
	currentStreak := 0
	longestStreak := 0
	tempStreak := 1

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	// 计算当前连续天数
	if dates[len(dates)-1] == today || dates[len(dates)-1] == yesterday {
		currentStreak = 1
		for i := len(dates) - 2; i >= 0; i-- {
			prevDate, _ := time.Parse("2006-01-02", dates[i])
			currDate, _ := time.Parse("2006-01-02", dates[i+1])

			if currDate.Sub(prevDate).Hours() == 24 {
				currentStreak++
			} else {
				break
			}
		}
	}

	// 计算最长连续天数
	for i := 1; i < len(dates); i++ {
		prevDate, _ := time.Parse("2006-01-02", dates[i-1])
		currDate, _ := time.Parse("2006-01-02", dates[i])

		if currDate.Sub(prevDate).Hours() == 24 {
			tempStreak++
		} else {
			if tempStreak > longestStreak {
				longestStreak = tempStreak
			}
			tempStreak = 1
		}
	}
	if tempStreak > longestStreak {
		longestStreak = tempStreak
	}

	stats.CurrentStreak = currentStreak
	stats.LongestStreak = longestStreak

	return stats, nil
}
