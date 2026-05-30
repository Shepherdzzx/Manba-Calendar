import '../models/parsed_command.dart';

class VoiceStubResult {
  const VoiceStubResult({required this.transcript, required this.command});

  final String transcript;
  final ParsedCommand command;
}

class VoiceService {
  Future<VoiceStubResult> parseAfterRelease() async {
    await Future<void>.delayed(const Duration(milliseconds: 700));
    return const VoiceStubResult(
      transcript: '明天下午三点和张三开会',
      command: ParsedCommand(
        intent: 'create_event',
        title: '和张三开会',
        dateLabel: '明天',
        timeLabel: '15:00',
        needReminder: true,
      ),
    );
  }
}
