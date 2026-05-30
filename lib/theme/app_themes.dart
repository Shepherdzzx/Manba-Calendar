import 'package:flutter/material.dart';

import 'color_schemes.dart';

ThemeData buildTheme(AppThemeOption option) {
  final colorScheme = ColorScheme.fromSeed(
    seedColor: option.seedColor,
    primary: option.primary,
    secondary: option.secondary,
    surface: option.surfaceTint,
  );

  return ThemeData(
    useMaterial3: true,
    colorScheme: colorScheme,
    scaffoldBackgroundColor: option.surfaceTint,
    appBarTheme: AppBarTheme(
      backgroundColor: option.surfaceTint,
      foregroundColor: option.primary,
      centerTitle: false,
      elevation: 0,
    ),
    cardTheme: CardThemeData(
      color: Colors.white,
      elevation: 0,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(20)),
      margin: EdgeInsets.zero,
    ),
    floatingActionButtonTheme: FloatingActionButtonThemeData(
      backgroundColor: option.primary,
      foregroundColor: Colors.white,
    ),
    chipTheme: ChipThemeData(
      backgroundColor: option.secondary.withValues(alpha: 0.16),
      selectedColor: option.secondary.withValues(alpha: 0.24),
      labelStyle: TextStyle(color: option.primary, fontWeight: FontWeight.w600),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(999)),
    ),
    inputDecorationTheme: InputDecorationTheme(
      filled: true,
      fillColor: Colors.white,
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(16),
        borderSide: BorderSide.none,
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(16),
        borderSide: BorderSide(color: option.primary, width: 1.2),
      ),
    ),
  );
}
