import '../models/parsed_command.dart';

class VoiceParseResult {
  const VoiceParseResult({
    required this.transcript,
    this.command,
    this.error,
  });

  final String transcript;
  final ParsedCommand? command;
  final String? error;

  bool get isSuccess => command != null && error == null;
  bool get isError => error != null;
}

abstract class VoiceParserService {
  Future<VoiceParseResult> parseAfterRelease();
}

class VoiceStubService implements VoiceParserService {
  @override
  Future<VoiceParseResult> parseAfterRelease() async {
    await Future<void>.delayed(const Duration(milliseconds: 700));
    return const VoiceParseResult(
      transcript: '明天下午三点和张三开会',
      command: ParsedCommand(
        intent: ParsedCommand.intentCreateEvent,
        title: '和张三开会',
        date: '2026-05-30',
        time: '15:00',
        needReminder: true,
      ),
    );
  }
}
