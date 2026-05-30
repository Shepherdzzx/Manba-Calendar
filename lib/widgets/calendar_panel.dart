import 'package:flutter/material.dart';
import 'package:table_calendar/table_calendar.dart';

import '../providers/event_state.dart';

class CalendarPanel extends StatelessWidget {
  const CalendarPanel({
    super.key,
    required this.state,
    required this.onSelect,
    required this.onPageChanged,
  });

  final EventState state;
  final ValueChanged<DateTime> onSelect;
  final ValueChanged<DateTime> onPageChanged;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: TableCalendar<Object>(
          firstDay: DateTime.utc(2020, 1, 1),
          lastDay: DateTime.utc(2035, 12, 31),
          focusedDay: state.focusedMonth,
          selectedDayPredicate: (day) => isSameDay(day, state.selectedDate),
          availableCalendarFormats: const {CalendarFormat.month: 'Month'},
          daysOfWeekHeight: 34,
          daysOfWeekStyle: const DaysOfWeekStyle(
            weekdayStyle: TextStyle(fontSize: 12),
            weekendStyle: TextStyle(fontSize: 12),
          ),
          calendarBuilders: CalendarBuilders(
            dowBuilder: (context, day) {
              const labels = ['一', '二', '三', '四', '五', '六', '日'];
              return Center(child: Text(labels[day.weekday - 1]));
            },
            headerTitleBuilder: (context, day) {
              return Text(
                '${day.year}年${day.month}月',
                style: theme.textTheme.titleMedium?.copyWith(
                  fontWeight: FontWeight.w700,
                ),
              );
            },
          ),
          onDaySelected: (selectedDay, focusedDay) {
            onSelect(selectedDay);
            onPageChanged(focusedDay);
          },
          onPageChanged: onPageChanged,
          eventLoader: (day) =>
              List.generate(state.countForDay(day), (_) => Object()),
          calendarStyle: CalendarStyle(
            markersMaxCount: 3,
            markerDecoration: BoxDecoration(
              color: theme.colorScheme.secondary,
              shape: BoxShape.circle,
            ),
            selectedDecoration: BoxDecoration(
              color: theme.colorScheme.primary,
              shape: BoxShape.circle,
            ),
            todayDecoration: BoxDecoration(
              color: theme.colorScheme.primary.withValues(alpha: 0.3),
              shape: BoxShape.circle,
            ),
          ),
          headerStyle: const HeaderStyle(
            formatButtonVisible: false,
            titleCentered: true,
          ),
        ),
      ),
    );
  }
}
