package parser

import "strings"

var weakDeleteTitles = map[string]struct{}{
	"事":   {},
	"我事":  {},
	"事情":  {},
	"安排":  {},
	"东西":  {},
	"那个":  {},
	"这个":  {},
}

func ExtractTitle(text string, intent Intent, dt DateTimeParseResult) string {
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
		"的事", "",
		"事情", "",
		"的", "",
	)
	if strings.Contains(remaining, "提醒我") {
		remaining = strings.ReplaceAll(remaining, "提醒我", "")
	}
	remaining = replacer.Replace(remaining)
	if intent == IntentQueryEvents {
		remaining = strings.TrimPrefix(remaining, "我")
		remaining = strings.TrimPrefix(remaining, "你")
		remaining = strings.TrimPrefix(remaining, "他")
		remaining = strings.TrimPrefix(remaining, "她")
	}
	remaining = strings.TrimSpace(remaining)
	remaining = strings.Trim(remaining, " ，。！？,.!?")
	if _, ok := weakDeleteTitles[remaining]; ok {
		return ""
	}
	return remaining
}
