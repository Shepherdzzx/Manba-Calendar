import 'dart:math';

import 'package:flutter/material.dart';

import '../models/calendar_event.dart';

class EventStore {
  final List<CalendarEvent> _events = [
    CalendarEvent(
      id: 'evt-demo-1',
      title: '和张三开会',
      date: DateTime.now(),
      time: const TimeOfDayValue(hour: 15, minute: 0),
      needReminder: true,
      createdAt: DateTime.now(),
    ),
    CalendarEvent(
      id: 'evt-demo-2',
      title: '去超市买东西',
      date: DateTime.now().add(const Duration(days: 1)),
      time: const TimeOfDayValue(hour: 10, minute: 0),
      createdAt: DateTime.now(),
    ),
  ];

  List<CalendarEvent> listActive() {
    final active = _events
        .where((event) => event.status == EventStatus.active)
        .toList();
    active.sort((a, b) {
      final dateCompare = _startOfDay(a.date).compareTo(_startOfDay(b.date));
      if (dateCompare != 0) {
        return dateCompare;
      }
      final aMinutes = a.time == null ? -1 : a.time!.hour * 60 + a.time!.minute;
      final bMinutes = b.time == null ? -1 : b.time!.hour * 60 + b.time!.minute;
      return aMinutes.compareTo(bMinutes);
    });
    return active;
  }

  List<CalendarEvent> eventsForDay(DateTime day) {
    final target = _startOfDay(day);
    return listActive()
        .where((event) => _startOfDay(event.date) == target)
        .toList();
  }

  CalendarEvent create({
    required String title,
    required DateTime date,
    TimeOfDayValue? time,
    required bool needReminder,
  }) {
    final event = CalendarEvent(
      id: _newId(),
      title: title,
      date: _startOfDay(date),
      time: time,
      needReminder: needReminder,
      createdAt: DateTime.now(),
    );
    _events.add(event);
    return event;
  }

  void update(CalendarEvent updated) {
    final index = _events.indexWhere((event) => event.id == updated.id);
    if (index == -1) {
      return;
    }
    _events[index] = updated;
  }

  void softDelete(String id) {
    final index = _events.indexWhere((event) => event.id == id);
    if (index == -1) {
      return;
    }
    _events[index] = _events[index].copyWith(
      status: EventStatus.deleted,
      deletedAt: DateTime.now(),
    );
  }

  String _newId() {
    return 'evt-${Random().nextInt(1 << 32).toRadixString(16)}';
  }

  DateTime _startOfDay(DateTime value) => DateUtils.dateOnly(value);
}

class EventState extends ChangeNotifier {
  EventState({EventStore? store}) : _store = store ?? EventStore();

  final EventStore _store;
  DateTime _focusedMonth = DateUtils.dateOnly(DateTime.now());
  DateTime _selectedDate = DateUtils.dateOnly(DateTime.now());

  DateTime get focusedMonth => _focusedMonth;
  DateTime get selectedDate => _selectedDate;
  List<CalendarEvent> get allEvents => _store.listActive();
  List<CalendarEvent> get selectedDateEvents =>
      _store.eventsForDay(_selectedDate);

  void selectDate(DateTime day) {
    _selectedDate = DateUtils.dateOnly(day);
    _focusedMonth = DateTime(day.year, day.month, 1);
    notifyListeners();
  }

  void setFocusedMonth(DateTime month) {
    _focusedMonth = DateTime(month.year, month.month, 1);
    notifyListeners();
  }

  void addEvent({
    required String title,
    required DateTime date,
    TimeOfDayValue? time,
    required bool needReminder,
  }) {
    _store.create(
      title: title,
      date: date,
      time: time,
      needReminder: needReminder,
    );
    _selectedDate = DateUtils.dateOnly(date);
    _focusedMonth = DateTime(date.year, date.month, 1);
    notifyListeners();
  }

  void updateEvent(CalendarEvent event) {
    _store.update(event);
    notifyListeners();
  }

  void deleteEvent(String id) {
    _store.softDelete(id);
    notifyListeners();
  }

  int countForDay(DateTime day) => _store.eventsForDay(day).length;
}
