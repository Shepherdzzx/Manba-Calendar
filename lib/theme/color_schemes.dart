import 'package:flutter/material.dart';

class AppThemeOption {
  const AppThemeOption({
    required this.id,
    required this.label,
    required this.seedColor,
    required this.primary,
    required this.secondary,
    required this.surfaceTint,
  });

  final String id;
  final String label;
  final Color seedColor;
  final Color primary;
  final Color secondary;
  final Color surfaceTint;
}

const purpleGoldTheme = AppThemeOption(
  id: 'purple_gold',
  label: '紫金',
  seedColor: Color(0xFF6B3FA0),
  primary: Color(0xFF5D2E8C),
  secondary: Color(0xFFD4A937),
  surfaceTint: Color(0xFFF4ECFB),
);

const tealCreamTheme = AppThemeOption(
  id: 'teal_cream',
  label: '青瓷',
  seedColor: Color(0xFF2E7D7A),
  primary: Color(0xFF256B68),
  secondary: Color(0xFFE7D9B7),
  surfaceTint: Color(0xFFF1FBFA),
);

const coralNightTheme = AppThemeOption(
  id: 'coral_night',
  label: '珊瑚夜',
  seedColor: Color(0xFF9C4A64),
  primary: Color(0xFF8A3553),
  secondary: Color(0xFFF4B183),
  surfaceTint: Color(0xFFFFF4F2),
);

const themeOptions = <AppThemeOption>[
  purpleGoldTheme,
  tealCreamTheme,
  coralNightTheme,
];
