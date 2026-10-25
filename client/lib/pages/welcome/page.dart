import 'package:client/components/locales_selector.dart';
import 'package:client/components/padded_card.dart';
import 'package:client/components/theme_selector.dart';
import 'package:client/theme/app_dimens.dart';
import 'package:flutter/material.dart';

import '../../gen/i18n/strings.g.dart';

class WelcomePage extends StatefulWidget {
  const WelcomePage({super.key});

  @override
  State<WelcomePage> createState() => _WelcomePageState();
}

class _WelcomePageState extends State<WelcomePage> {
  ThemeMode selectedThemeMode = ThemeMode.system;
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(t.Welcome)),
      body: Center(
        child: SizedBox(
          width: MediaQuery.of(context).size.height / 100 * 60,
          height: MediaQuery.of(context).size.height / 100 * 32,
          child: PaddedCard(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(AppRadius.small),
            ),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const LanguageSelector(),
                const SizedBox(height: AppSpacing.marginTablet),
                ThemeSelector(),
                Expanded(child: Container()),
                Row(
                  mainAxisAlignment: MainAxisAlignment.end,
                  children: [
                    FilledButton(
                      onPressed: () {
                        // Handle button press
                      },
                      child: Text(t.labels.kContinue),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
