import 'package:flutter/material.dart';

import '../models/calendar_event.dart';

class EventEditorScreen extends StatefulWidget {
  const EventEditorScreen({super.key, this.initialDate, this.initialEvent});

  final DateTime? initialDate;
  final CalendarEvent? initialEvent;

  @override
  State<EventEditorScreen> createState() => _EventEditorScreenState();
}

class _EventEditorScreenState extends State<EventEditorScreen> {
  final _formKey = GlobalKey<FormState>();
  late final TextEditingController _titleController;
  late DateTime _selectedDate;
  TimeOfDay? _selectedTime;
  late bool _needReminder;

  @override
  void initState() {
    super.initState();
    _titleController = TextEditingController(
      text: widget.initialEvent?.title ?? '',
    );
    _selectedDate = DateUtils.dateOnly(
      widget.initialEvent?.date ?? widget.initialDate ?? DateTime.now(),
    );
    final time = widget.initialEvent?.time;
    if (time != null) {
      _selectedTime = TimeOfDay(hour: time.hour, minute: time.minute);
    }
    _needReminder = widget.initialEvent?.needReminder ?? false;
  }

  @override
  void dispose() {
    _titleController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final editing = widget.initialEvent != null;

    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back_rounded, size: 30),
          onPressed: () => Navigator.of(context).pop(),
        ),
        title: Text(editing ? '编辑事件' : '新增事件'),
      ),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(20),
          children: [
            TextFormField(
              controller: _titleController,
              decoration: const InputDecoration(labelText: '事件标题'),
              validator: (value) {
                if (value == null || value.trim().isEmpty) {
                  return '请输入事件标题';
                }
                return null;
              },
            ),
            const SizedBox(height: 16),
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  children: [
                    ListTile(
                      contentPadding: EdgeInsets.zero,
                      leading: const Icon(Icons.calendar_month_rounded, size: 30),
                      title: const Text('日期'),
                      subtitle: Text(
                        '${_selectedDate.year}-${_selectedDate.month.toString().padLeft(2, '0')}-${_selectedDate.day.toString().padLeft(2, '0')}',
                      ),
                      trailing: const Icon(Icons.chevron_right_rounded),
                      onTap: _pickDate,
                    ),
                    ListTile(
                      contentPadding: EdgeInsets.zero,
                      leading: const Icon(Icons.schedule_rounded, size: 30),
                      title: const Text('时间'),
                      subtitle: Text(
                        _selectedTime == null ? '全天' : _selectedTime!.format(context),
                      ),
                      trailing: const Icon(Icons.chevron_right_rounded),
                      onTap: _pickTime,
                    ),
                    SwitchListTile(
                      contentPadding: EdgeInsets.zero,
                      secondary: const Icon(Icons.notifications_rounded, size: 30),
                      value: _needReminder,
                      onChanged: (value) => setState(() => _needReminder = value),
                      title: const Text('开启提醒'),
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 24),
            FilledButton(
              style: FilledButton.styleFrom(
                minimumSize: const Size.fromHeight(56),
                padding: const EdgeInsets.symmetric(vertical: 16),
              ),
              onPressed: () {
                if (!_formKey.currentState!.validate()) {
                  return;
                }
                Navigator.of(context).pop(
                  EventEditorResult(
                    id: widget.initialEvent?.id,
                    title: _titleController.text.trim(),
                    date: _selectedDate,
                    time: _selectedTime == null
                        ? null
                        : TimeOfDayValue(
                            hour: _selectedTime!.hour,
                            minute: _selectedTime!.minute,
                          ),
                    needReminder: _needReminder,
                    createdAt: widget.initialEvent?.createdAt,
                  ),
                );
              },
              child: Text(
                editing ? '保存修改' : '创建事件',
                style: theme.textTheme.titleMedium?.copyWith(
                  fontWeight: FontWeight.w600,
                  color: Colors.white,
                ),
              ),
            ),
            const SizedBox(height: 12),
            Text(
              '第一版聚焦高频日程记录场景，支持标题、日期、时间与提醒开关。',
              style: theme.textTheme.bodySmall,
            ),
          ],
        ),
      ),
    );
  }

  Future<void> _pickDate() async {
    final picked = await showDatePicker(
      context: context,
      initialDate: _selectedDate,
      firstDate: DateTime(2020),
      lastDate: DateTime(2035),
    );
    if (picked != null) {
      setState(() => _selectedDate = DateUtils.dateOnly(picked));
    }
  }

  Future<void> _pickTime() async {
    final picked = await showTimePicker(
      context: context,
      initialTime: _selectedTime ?? TimeOfDay.now(),
    );
    if (picked != null) {
      setState(() => _selectedTime = picked);
    }
  }
}

class EventEditorResult {
  const EventEditorResult({
    this.id,
    required this.title,
    required this.date,
    this.time,
    required this.needReminder,
    this.createdAt,
  });

  final String? id;
  final String title;
  final DateTime date;
  final TimeOfDayValue? time;
  final bool needReminder;
  final DateTime? createdAt;
}
