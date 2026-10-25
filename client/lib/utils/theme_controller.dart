import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';

final themeMode = ValueNotifier<ThemeMode>(ThemeMode.system);

Future<void> loadTheme() async {
  final prefs = await SharedPreferences.getInstance();
  final saved = prefs.getString('themeMode');
  themeMode.value = ThemeMode.values.firstWhere(
    (m) => m.name == saved,
    orElse: () => ThemeMode.system,
  );
}

Future<void> setTheme(ThemeMode mode) async {
  themeMode.value = mode;
  final prefs = await SharedPreferences.getInstance();
  await prefs.setString('themeMode', mode.name);
}
