import 'package:flutter_test/flutter_test.dart';
import 'package:manba_alert/models/parsed_command.dart';
import 'package:manba_alert/utils/parsed_command_mapper.dart';

void main() {
  test('parsedDateFromCommand parses ISO date strings', () {
    const command = ParsedCommand(
      intent: ParsedCommand.intentCreateEvent,
      date: '2026-05-30',
    );

    final parsed = parsedDateFromCommand(command);

    expect(parsed, isNotNull);
    expect(parsed, DateTime(2026, 5, 30));
  });

  test('parsedTimeFromCommand parses HH:MM strings', () {
    const command = ParsedCommand(
      intent: ParsedCommand.intentCreateEvent,
      time: '15:00',
    );

    final parsed = parsedTimeFromCommand(command);

    expect(parsed, isNotNull);
    expect(parsed?.hour, 15);
    expect(parsed?.minute, 0);
  });

  test('parsedTimeFromCommand returns null for invalid values', () {
    const command = ParsedCommand(
      intent: ParsedCommand.intentCreateEvent,
      time: '25:99',
    );

    final parsed = parsedTimeFromCommand(command);

    expect(parsed, isNull);
  });
}
