enum EventStatus { active, deleted }

class CalendarEvent {
  CalendarEvent({
    required this.id,
    required this.title,
    required this.date,
    this.time,
    this.needReminder = false,
    this.status = EventStatus.active,
    required this.createdAt,
    this.deletedAt,
  });

  final String id;
  final String title;
  final DateTime date;
  final TimeOfDayValue? time;
  final bool needReminder;
  final EventStatus status;
  final DateTime createdAt;
  final DateTime? deletedAt;

  CalendarEvent copyWith({
    String? id,
    String? title,
    DateTime? date,
    TimeOfDayValue? time,
    bool? needReminder,
    EventStatus? status,
    DateTime? createdAt,
    DateTime? deletedAt,
  }) {
    return CalendarEvent(
      id: id ?? this.id,
      title: title ?? this.title,
      date: date ?? this.date,
      time: time ?? this.time,
      needReminder: needReminder ?? this.needReminder,
      status: status ?? this.status,
      createdAt: createdAt ?? this.createdAt,
      deletedAt: deletedAt ?? this.deletedAt,
    );
  }
}

class TimeOfDayValue {
  const TimeOfDayValue({required this.hour, required this.minute});

  final int hour;
  final int minute;

  String format() {
    final hh = hour.toString().padLeft(2, '0');
    final mm = minute.toString().padLeft(2, '0');
    return '$hh:$mm';
  }
}
