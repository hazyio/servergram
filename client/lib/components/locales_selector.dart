import 'package:flutter/material.dart';

import '../gen/i18n/strings.g.dart';

class LanguageSelector extends StatelessWidget {
  const LanguageSelector({super.key});

  @override
  Widget build(BuildContext context) {
    final currentLocale = LocaleSettings.currentLocale;

    return DropdownButtonFormField<AppLocale>(
      key: ValueKey(currentLocale),
      initialValue: currentLocale,
      decoration: InputDecoration(labelText: context.t.labels.language),
      isExpanded: true,
      items: AppLocale.values
          .map(
            (locale) => DropdownMenuItem(
              value: locale,
              child: Text(_displayName(locale)),
            ),
          )
          .toList(),
      onChanged: (locale) async {
        if (locale != null) {
          await LocaleSettings.setLocale(locale);
        }
      },
    );
  }

  String _displayName(AppLocale locale) {
    switch (locale) {
      case AppLocale.en:
        return 'English';
    }
  }
}
