package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type DateTimeParser struct {
	now time.Time
}

type DateTimeParseResult struct {
	Date      string
	Time      string
	Remaining string
}

func NewDateTimeParser(now time.Time) DateTimeParser {
	return DateTimeParser{now: now}
}

var (
	reNextWeekday = regexp.MustCompile(`下(?:周|星期)([一二三四五六日天])`)
	reWeekday     = regexp.MustCompile(`(?:周|星期)([一二三四五六日天])`)
	reHalfHour    = regexp.MustCompile(`([零一二两三四五六七八九十\d]+)点半`)
	reHour        = regexp.MustCompile(`([零一二两三四五六七八九十\d]+)点`)
)

func (p DateTimeParser) Parse(text string) DateTimeParseResult {
	remaining := text
	date := ""
	tm := ""

	switch {
	case strings.Contains(remaining, "后天"):
		date = p.now.AddDate(0, 0, 2).Format("2006-01-02")
		remaining = strings.ReplaceAll(remaining, "后天", "")
	case strings.Contains(remaining, "明天"):
		date = p.now.AddDate(0, 0, 1).Format("2006-01-02")
		remaining = strings.ReplaceAll(remaining, "明天", "")
	case strings.Contains(remaining, "今天"):
		date = p.now.Format("2006-01-02")
		remaining = strings.ReplaceAll(remaining, "今天", "")
	case strings.Contains(remaining, "今晚"):
		date = p.now.Format("2006-01-02")
		remaining = strings.ReplaceAll(remaining, "今晚", "晚上")
	default:
		if m := reNextWeekday.FindStringSubmatch(remaining); len(m) == 2 {
			date = resolveWeekday(p.now, m[1], true)
			remaining = reNextWeekday.ReplaceAllString(remaining, "")
		} else if m := reWeekday.FindStringSubmatch(remaining); len(m) == 2 {
			date = resolveWeekday(p.now, m[1], false)
			remaining = reWeekday.ReplaceAllString(remaining, "")
		}
	}

	period := detectPeriod(remaining)
	remaining = stripPeriods(remaining)

	if m := reHalfHour.FindStringSubmatch(remaining); len(m) == 2 {
		if hour, ok := parseChineseHour(m[1]); ok {
			tm = formatHour(hour, 30, period)
		}
		remaining = reHalfHour.ReplaceAllString(remaining, "")
	} else if m := reHour.FindStringSubmatch(remaining); len(m) == 2 {
		if hour, ok := parseChineseHour(m[1]); ok {
			tm = formatHour(hour, 0, period)
		}
		remaining = reHour.ReplaceAllString(remaining, "")
	}

	remaining = strings.TrimSpace(remaining)
	return DateTimeParseResult{Date: date, Time: tm, Remaining: remaining}
}

func detectPeriod(text string) string {
	switch {
	case strings.Contains(text, "上午") || strings.Contains(text, "早上"):
		return "morning"
	case strings.Contains(text, "下午"):
		return "afternoon"
	case strings.Contains(text, "晚上"):
		return "evening"
	default:
		return ""
	}
}

func stripPeriods(text string) string {
	replacer := strings.NewReplacer("上午", "", "早上", "", "下午", "", "晚上", "")
	return replacer.Replace(text)
}

func resolveWeekday(now time.Time, zh string, nextWeek bool) string {
	target, ok := weekdayMap[zh]
	if !ok {
		return ""
	}
	current := int(now.Weekday())
	if current == 0 {
		current = 7
	}
	if nextWeek {
		mondayCurrentWeek := now.AddDate(0, 0, -(current - 1))
		mondayNextWeek := mondayCurrentWeek.AddDate(0, 0, 7)
		return mondayNextWeek.AddDate(0, 0, target-1).Format("2006-01-02")
	}
	delta := target - current
	if delta <= 0 {
		delta += 7
	}
	return now.AddDate(0, 0, delta).Format("2006-01-02")
}

var weekdayMap = map[string]int{
	"一": 1,
	"二": 2,
	"三": 3,
	"四": 4,
	"五": 5,
	"六": 6,
	"日": 7,
	"天": 7,
}

func parseChineseHour(raw string) (int, bool) {
	if raw == "" {
		return 0, false
	}
	if n, err := strconv.Atoi(raw); err == nil {
		return n, true
	}
	if raw == "十" {
		return 10, true
	}
	if strings.HasPrefix(raw, "十") {
		ones, ok := chineseDigit(strings.TrimPrefix(raw, "十"))
		if !ok {
			return 0, false
		}
		return 10 + ones, true
	}
	if strings.HasSuffix(raw, "十") {
		tens, ok := chineseDigit(strings.TrimSuffix(raw, "十"))
		if !ok {
			return 0, false
		}
		return tens * 10, true
	}
	parts := strings.Split(raw, "十")
	if len(parts) == 2 {
		tens, ok1 := chineseDigit(parts[0])
			ones, ok2 := chineseDigit(parts[1])
		if ok1 && ok2 {
			return tens*10 + ones, true
		}
	}
	return chineseDigit(raw)
}

func chineseDigit(raw string) (int, bool) {
	switch raw {
	case "零":
		return 0, true
	case "一":
		return 1, true
	case "二", "两":
		return 2, true
	case "三":
		return 3, true
	case "四":
		return 4, true
	case "五":
		return 5, true
	case "六":
		return 6, true
	case "七":
		return 7, true
	case "八":
		return 8, true
	case "九":
		return 9, true
	default:
		return 0, false
	}
}

func formatHour(hour, minute int, period string) string {
	switch period {
	case "afternoon", "evening":
		if hour < 12 {
			hour += 12
		}
	case "morning":
		if hour == 12 {
			hour = 0
		}
	}
	return fmt.Sprintf("%02d:%02d", hour, minute)
}
