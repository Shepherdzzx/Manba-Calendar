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

const skyMintTheme = AppThemeOption(
  id: 'sky_mint',
  label: '天青薄荷',
  seedColor: Color(0xFF4C8DAE),
  primary: Color(0xFF3B7493),
  secondary: Color(0xFF86D5C2),
  surfaceTint: Color(0xFFF2FBFD),
);

const roseSandTheme = AppThemeOption(
  id: 'rose_sand',
  label: '玫砂',
  seedColor: Color(0xFFB05A72),
  primary: Color(0xFF96485E),
  secondary: Color(0xFFE7C8A0),
  surfaceTint: Color(0xFFFFF5F6),
);

const forestAmberTheme = AppThemeOption(
  id: 'forest_amber',
  label: '森琥珀',
  seedColor: Color(0xFF3F6C4E),
  primary: Color(0xFF335A40),
  secondary: Color(0xFFD7A441),
  surfaceTint: Color(0xFFF5FAF6),
);

const duskLavenderTheme = AppThemeOption(
  id: 'dusk_lavender',
  label: '暮薰衣',
  seedColor: Color(0xFF6E5EA8),
  primary: Color(0xFF5A4D90),
  secondary: Color(0xFFC9A8E5),
  surfaceTint: Color(0xFFF7F3FC),
);

const themeOptions = <AppThemeOption>[
  purpleGoldTheme,
  tealCreamTheme,
  coralNightTheme,
  skyMintTheme,
  roseSandTheme,
  forestAmberTheme,
  duskLavenderTheme,
];
