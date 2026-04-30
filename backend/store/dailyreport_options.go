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

// GetAllDailyReportTags 获取所有日报中使用过的标签（去重排序）
func GetAllDailyReportTags() ([]string, error) {
	seen := make(map[string]struct{})
	err := DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(dailyReportBucket)
		if b == nil {
			return nil
		}
		return b.ForEach(func(k, v []byte) error {
			var r DailyReport
			if err := json.Unmarshal(v, &r); err != nil {
				return nil
			}
			for _, tag := range r.Tags {
				if tag != "" {
					seen[tag] = struct{}{}
				}
			}
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	tags := make([]string, 0, len(seen))
	for tag := range seen {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags, nil
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

// MonthTagStat store 层的月标签统计
type MonthTagStat struct {
	Tag  string
	Days float64 // 支持小数（多标签按比例拆分）
}

// MonthStat store 层的月统计
type MonthStat struct {
	Month         string
	TagStats      []MonthTagStat
	UntaggedDates []string
	TotalDays     int
}

// DailyTagRatio 某天各标签的工时比例，所有 tag 的 ratio 之和 = 1.0
// key: date (YYYY-MM-DD), value: JSON of map[string]float64
type DailyTagRatio struct {
	Date   string             `json:"date"`
	Ratios map[string]float64 `json:"ratios"` // tag -> ratio (0~1)
}

// SaveTagRatios 保存某天的标签工时比例
func SaveTagRatios(date string, ratios map[string]float64) error {
	data, err := json.Marshal(DailyTagRatio{Date: date, Ratios: ratios})
	if err != nil {
		return err
	}
	return DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(dailyReportTagRatioBucket)
		return b.Put([]byte(date), data)
	})
}

// GetTagRatios 获取某天的标签工时比例，不存在则返回 nil
func GetTagRatios(date string) (map[string]float64, error) {
	var result map[string]float64
	err := DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(dailyReportTagRatioBucket)
		if b == nil {
			return nil
		}
		v := b.Get([]byte(date))
		if v == nil {
			return nil
		}
		var r DailyTagRatio
		if err := json.Unmarshal(v, &r); err != nil {
			return err
		}
		result = r.Ratios
		return nil
	})
	return result, err
}

// getAllTagRatios 内部：获取所有已保存的比例，返回 date -> ratios
func getAllTagRatios(tx *bbolt.Tx) map[string]map[string]float64 {
	out := make(map[string]map[string]float64)
	b := tx.Bucket(dailyReportTagRatioBucket)
	if b == nil {
		return out
	}
	b.ForEach(func(k, v []byte) error {
		var r DailyTagRatio
		if json.Unmarshal(v, &r) == nil {
			out[r.Date] = r.Ratios
		}
		return nil
	})
	return out
}

// GetMonthlyTagStats 按月统计每个标签工时（天数，支持小数比例）及未打标签日期
func GetMonthlyTagStats() ([]MonthStat, error) {
	type monthData struct {
		tagDays       map[string]float64
		untaggedDates []string
		totalDays     int
	}

	monthMap := make(map[string]*monthData)
	var allRatios map[string]map[string]float64

	err := DB.View(func(tx *bbolt.Tx) error {
		allRatios = getAllTagRatios(tx)

		b := tx.Bucket(dailyReportBucket)
		if b == nil {
			return nil
		}
		return b.ForEach(func(k, v []byte) error {
			var r DailyReport
			if err := json.Unmarshal(v, &r); err != nil {
				return nil
			}
			if len(r.Date) < 7 {
				return nil
			}
			month := r.Date[:7]
			if _, ok := monthMap[month]; !ok {
				monthMap[month] = &monthData{
					tagDays: make(map[string]float64),
				}
			}
			md := monthMap[month]
			md.totalDays++

			validTags := make([]string, 0, len(r.Tags))
			for _, t := range r.Tags {
				if t != "" {
					validTags = append(validTags, t)
				}
			}

			if len(validTags) == 0 {
				md.untaggedDates = append(md.untaggedDates, r.Date)
				return nil
			}

			// 取用户自定义比例，若无则均分
			ratios, hasCustom := allRatios[r.Date]
			for _, tag := range validTags {
				var ratio float64
				if hasCustom {
					if rv, ok := ratios[tag]; ok {
						ratio = rv
					} else {
						ratio = 1.0 / float64(len(validTags))
					}
				} else {
					ratio = 1.0 / float64(len(validTags))
				}
				md.tagDays[tag] += ratio
			}
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	months := make([]string, 0, len(monthMap))
	for m := range monthMap {
		months = append(months, m)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(months)))

	result := make([]MonthStat, 0, len(months))
	for _, month := range months {
		md := monthMap[month]

		tagStats := make([]MonthTagStat, 0, len(md.tagDays))
		for tag, days := range md.tagDays {
			tagStats = append(tagStats, MonthTagStat{Tag: tag, Days: days})
		}
		sort.Slice(tagStats, func(i, j int) bool {
			if tagStats[i].Days != tagStats[j].Days {
				return tagStats[i].Days > tagStats[j].Days
			}
			return tagStats[i].Tag < tagStats[j].Tag
		})

		sort.Strings(md.untaggedDates)

		result = append(result, MonthStat{
			Month:         month,
			TagStats:      tagStats,
			UntaggedDates: md.untaggedDates,
			TotalDays:     md.totalDays,
		})
	}

	return result, nil
}
