package parser

import "strings"

var (
	deleteKeywords = []string{"删除", "取消", "移除", "去掉", "删掉", "删了", "删除掉"}
	createKeywords = []string{"添加", "新增", "创建", "安排", "记下", "提醒我"}
	queryKeywords  = []string{"查看", "查询", "看看", "有什么安排", "有哪些安排", "有什么", "还有什么安排", "还有什么"}
)

func DetectIntent(text string) Intent {
	normalized := strings.TrimSpace(text)
	if normalized == "" {
		return ""
	}

	for _, kw := range deleteKeywords {
		if strings.Contains(normalized, kw) {
			return IntentDeleteEvent
		}
	}
	for _, kw := range queryKeywords {
		if strings.Contains(normalized, kw) {
			return IntentQueryEvents
		}
	}
	for _, kw := range createKeywords {
		if strings.Contains(normalized, kw) {
			return IntentCreateEvent
		}
	}

	return ""
}
