# JSON 协议规范（v1）

## 目标

本文定义 Manba Alert 第一版中，本地解析器输出给应用侧的结构化 JSON 指令格式。

目标是让解析结果稳定、应用可校验、关键操作可确认。当前 v1 只覆盖三个高频动作：

- `create_event`
- `query_events`
- `delete_event`

基础示例：

```json
{
  "intent": "create_event",
  "title": "开会",
  "date": "2026-05-30",
  "time": "15:00",
  "need_reminder": true
}
```

## 总体约束

- 解析器输出必须是 **单个 JSON 对象**
- 不允许输出解释、注释、Markdown 代码块或多余文本
- 应用侧**不能直接信任解析结果**，必须先校验后执行
- 对删除类操作，应用侧应增加确认步骤，不能仅凭解析结果直接执行

## 字段定义

| 字段名 | 类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| `intent` | `string` | 是 | 操作类型，必须是允许列表中的一个 |
| `title` | `string` | 视 intent 而定 | 事件标题或待匹配的事件名称 |
| `date` | `string` | 视 intent 而定 | 日期，格式为 `YYYY-MM-DD` |
| `time` | `string` | 视 intent 而定 | 时间，格式为 `HH:MM`，24 小时制 |
| `need_reminder` | `boolean` | 否 | 是否需要提醒，主要用于创建事件 |

## intent 取值

### 1. `create_event`

用于新增日程。

必填字段：
- `intent`
- `title`
- `date`

可选字段：
- `time`
- `need_reminder`

示例：

```json
{
  "intent": "create_event",
  "title": "和张三开会",
  "date": "2026-05-30",
  "time": "15:00",
  "need_reminder": true
}
```

仅日期示例：

```json
{
  "intent": "create_event",
  "title": "要写作业",
  "date": "2026-05-29"
}
```

### 2. `query_events`

用于查询日程。

必填字段：
- `intent`

可选字段：
- `title`
- `date`
- `time`

说明：
- 如果用户只说“今天有什么安排”，可只提供 `date`
- 如果用户说“下周一有什么会议”，可以同时提供 `date` 和一个较宽泛的 `title`
- `query_events` 不使用 `need_reminder`

示例：

```json
{
  "intent": "query_events",
  "date": "2026-05-29"
}
```

### 3. `delete_event`

用于删除或取消日程。

必填字段：
- `intent`
- `title`

可选字段：
- `date`
- `time`

说明：
- 删除动作应尽量携带更多定位信息，以减少误删
- 当用户表达里包含日期或时间时，解析器应尽量补全
- 应用侧收到该指令后，仍需要向用户确认

示例：

```json
{
  "intent": "delete_event",
  "title": "健身",
  "date": "2026-05-29",
  "time": "19:00"
}
```

## 各 intent 的字段要求矩阵

| intent | `title` | `date` | `time` | `need_reminder` |
| --- | --- | --- | --- | --- |
| `create_event` | 必填 | 必填 | 可选 | 可选 |
| `query_events` | 可选 | 可选 | 可选 | 忽略 |
| `delete_event` | 必填 | 可选 | 可选 | 忽略 |

## 日期与时间格式

### 日期
- 格式：`YYYY-MM-DD`
- 示例：`2026-05-30`

### 时间
- 格式：`HH:MM`
- 使用 24 小时制
- 示例：`09:00`、`15:30`、`20:00`

## 应用侧校验规则

应用在执行任何指令前，至少应校验：

1. 返回值是否为合法 JSON
2. `intent` 是否属于以下允许范围：
   - `create_event`
   - `query_events`
   - `delete_event`
3. `date` 是否满足 `YYYY-MM-DD` 格式，且是合法日期
4. `time` 是否满足 `HH:MM` 格式，且小时在 `00-23`、分钟在 `00-59`
5. `create_event` 的 `title`、`date` 是否完整
6. `delete_event` 的 `title` 是否非空
7. 未知字段是否需要忽略或记录日志（建议忽略执行，但保留调试信息）

## v1 范围边界

当前协议 **不包含**：

- 修改事件
- 重复事件
- 持续时长
- 地点、参与人等扩展字段
- 复杂相对时间（如“下下周二傍晚”）

这些能力可以在后续版本中扩展，但不应进入当前第一版协议。 
