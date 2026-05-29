package parser

import (
	"testing"
	"time"
)

func TestParse_Examples(t *testing.T) {
	now := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)
	parser := New(now)
	trueVal := true

	tests := []struct {
		name  string
		input string
		want  ParsedCommand
	}{
		{
			name:  "example1_create_meeting",
			input: "明天下午三点和张三开会",
			want: ParsedCommand{Intent: IntentCreateEvent, Title: "和张三开会", Date: "2026-05-30", Time: "15:00"},
		},
		{
			name:  "example2_create_shopping",
			input: "后天上午十点去超市买东西",
			want: ParsedCommand{Intent: IntentCreateEvent, Title: "去超市买东西", Date: "2026-05-31", Time: "10:00"},
		},
		{
			name:  "example3_query_today",
			input: "帮我看看今天还有什么安排",
			want: ParsedCommand{Intent: IntentQueryEvents, Date: "2026-05-29"},
		},
		{
			name:  "example4_query_meetings_next_monday",
			input: "下周一有什么会议",
			want: ParsedCommand{Intent: IntentQueryEvents, Title: "会议", Date: "2026-06-01"},
		},
		{
			name:  "example5_delete_gym_tonight",
			input: "取消今晚七点的健身",
			want: ParsedCommand{Intent: IntentDeleteEvent, Title: "健身", Date: "2026-05-29", Time: "19:00"},
		},
		{
			name:  "example6_delete_interview_tomorrow_morning",
			input: "删除明天上午的面试",
			want: ParsedCommand{Intent: IntentDeleteEvent, Title: "面试", Date: "2026-05-30"},
		},
		{
			name:  "example7_create_meeting_next_wednesday",
			input: "下周三下午两点半开会",
			want: ParsedCommand{Intent: IntentCreateEvent, Title: "开会", Date: "2026-06-03", Time: "14:30"},
		},
		{
			name:  "example8_create_movie_friday_evening",
			input: "星期五晚上八点看电影",
			want: ParsedCommand{Intent: IntentCreateEvent, Title: "看电影", Date: "2026-06-05", Time: "20:00"},
		},
		{
			name:  "create_event_with_reminder",
			input: "明天下午三点提醒我和张三开会",
			want: ParsedCommand{Intent: IntentCreateEvent, Title: "和张三开会", Date: "2026-05-30", Time: "15:00", NeedReminder: &trueVal},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.Parse(tt.input)
			if result.Err != nil {
				t.Fatalf("unexpected error: %v", result.Err)
			}
			if result.Command == nil {
				t.Fatal("expected command, got nil")
			}
			got := *result.Command
			if got.Intent != tt.want.Intent || got.Title != tt.want.Title || got.Date != tt.want.Date || got.Time != tt.want.Time {
				t.Fatalf("unexpected command: got %+v want %+v", got, tt.want)
			}
			if (got.NeedReminder == nil) != (tt.want.NeedReminder == nil) {
				t.Fatalf("unexpected need_reminder presence: got %+v want %+v", got.NeedReminder, tt.want.NeedReminder)
			}
		})
		}
}

func TestParse_Errors(t *testing.T) {
	now := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)
	parser := New(now)

	tests := []struct {
		name   string
		input  string
		reason string
	}{
		{name: "empty_input", input: "   ", reason: ErrEmptyInput},
		{name: "unknown_intent", input: "你好世界", reason: ErrUnknownIntent},
		{name: "missing_time", input: "添加明天和张三开会", reason: ErrMissingTime},
		{name: "missing_title_for_delete", input: "删除明天上午的", reason: ErrMissingTitle},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.Parse(tt.input)
			if result.Err == nil {
				t.Fatal("expected error, got nil")
			}
			if result.Err.Reason != tt.reason {
				t.Fatalf("unexpected error reason: got %s want %s", result.Err.Reason, tt.reason)
			}
		})
	}
}
