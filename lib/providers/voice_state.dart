import 'package:flutter/foundation.dart';

import '../models/parsed_command.dart';
import '../services/voice_service.dart';

enum VoiceStatus { idle, recording, processing, ready }

class VoiceState extends ChangeNotifier {
  VoiceState({VoiceService? service}) : _service = service ?? VoiceService();

  final VoiceService _service;
  VoiceStatus _status = VoiceStatus.idle;
  String? _transcript;
  ParsedCommand? _command;

  VoiceStatus get status => _status;
  String? get transcript => _transcript;
  ParsedCommand? get command => _command;

  void startRecording() {
    _status = VoiceStatus.recording;
    notifyListeners();
  }

  Future<void> finishRecording() async {
    _status = VoiceStatus.processing;
    notifyListeners();
    final result = await _service.parseAfterRelease();
    _transcript = result.transcript;
    _command = result.command;
    _status = VoiceStatus.ready;
    notifyListeners();
  }

  void reset() {
    _status = VoiceStatus.idle;
    _transcript = null;
    _command = null;
    notifyListeners();
  }
}
