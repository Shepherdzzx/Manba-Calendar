class ParsedCommand {
  const ParsedCommand({
    required this.intent,
    this.title,
    this.dateLabel,
    this.timeLabel,
    this.needReminder,
  });

  final String intent;
  final String? title;
  final String? dateLabel;
  final String? timeLabel;
  final bool? needReminder;
}
