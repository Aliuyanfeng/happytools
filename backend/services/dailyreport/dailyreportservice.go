// Package dailyreport
// @author: liuyanfeng
// @date: 2026/2/3
// @description: 日报服务
package dailyreport

import (
	"errors"

	"github.com/Aliuyanfeng/happytools/backend/store"
)

type DailyReport struct {
	ID      int      `json:"id"`
	Date    string   `json:"date"`
	Content string   `json:"content"`
	Summary string   `json:"summary"`
	Tags    []string `json:"tags"`
}

type DailyReportStats struct {
	TotalReports  int `json:"total_reports"`
	CurrentStreak int `json:"current_streak"`
	LongestStreak int `json:"longest_streak"`
}

// MonthTagStat 某月某标签的统计
type MonthTagStat struct {
	Tag  string  `json:"tag"`
	Days float64 `json:"days"` // 支持小数（多标签按比例拆分）
}

// MonthStat 某月的统计数据
type MonthStat struct {
	Month         string         `json:"month"`          // YYYY-MM
	TagStats      []MonthTagStat `json:"tag_stats"`
	UntaggedDates []string       `json:"untagged_dates"` // 未打标签的日期列表
	TotalDays     int            `json:"total_days"`
}

type DailyReportService struct{}

func NewDailyReportService() *DailyReportService {
	return &DailyReportService{}
}

// Get 获取指定日期的日报
func (d *DailyReportService) Get(date string) (*DailyReport, error) {
	report, err := store.GetDailyReport(date)
	if err != nil {
		return nil, err
	}
	if report == nil {
		return nil, nil
	}

	return &DailyReport{
		ID:      report.ID,
		Date:    report.Date,
		Content: report.Content,
		Summary: report.Summary,
		Tags:    report.Tags,
	}, nil
}

// GetRange 获取日期范围内的日报
func (d *DailyReportService) GetRange(startDate, endDate string) ([]DailyReport, error) {
	reports, err := store.GetDailyReports(startDate, endDate)
	if err != nil {
		return nil, err
	}

	result := make([]DailyReport, len(reports))
	for i, r := range reports {
		result[i] = DailyReport{
			ID:      r.ID,
			Date:    r.Date,
			Content: r.Content,
			Summary: r.Summary,
			Tags:    r.Tags,
		}
	}
	return result, nil
}

// GetAll 获取所有日报 (分页)
func (d *DailyReportService) GetAll(limit, offset int) ([]DailyReport, int, error) {
	reports, total, err := store.GetAllDailyReports(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	result := make([]DailyReport, len(reports))
	for i, r := range reports {
		result[i] = DailyReport{
			ID:      r.ID,
			Date:    r.Date,
			Content: r.Content,
			Summary: r.Summary,
			Tags:    r.Tags,
		}
	}
	return result, total, nil
}

// Save 保存日报 (创建或更新)
func (d *DailyReportService) Save(date, content, summary string, tags []string) (*DailyReport, error) {
	if date == "" || content == "" {
		return nil, errors.New("date and content cannot be empty")
	}

	report, err := store.CreateOrUpdateDailyReport(date, content, summary, tags)
	if err != nil {
		return nil, err
	}

	return &DailyReport{
		ID:      report.ID,
		Date:    report.Date,
		Content: report.Content,
		Summary: report.Summary,
		Tags:    report.Tags,
	}, nil
}

// Delete 删除日报
func (d *DailyReportService) Delete(id int) error {
	return store.DeleteDailyReport(id)
}

// GetAllTags 获取所有日报中使用过的标签（去重）
func (d *DailyReportService) GetAllTags() ([]string, error) {
	return store.GetAllDailyReportTags()
}

// GetStats 获取日报统计
func (d *DailyReportService) GetStats() (*DailyReportStats, error) {
	stats, err := store.GetDailyReportStats()
	if err != nil {
		return nil, err
	}

	return &DailyReportStats{
		TotalReports:  stats.TotalReports,
		CurrentStreak: stats.CurrentStreak,
		LongestStreak: stats.LongestStreak,
	}, nil
}

// GetMonthlyTagStats 获取所有月份的标签工时统计
func (d *DailyReportService) GetMonthlyTagStats() ([]MonthStat, error) {
	monthStats, err := store.GetMonthlyTagStats()
	if err != nil {
		return nil, err
	}

	result := make([]MonthStat, len(monthStats))
	for i, ms := range monthStats {
		tagStats := make([]MonthTagStat, len(ms.TagStats))
		for j, ts := range ms.TagStats {
			tagStats[j] = MonthTagStat{Tag: ts.Tag, Days: ts.Days}
		}
		result[i] = MonthStat{
			Month:         ms.Month,
			TagStats:      tagStats,
			UntaggedDates: ms.UntaggedDates,
			TotalDays:     ms.TotalDays,
		}
	}
	return result, nil
}

// SaveTagRatios 保存某天各标签的工时比例（所有比例之和应为 1.0）
func (d *DailyReportService) SaveTagRatios(date string, ratios map[string]float64) error {
	return store.SaveTagRatios(date, ratios)
}

// GetTagRatios 获取某天各标签的工时比例，不存在则返回空 map
func (d *DailyReportService) GetTagRatios(date string) (map[string]float64, error) {
	ratios, err := store.GetTagRatios(date)
	if err != nil {
		return nil, err
	}
	if ratios == nil {
		return map[string]float64{}, nil
	}
	return ratios, nil
}
