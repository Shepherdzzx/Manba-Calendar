import '../models/calendar_event.dart';
import '../models/parsed_command.dart';

DateTime? parsedDateFromCommand(ParsedCommand command) {
  final value = command.date;
  if (value == null || value.isEmpty) {
    return null;
  }
  return DateTime.tryParse(value);
}

TimeOfDayValue? parsedTimeFromCommand(ParsedCommand command) {
  final value = command.time;
  if (value == null || value.isEmpty) {
    return null;
  }
  final parts = value.split(':');
  if (parts.length != 2) {
    return null;
  }
  final hour = int.tryParse(parts[0]);
  final minute = int.tryParse(parts[1]);
  if (hour == null || minute == null) {
    return null;
  }
  if (hour < 0 || hour > 23 || minute < 0 || minute > 59) {
    return null;
  }
  return TimeOfDayValue(hour: hour, minute: minute);
}
