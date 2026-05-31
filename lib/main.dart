import 'package:flutter/widgets.dart';

import 'app.dart';
import 'providers/app_state.dart';
import 'services/settings_service.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  final settingsService = SettingsService();
  final initialThemeId = await settingsService.loadThemeId();
  final appState = AppState(
    settingsService: settingsService,
    initialThemeId: initialThemeId,
  );
  runApp(ManbaCalendarApp(appState: appState));
}
