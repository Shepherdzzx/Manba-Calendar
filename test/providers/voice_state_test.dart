import 'package:flutter_test/flutter_test.dart';
import 'package:manba_alert/models/parsed_command.dart';
import 'package:manba_alert/providers/voice_state.dart';
import 'package:manba_alert/services/voice_service.dart';

class _FakeVoiceService extends VoiceService {
  @override
  Future<VoiceStubResult> parseAfterRelease() async {
    return const VoiceStubResult(
      transcript: '帮我看看今天还有什么安排',
      command: ParsedCommand(intent: 'query_events', dateLabel: '今天'),
    );
  }
}

void main() {
  test('voice state transitions from recording to ready', () async {
    final state = VoiceState(service: _FakeVoiceService());

    state.startRecording();
    expect(state.status, VoiceStatus.recording);

    await state.finishRecording();

    expect(state.status, VoiceStatus.ready);
    expect(state.transcript, '帮我看看今天还有什么安排');
    expect(state.command?.intent, 'query_events');
  });
}
