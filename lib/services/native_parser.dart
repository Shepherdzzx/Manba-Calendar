import 'dart:convert';

import 'package:flutter/services.dart';

import '../models/parsed_command.dart';

class NativeParser {
  const NativeParser();

  static const _channel = MethodChannel('com.example.manba_alert/parser');

  Future<ParsedCommand> parseText(String input, {DateTime? now}) async {
    final payload = await _channel.invokeMethod<String>('parseText', {
      'input': input,
      'now': (now ?? DateTime.now()).toUtc().toIso8601String(),
    });
    if (payload == null || payload.isEmpty) {
      throw const FormatException('parseText returned empty payload');
    }

    final json = jsonDecode(payload) as Map<String, dynamic>;
    return ParsedCommand.fromJson(json);
  }
}
