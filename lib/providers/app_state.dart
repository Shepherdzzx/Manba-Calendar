import 'package:flutter/foundation.dart';

import '../services/settings_service.dart';
import '../theme/color_schemes.dart';

class AppState extends ChangeNotifier {
  AppState({required SettingsService settingsService, String? initialThemeId})
    : _settingsService = settingsService,
      _currentTheme = _resolveTheme(initialThemeId);

  final SettingsService _settingsService;
  AppThemeOption _currentTheme;

  AppThemeOption get currentTheme => _currentTheme;
  List<AppThemeOption> get availableThemes => themeOptions;

  Future<void> selectTheme(String themeId) async {
    final nextTheme = _resolveTheme(themeId);
    if (nextTheme.id == _currentTheme.id) {
      return;
    }
    _currentTheme = nextTheme;
    notifyListeners();
    await _settingsService.saveThemeId(themeId);
  }

  static AppThemeOption _resolveTheme(String? themeId) {
    return themeOptions.firstWhere(
      (option) => option.id == themeId,
      orElse: () => purpleGoldTheme,
    );
  }
}
