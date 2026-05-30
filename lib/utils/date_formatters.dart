import '../models/calendar_event.dart';

const _weekdayLabels = <int, String>{
  DateTime.monday: '周一',
  DateTime.tuesday: '周二',
  DateTime.wednesday: '周三',
  DateTime.thursday: '周四',
  DateTime.friday: '周五',
  DateTime.saturday: '周六',
  DateTime.sunday: '周日',
};

String formatDateLabel(DateTime date) {
  final weekday = _weekdayLabels[date.weekday] ?? '';
  return '${date.month}月${date.day}日 $weekday';
}

String formatMonthLabel(DateTime date) {
  return '${date.year}年${date.month}月';
}

String formatEventTime(TimeOfDayValue? time) {
  if (time == null) {
    return '全天';
  }
  return time.format();
}
