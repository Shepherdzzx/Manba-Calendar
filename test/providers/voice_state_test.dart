import 'package:flutter_test/flutter_test.dart';
import 'package:manba_alert/models/parsed_command.dart';
import 'package:manba_alert/providers/voice_state.dart';
import 'package:manba_alert/services/voice_service.dart';

class _FakeVoiceParserService implements VoiceParserService {
  @override
  Future<VoiceParseResult> parseAfterRelease() async {
    return const VoiceParseResult(
      transcript: '帮我看看今天还有什么安排',
      command: ParsedCommand(
        intent: ParsedCommand.intentQueryEvents,
        date: '2026-05-30',
      ),
    );
  }
}

void main() {
  test('voice state transitions from recording to ready', () async {
    final state = VoiceState(service: _FakeVoiceParserService());

    state.startRecording();
    expect(state.status, VoiceStatus.recording);

    await state.finishRecording();

    expect(state.status, VoiceStatus.ready);
    expect(state.transcript, '帮我看看今天还有什么安排');
    expect(state.command?.intent, ParsedCommand.intentQueryEvents);
  });
}
