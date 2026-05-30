import 'package:flutter_test/flutter_test.dart';
import 'package:manba_alert/providers/app_state.dart';
import 'package:manba_alert/services/settings_service.dart';
import 'package:manba_alert/theme/color_schemes.dart';

class _FakeSettingsService extends SettingsService {
  String? storedThemeId;

  @override
  Future<String?> loadThemeId() async => storedThemeId;

  @override
  Future<void> saveThemeId(String themeId) async {
    storedThemeId = themeId;
  }
}

void main() {
  test('defaults to purple-gold theme', () {
    final state = AppState(
      settingsService: _FakeSettingsService(),
      initialThemeId: null,
    );
    expect(state.currentTheme.id, purpleGoldTheme.id);
  });

  test('selectTheme updates current theme and persists it', () async {
    final settings = _FakeSettingsService();
    final state = AppState(settingsService: settings, initialThemeId: null);

    await state.selectTheme(tealCreamTheme.id);

    expect(state.currentTheme.id, tealCreamTheme.id);
    expect(settings.storedThemeId, tealCreamTheme.id);
  });
}
