import 'package:flutter/foundation.dart';

import '../models/parsed_command.dart';
import 'native_parser.dart';

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
  VoiceStubService({NativeParser? nativeParser, bool? useNativeParser})
    : _nativeParser = nativeParser ?? const NativeParser(),
      _useNativeParser = useNativeParser ?? defaultTargetPlatform == TargetPlatform.android;

  final NativeParser _nativeParser;
  final bool _useNativeParser;

  @override
  Future<VoiceParseResult> parseAfterRelease() async {
    await Future<void>.delayed(const Duration(milliseconds: 700));
    const transcript = '明天下午三点和张三开会';

    if (_useNativeParser) {
      try {
        final command = await _nativeParser.parseText(transcript);
        return VoiceParseResult(transcript: transcript, command: command);
      } catch (error) {
        return const VoiceParseResult(
          transcript: transcript,
          command: ParsedCommand(
            intent: ParsedCommand.intentCreateEvent,
            title: '和张三开会',
            date: '2026-05-30',
            time: '15:00',
            needReminder: true,
          ),
          error: 'native parser unavailable, fell back to stub result',
        );
      }
    }

    return const VoiceParseResult(
      transcript: transcript,
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
