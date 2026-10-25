import 'package:client/gen/i18n/strings.g.dart';
import 'package:flutter/material.dart';
import 'package:client/utils/theme_controller.dart';
import 'package:dynamic_color/dynamic_color.dart';
import 'package:client/theme/app_theme.dart';
import 'package:client/pages/welcome/page.dart';
import 'package:flutter_localizations/flutter_localizations.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  LocaleSettings.useDeviceLocale();
  await loadTheme();
  runApp(TranslationProvider(child: MyApp(initialRoute: '/welcome')));
}

class MyApp extends StatelessWidget {
  const MyApp({super.key, required this.initialRoute});
  final String initialRoute;

  @override
  Widget build(BuildContext context) {
    return DynamicColorBuilder(
      builder: (ColorScheme? lightDynamic, ColorScheme? darkDynamic) {
        return ValueListenableBuilder<ThemeMode>(
          valueListenable: themeMode,
          builder: (context, themeMode, child) => MaterialApp(
            theme: AppTheme.light(dynamicScheme: lightDynamic),
            darkTheme: AppTheme.dark(dynamicScheme: darkDynamic),
            themeMode: themeMode,
            initialRoute: initialRoute,
            locale: TranslationProvider.of(context).flutterLocale,
            supportedLocales: AppLocaleUtils.supportedLocales,
            localizationsDelegates: GlobalMaterialLocalizations.delegates,
            title: context.t.titles.appTitle,
            routes: {'/welcome': (context) => const WelcomePage()},
          ),
        );
      },
    );
  }
}
