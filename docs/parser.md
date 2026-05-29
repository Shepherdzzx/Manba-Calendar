# Android 本地解析器设计说明（v1）

## 目标

本文描述 Manba Alert 第一版 Android 端本地解析器的设计方案。

解析器的职责是把 ASR 输出的中文文本直接转换为 `docs/protocol.md` 定义的结构化 JSON 指令，并在应用执行前提供统一的字段校验与确认入口。

## 适用范围

当前第一版只支持以下三类意图：

- `create_event`
- `query_events`
- `delete_event`

当前第一版只覆盖以下高频中文时间表达：

- 今天、明天、后天
- 上午、下午、晚上
- 几点、几点半
- 周几、下周几

以下内容暂不纳入第一版解析器范围：

- 修改事件
- 重复事件
- 持续时长
- 地点、参与人等扩展字段
- 复杂相对时间（如“下下周二傍晚”）

## 总体思路

第一版解析器采用**规则优先**设计，不依赖端侧大模型。

推荐使用纯 Kotlin 模块实现，解析链路如下：

```text
用户语音 → Moonshine ASR 转写文本 → 本地解析器 → 结构化 JSON → 校验 → 确认/执行
```

解析器输出结构必须与 `docs/protocol.md` 保持一致，确保协议层与执行层解耦。

## 模块拆分

### 1. `IntentDetector`

负责识别用户当前意图。

支持的 intent：
- `create_event`
- `query_events`
- `delete_event`

建议采用关键字优先匹配：

- 删除类：`删除`、`取消`、`移除`、`去掉`
- 创建类：`添加`、`新增`、`创建`、`安排`、`记下`
- 查询类：`查看`、`查询`、`看看`、`有什么安排`、`有哪些安排`

建议优先级：

```text
delete_event > create_event > query_events
```

原因：删除类操作风险最高，应优先识别并走确认流程。

### 2. `DateTimeParser`

负责从 ASR 文本中提取并归一化日期、时间信息。

建议支持的规则：

#### 日期规则
- `今天` → 当前日期
- `明天` → 当前日期 + 1 天
- `后天` → 当前日期 + 2 天
- `周X` / `星期X` → 最近一次将来的对应星期日期
- `下周X` / `下星期X` → 下一自然周中的对应星期日期

#### 时间规则
- `X点` → `HH:00`
- `X点半` → `HH:30`
- `上午` / `早上` + `X点` → 上午时间
- `下午` + `X点` → 小时转为 24 小时制
- `晚上` + `X点` → 晚间时间，按 24 小时制输出

建议解析结果输出为：
- `date`: `YYYY-MM-DD`
- `time`: `HH:MM`

日期换算应依赖设备当前时间，不应写死示例日期。

### 3. `TitleExtractor`

负责从原始文本中提取事件标题。

推荐策略：
1. 先识别并移除 intent 关键词
2. 再移除已识别的日期/时间短语
3. 对剩余文本进行 trim 和简单清洗
4. 剩余文本作为标题候选

示例：

- `明天下午三点和张三开会` → `和张三开会`
- `后天上午十点去超市买东西` → `去超市买东西`
- `取消今晚七点的健身` → `健身`

如果剩余文本为空，则应返回缺失标题错误。

### 4. `ParserValidator`

负责对解析结果做字段级校验。

校验规则应直接复用 `docs/protocol.md`：

- `intent` 必须属于允许范围
- `create_event` 必须包含 `title`、`date`
- `time` 可选
- `delete_event` 必须包含 `title`
- `date` 必须符合 `YYYY-MM-DD`
- `time` 必须符合 `HH:MM`

如果校验失败，解析器不应直接返回“可执行结果”，而应返回明确的错误类型或缺失信息。

### 5. `CommandParser`

作为总入口，按顺序协调解析流程：

1. 识别 intent
2. 提取日期/时间
3. 提取标题
4. 校验字段
5. 输出结构化结果或错误

建议对外提供统一接口，例如：

```kotlin
parse(asrText: String): ParserResult
```

## 推荐的数据结构

可在实现时采用类似的数据结构：

```kotlin
data class ParsedCommand(
    val intent: String,
    val title: String?,
    val date: String?,
    val time: String?,
    val needReminder: Boolean?
)
```

以及统一错误结果：

```kotlin
sealed class ParserResult {
    data class Success(val command: ParsedCommand) : ParserResult()
    data class Error(val reason: String, val message: String) : ParserResult()
}
```

## 错误处理

第一版至少需要覆盖以下错误情况：

- 输入为空
- 无法识别 intent
- 创建事件缺少日期
- 删除事件缺少标题
- 时间表达无法解析
- 标题为空

对于这些情况，建议返回结构化错误，而不是勉强输出不完整 JSON。

## 确认策略

解析器负责产出结构化结果，执行层负责确认。

建议策略：

- `delete_event`：始终二次确认
- `create_event`：若字段完整可直接进入确认页
- `query_events`：通常无需确认，可直接查询

也就是说，本地解析器不直接决定是否执行，只负责把结果整理为可确认、可校验的数据。

## 与现有文档的关系

- `docs/protocol.md`：定义解析器输出契约
- `docs/examples.md`：定义验收输入/输出样例
- `README.md`：描述整体架构与开发路径

本文件用于替代旧的 LLM Prompt 文档，成为本地解析层的设计说明。

## 测试建议

`docs/examples.md` 中现有 8 条示例可直接转化为解析器测试用例，至少验证：

- intent 是否正确
- 日期是否正确换算
- 时间是否正确归一化
- 标题是否正确提取
- 输出是否符合协议

建议实现后按以下维度补充单元测试：

- 正常路径：新增 / 查询 / 删除
- 缺字段路径：缺时间、缺标题
- 边界路径：只有日期、只有时间、歧义表达
- 错误路径：空输入、无法识别意图
