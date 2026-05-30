import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../models/calendar_event.dart';
import '../providers/event_state.dart';
import '../providers/voice_state.dart';
import '../utils/date_formatters.dart';
import '../widgets/calendar_panel.dart';
import '../widgets/event_card.dart';
import '../widgets/voice_widgets.dart';
import 'event_editor_screen.dart';
import 'settings_screen.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => EventState()),
        ChangeNotifierProvider(create: (_) => VoiceState()),
      ],
      child: const _HomeScaffold(),
    );
  }
}

class _HomeScaffold extends StatelessWidget {
  const _HomeScaffold();

  @override
  Widget build(BuildContext context) {
    final eventState = context.watch<EventState>();
    final theme = Theme.of(context);
    final voiceState = context.watch<VoiceState>();

    return Scaffold(
      appBar: AppBar(
        toolbarHeight: 72,
        titleSpacing: 0,
        title: Transform.translate(
          offset: const Offset(0, 0),
          child: Align(
            alignment: Alignment.centerLeft,
            child: SizedBox(
              height: 250,
              child: Image.asset(
                'assets/images/manba_calendar_logo.png',
                fit: BoxFit.contain,
                alignment: Alignment.centerLeft,
              ),
            ),
          ),
        ),
        actions: [
          IconButton(
            onPressed: () {
              Navigator.of(
                context,
              ).push(MaterialPageRoute(builder: (_) => const SettingsScreen()));
            },
            icon: const Icon(Icons.palette_outlined, size: 40),
          ),
        ],
      ),
      body: ListView(
        padding: const EdgeInsets.fromLTRB(20, 12, 20, 24),
        children: [
          CalendarPanel(
            state: eventState,
            onSelect: eventState.selectDate,
            onPageChanged: eventState.setFocusedMonth,
          ),
          const SizedBox(height: 18),
          Text(
            '${formatDateLabel(eventState.selectedDate)} 的安排',
            style: theme.textTheme.titleLarge?.copyWith(
              fontWeight: FontWeight.w700,
            ),
          ),
          const SizedBox(height: 12),
          if (eventState.selectedDateEvents.isEmpty)
            Card(
              child: Padding(
                padding: const EdgeInsets.all(20),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: const [
                    Text('这一天还没有安排'),
                    SizedBox(height: 8),
                    Text('你可以手动新增事件，或者长按下方按钮说一句自然语言来创建安排。'),
                  ],
                ),
              ),
            )
          else
            ...eventState.selectedDateEvents.map(
              (event) => Padding(
                padding: const EdgeInsets.only(bottom: 12),
                child: EventCard(
                  event: event,
                  onEdit: () =>
                      _openEditor(context, eventState, initialEvent: event),
                  onDelete: () => _confirmDelete(context, eventState, event),
                ),
              ),
            ),
          const SizedBox(height: 20),
          VoiceRecordButton(
            onLongPressStart: voiceState.startRecording,
            onLongPressEnd: () => voiceState.finishRecording(),
            label: _voiceLabel(voiceState.status),
            icon: _voiceIcon(voiceState.status),
          ),
          const SizedBox(height: 12),
          Text(
            '按住识别，松开后进行解析。当前版本使用 stub 结果打通演示链路。',
            style: theme.textTheme.bodySmall,
          ),
          const SizedBox(height: 16),
          if (voiceState.status == VoiceStatus.ready &&
              voiceState.transcript != null &&
              voiceState.command != null)
            VoiceResultCard(
              transcript: voiceState.transcript!,
              command: voiceState.command!,
              onConfirm: () {
                eventState.addEvent(
                  title: voiceState.command!.title ?? '新的安排',
                  date: DateTime.now().add(const Duration(days: 1)),
                  time: const TimeOfDayValue(hour: 15, minute: 0),
                  needReminder: voiceState.command!.needReminder ?? false,
                );
                voiceState.reset();
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('已根据解析结果写入一条演示事件')),
                );
              },
              onCancel: voiceState.reset,
            )
          else if (voiceState.status == VoiceStatus.processing)
            const Card(
              child: Padding(
                padding: EdgeInsets.all(20),
                child: Row(
                  children: [
                    SizedBox(
                      width: 20,
                      height: 20,
                      child: CircularProgressIndicator(strokeWidth: 2),
                    ),
                    SizedBox(width: 12),
                    Expanded(child: Text('正在解析语音内容，请稍候…')),
                  ],
                ),
              ),
            ),
        ],
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => _openEditor(
          context,
          eventState,
          initialDate: eventState.selectedDate,
        ),
        icon: const Icon(Icons.add),
        label: const Text('新增事件'),
      ),
    );
  }

  Future<void> _openEditor(
    BuildContext context,
    EventState state, {
    DateTime? initialDate,
    CalendarEvent? initialEvent,
  }) async {
    final result = await Navigator.of(context).push<EventEditorResult>(
      MaterialPageRoute(
        builder: (_) => EventEditorScreen(
          initialDate: initialDate,
          initialEvent: initialEvent,
        ),
      ),
    );
    if (result == null) {
      return;
    }
    if (result.id == null) {
      state.addEvent(
        title: result.title,
        date: result.date,
        time: result.time,
        needReminder: result.needReminder,
      );
      return;
    }
    state.updateEvent(
      CalendarEvent(
        id: result.id!,
        title: result.title,
        date: result.date,
        time: result.time,
        needReminder: result.needReminder,
        createdAt: result.createdAt ?? DateTime.now(),
      ),
    );
  }

  Future<void> _confirmDelete(
    BuildContext context,
    EventState state,
    CalendarEvent event,
  ) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (dialogContext) {
        return AlertDialog(
          title: const Text('确认删除'),
          content: Text('确定删除“${event.title}”吗？这是一个软删除操作。'),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(dialogContext).pop(false),
              child: const Text('取消'),
            ),
            FilledButton(
              onPressed: () => Navigator.of(dialogContext).pop(true),
              child: const Text('删除'),
            ),
          ],
        );
      },
    );
    if (confirmed == true) {
      state.deleteEvent(event.id);
    }
  }

  static String _voiceLabel(VoiceStatus status) {
    switch (status) {
      case VoiceStatus.idle:
        return '按住说话';
      case VoiceStatus.recording:
        return '松开进行解析';
      case VoiceStatus.processing:
        return '解析中…';
      case VoiceStatus.ready:
        return '再次按住重新识别';
    }
  }

  static IconData _voiceIcon(VoiceStatus status) {
    switch (status) {
      case VoiceStatus.idle:
        return Icons.mic_none_rounded;
      case VoiceStatus.recording:
        return Icons.mic_rounded;
      case VoiceStatus.processing:
        return Icons.auto_awesome;
      case VoiceStatus.ready:
        return Icons.check_circle_outline;
    }
  }
}
