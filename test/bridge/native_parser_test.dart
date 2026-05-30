import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:manba_alert/models/parsed_command.dart';
import 'package:manba_alert/services/native_parser.dart';

void main() {
  TestWidgetsFlutterBinding.ensureInitialized();

  const channel = MethodChannel('com.example.manba_alert/parser');
  final log = <MethodCall>[];

  setUp(() {
    log.clear();
    TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger
        .setMockMethodCallHandler(channel, (call) async {
          log.add(call);
          return '{"intent":"create_event","title":"和张三开会","date":"2026-05-30","time":"15:00","need_reminder":true}';
        });
  });

  tearDown(() {
    TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger
        .setMockMethodCallHandler(channel, null);
  });

  test('parseText returns ParsedCommand from channel payload', () async {
    const parser = NativeParser();

    final command = await parser.parseText('明天下午三点和张三开会');

    expect(command.intent, ParsedCommand.intentCreateEvent);
    expect(command.title, '和张三开会');
    expect(command.date, '2026-05-30');
    expect(command.time, '15:00');
    expect(command.needReminder, isTrue);
    expect(log.single.method, 'parseText');
    expect(log.single.arguments['input'], '明天下午三点和张三开会');
  });
}
