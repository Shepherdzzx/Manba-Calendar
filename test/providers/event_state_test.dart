import 'package:flutter_test/flutter_test.dart';
import 'package:manba_alert/providers/event_state.dart';

void main() {
  test('can add and soft delete events', () {
    final state = EventState();
    final initialCount = state.allEvents.length;

    state.addEvent(
      title: '测试新增事件',
      date: DateTime(2026, 5, 29),
      needReminder: true,
    );

    expect(state.allEvents.length, initialCount + 1);
    final created = state.allEvents.firstWhere(
      (event) => event.title == '测试新增事件',
    );

    state.deleteEvent(created.id);

    expect(state.allEvents.any((event) => event.id == created.id), isFalse);
  });

  test('selectDate updates selected day events', () {
    final state = EventState();
    state.addEvent(
      title: '当天事件',
      date: DateTime(2026, 6, 1),
      needReminder: false,
    );

    state.selectDate(DateTime(2026, 6, 1));

    expect(
      state.selectedDateEvents.any((event) => event.title == '当天事件'),
      isTrue,
    );
  });
}
