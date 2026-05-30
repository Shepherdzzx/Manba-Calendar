class ParsedCommand {
  const ParsedCommand({
    required this.intent,
    this.title,
    this.date,
    this.time,
    this.needReminder,
  });

  static const intentCreateEvent = 'create_event';
  static const intentQueryEvents = 'query_events';
  static const intentDeleteEvent = 'delete_event';

  final String intent;
  final String? title;
  final String? date;
  final String? time;
  final bool? needReminder;

  factory ParsedCommand.fromJson(Map<String, dynamic> json) {
    return ParsedCommand(
      intent: json['intent'] as String,
      title: json['title'] as String?,
      date: json['date'] as String?,
      time: json['time'] as String?,
      needReminder: json['need_reminder'] as bool?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'intent': intent,
      if (title != null) 'title': title,
      if (date != null) 'date': date,
      if (time != null) 'time': time,
      if (needReminder != null) 'need_reminder': needReminder,
    };
  }

  String? validate() {
    const validIntents = {
      intentCreateEvent,
      intentQueryEvents,
      intentDeleteEvent,
    };
    if (!validIntents.contains(intent)) {
      return 'invalid intent: $intent';
    }

    if (date != null && date!.isNotEmpty) {
      final datePattern = RegExp(r'^\d{4}-\d{2}-\d{2}$');
      if (!datePattern.hasMatch(date!)) {
        return 'date must be in YYYY-MM-DD format';
      }
      if (DateTime.tryParse(date!) == null) {
        return 'date is not valid';
      }
    }

    if (time != null && time!.isNotEmpty) {
      final timePattern = RegExp(r'^\d{2}:\d{2}$');
      if (!timePattern.hasMatch(time!)) {
        return 'time must be in HH:MM format';
      }
      final parts = time!.split(':');
      final hour = int.tryParse(parts[0]);
      final minute = int.tryParse(parts[1]);
      if (hour == null || minute == null) {
        return 'time is not valid';
      }
      if (hour < 0 || hour > 23 || minute < 0 || minute > 59) {
        return 'time is not valid';
      }
    }

    if (intent == intentCreateEvent) {
      if (title == null || title!.isEmpty) {
        return 'create_event requires title';
      }
      if (date == null || date!.isEmpty) {
        return 'create_event requires date';
      }
    }

    if (intent == intentDeleteEvent) {
      final hasTitle = title != null && title!.isNotEmpty;
      final hasDate = date != null && date!.isNotEmpty;
      if (!hasTitle && !hasDate) {
        return 'delete_event requires title or date';
      }
    }

    return null;
  }
}
