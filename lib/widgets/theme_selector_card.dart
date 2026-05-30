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
        padding: const EdgeInsets.all(10),
        child: Column(
          children: themes.map((theme) {
            final selected = theme.id == currentThemeId;
            return ListTile(
              visualDensity: VisualDensity.compact,
              contentPadding: const EdgeInsets.symmetric(
                horizontal: 4,
                vertical: 2,
              ),
              onTap: () => onThemeSelected(theme.id),
              leading: Icon(
                selected ? Icons.radio_button_checked : Icons.radio_button_off,
              ),
              title: Row(
                children: [
                  Expanded(child: Text(theme.label)),
                  _ColorDot(color: theme.primary),
                  const SizedBox(width: 4),
                  _ColorDot(color: theme.secondary),
                  if (selected) ...[
                    const SizedBox(width: 4),
                    const Icon(Icons.check_circle, size: 16),
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
      width: 14,
      height: 14,
      decoration: BoxDecoration(color: color, shape: BoxShape.circle),
    );
  }
}
