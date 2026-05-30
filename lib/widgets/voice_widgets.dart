import 'package:flutter/material.dart';

import '../models/parsed_command.dart';

class VoiceRecordButton extends StatelessWidget {
  const VoiceRecordButton({
    super.key,
    required this.onLongPressStart,
    required this.onLongPressEnd,
    required this.label,
    required this.icon,
  });

  final VoidCallback onLongPressStart;
  final VoidCallback onLongPressEnd;
  final String label;
  final IconData icon;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return GestureDetector(
      onLongPressStart: (_) => onLongPressStart(),
      onLongPressEnd: (_) => onLongPressEnd(),
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 180),
        width: double.infinity,
        padding: const EdgeInsets.symmetric(vertical: 18, horizontal: 20),
        decoration: BoxDecoration(
          gradient: LinearGradient(
            colors: [theme.colorScheme.primary, theme.colorScheme.secondary],
          ),
          borderRadius: BorderRadius.circular(24),
          boxShadow: [
            BoxShadow(
              color: theme.colorScheme.primary.withValues(alpha: 0.24),
              blurRadius: 18,
              offset: const Offset(0, 10),
            ),
          ],
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(icon, color: Colors.white),
            const SizedBox(width: 10),
            Text(
              label,
              style: theme.textTheme.titleMedium?.copyWith(
                color: Colors.white,
                fontWeight: FontWeight.w700,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class VoiceResultCard extends StatelessWidget {
  const VoiceResultCard({
    super.key,
    required this.transcript,
    required this.command,
    required this.onConfirm,
    required this.onCancel,
  });

  final String transcript;
  final ParsedCommand command;
  final VoidCallback onConfirm;
  final VoidCallback onCancel;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              '识别结果',
              style: theme.textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.w700,
              ),
            ),
            const SizedBox(height: 8),
            Text(transcript),
            const SizedBox(height: 16),
            Text(
              '解析结果',
              style: theme.textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.w700,
              ),
            ),
            const SizedBox(height: 8),
            Text('Intent: ${command.intent}'),
            if (command.title != null) Text('标题: ${command.title}'),
            if (command.dateLabel != null) Text('日期: ${command.dateLabel}'),
            if (command.timeLabel != null) Text('时间: ${command.timeLabel}'),
            if (command.needReminder != null)
              Text('提醒: ${command.needReminder! ? '开启' : '关闭'}'),
            const SizedBox(height: 20),
            Row(
              children: [
                Expanded(
                  child: FilledButton(
                    onPressed: onConfirm,
                    child: const Text('确认写入日历'),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: OutlinedButton(
                    onPressed: onCancel,
                    child: const Text('取消'),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
