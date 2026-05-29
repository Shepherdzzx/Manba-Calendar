package parser

import "strings"

func ExtractTitle(text string, dt DateTimeParseResult) string {
	remaining := dt.Remaining
	for _, kw := range deleteKeywords {
		remaining = strings.ReplaceAll(remaining, kw, "")
	}
	for _, kw := range createKeywords {
		remaining = strings.ReplaceAll(remaining, kw, "")
	}
	for _, kw := range queryKeywords {
		remaining = strings.ReplaceAll(remaining, kw, "")
	}

	replacer := strings.NewReplacer(
		"帮我", "",
		"还有什么", "",
		"还有", "",
		"还", "",
		"什么", "",
		"有什么", "",
		"有哪些", "",
		"安排", "",
		"有", "",
		"的", "",
	)
	if strings.Contains(remaining, "提醒我") {
		remaining = strings.ReplaceAll(remaining, "提醒我", "")
	}
	remaining = replacer.Replace(remaining)
	remaining = strings.TrimSpace(remaining)
	remaining = strings.Trim(remaining, " ，。！？,.!?")
	return remaining
}
