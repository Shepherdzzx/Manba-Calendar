import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import 'providers/app_state.dart';
import 'screens/home_screen.dart';
import 'theme/app_themes.dart';

class ManbaCalendarApp extends StatelessWidget {
  const ManbaCalendarApp({super.key, required this.appState});

  final AppState appState;

  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider.value(
      value: appState,
      child: Consumer<AppState>(
        builder: (context, state, _) {
          return MaterialApp(
            title: 'Manba Calendar',
            debugShowCheckedModeBanner: false,
            theme: buildTheme(state.currentTheme),
            home: const HomeScreen(),
          );
        },
      ),
    );
  }
}
