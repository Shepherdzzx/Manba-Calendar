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
			name:  "query_next_friday_my_things",
			input: "下周五我有什么事",
			want: ParsedCommand{Intent: IntentQueryEvents, Date: "2026-06-05"},
		},
		{
			name:  "query_today_my_arrangements",
			input: "今天我有什么安排",
			want: ParsedCommand{Intent: IntentQueryEvents, Date: "2026-05-29"},
		},
		{
			name:  "query_today_my_matters",
			input: "今天我有什么事情",
			want: ParsedCommand{Intent: IntentQueryEvents, Date: "2026-05-29"},
		},
		{
			name:  "query_today_specific_meetings",
			input: "今天我有什么会议",
			want: ParsedCommand{Intent: IntentQueryEvents, Title: "会议", Date: "2026-05-29"},
		},
		{
			name:  "example5_delete_gym_tonight",
			input: "取消今晚七点的健身",
			want: ParsedCommand{Intent: IntentDeleteEvent, Title: "健身", Date: "2026-05-29", Time: "19:00"},
		},
		{
			name:  "delete_event_time_only_tonight_thing",
			input: "删除今晚七点的事",
			want: ParsedCommand{Intent: IntentDeleteEvent, Date: "2026-05-29", Time: "19:00"},
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
		{
			name:  "create_event_date_only_today",
			input: "今天要写作业",
			want: ParsedCommand{Intent: IntentCreateEvent, Title: "要写作业", Date: "2026-05-29"},
		},
		{
			name:  "delete_event_date_only_today_thing",
			input: "删除今天的事",
			want: ParsedCommand{Intent: IntentDeleteEvent, Date: "2026-05-29"},
		},
		{
			name:  "delete_event_with_colloquial_shandiao_today",
			input: "删掉今天的事",
			want: ParsedCommand{Intent: IntentDeleteEvent, Date: "2026-05-29"},
		},
		{
			name:  "delete_event_with_colloquial_shandiao_tomorrow",
			input: "删掉明天上午的会议",
			want: ParsedCommand{Intent: IntentDeleteEvent, Title: "会议", Date: "2026-05-30"},
		},
		{
			name:  "create_event_play_game_with_time",
			input: "下周五五点要打游戏",
			want: ParsedCommand{Intent: IntentCreateEvent, Title: "要打游戏", Date: "2026-06-05", Time: "05:00"},
		},
		{
			name:  "create_event_play_game_date_only",
			input: "下周五要打游戏",
			want: ParsedCommand{Intent: IntentCreateEvent, Title: "要打游戏", Date: "2026-06-05"},
		},
		{
			name:  "create_event_play_game_date_only_with_wo",
			input: "下周五我要打游戏",
			want: ParsedCommand{Intent: IntentCreateEvent, Title: "我要打游戏", Date: "2026-06-05"},
		},
		{
			name:  "create_event_review_date_only",
			input: "明天要复习",
			want: ParsedCommand{Intent: IntentCreateEvent, Title: "要复习", Date: "2026-05-30"},
		},
		{
			name:  "create_event_class_date_only",
			input: "后天上课",
			want: ParsedCommand{Intent: IntentCreateEvent, Title: "上课", Date: "2026-05-31"},
		},
		{
			name:  "create_event_go_gym_weekend",
			input: "周六去健身",
			want: ParsedCommand{Intent: IntentCreateEvent, Title: "去健身", Date: "2026-05-30"},
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
		{name: "missing_title_for_delete", input: "删除", reason: ErrMissingTitle},
		{name: "date_only_but_unknown_action", input: "明天那个", reason: ErrUnknownIntent},
		{name: "date_only_but_unknown_action_this", input: "明天这个", reason: ErrUnknownIntent},
		{name: "date_only_but_unknown_action_next_friday", input: "下周五那个", reason: ErrUnknownIntent},
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
