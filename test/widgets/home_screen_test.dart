import 'package:flutter_test/flutter_test.dart';
import 'package:manba_alert/providers/app_state.dart';
import 'package:manba_alert/screens/home_screen.dart';
import 'package:manba_alert/services/settings_service.dart';
import 'package:provider/provider.dart';
import 'package:flutter/material.dart';

class _FakeSettingsService extends SettingsService {
  @override
  Future<String?> loadThemeId() async => null;

  @override
  Future<void> saveThemeId(String themeId) async {}
}

void main() {
  testWidgets('home screen shows calendar and voice entry', (tester) async {
    final appState = AppState(
      settingsService: _FakeSettingsService(),
      initialThemeId: null,
    );

    await tester.pumpWidget(
      ChangeNotifierProvider.value(
        value: appState,
        child: MaterialApp(home: const HomeScreen()),
      ),
    );

    await tester.pumpAndSettle();

    expect(find.byType(GestureDetector), findsWidgets);
    expect(find.byType(FloatingActionButton), findsOneWidget);
    expect(find.text('新增事件'), findsOneWidget);
    expect(find.byIcon(Icons.palette_outlined), findsOneWidget);
  });
}
