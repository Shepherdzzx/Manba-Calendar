import 'package:flutter/material.dart';

import '../theme/color_schemes.dart';

class ThemeSelectorCard extends StatelessWidget {
  const ThemeSelectorCard({
    super.key,
    required this.currentThemeId,
    required this.themes,
    required this.onThemeSelected,
  });

  final String currentThemeId;
  final List<AppThemeOption> themes;
  final ValueChanged<String> onThemeSelected;

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          children: themes.map((theme) {
            final selected = theme.id == currentThemeId;
            return ListTile(
              onTap: () => onThemeSelected(theme.id),
              leading: Icon(
                selected ? Icons.radio_button_checked : Icons.radio_button_off,
              ),
              title: Row(
                children: [
                  Expanded(child: Text(theme.label)),
                  _ColorDot(color: theme.primary),
                  const SizedBox(width: 8),
                  _ColorDot(color: theme.secondary),
                  if (selected) ...[
                    const SizedBox(width: 8),
                    const Icon(Icons.check_circle, size: 18),
                  ],
                ],
              ),
            );
          }).toList(),
        ),
      ),
    );
  }
}

class _ColorDot extends StatelessWidget {
  const _ColorDot({required this.color});

  final Color color;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 18,
      height: 18,
      decoration: BoxDecoration(color: color, shape: BoxShape.circle),
    );
  }
}
