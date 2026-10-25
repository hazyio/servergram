import 'package:client/gen/i18n/strings.g.dart';
import 'package:client/theme/app_dimens.dart';
import 'package:flutter/material.dart';
import 'package:client/utils/theme_controller.dart';

class ThemeSelector extends StatefulWidget {
  const ThemeSelector({super.key});

  @override
  State<ThemeSelector> createState() => _ThemeSelectorState();
}

class _ThemeSelectorState extends State<ThemeSelector> {
  ThemeMode selectedThemeMode = themeMode.value;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(t.labels.theme),
        const SizedBox(height: AppSpacing.baseline),
        SizedBox(
          width: double.infinity,
          child: SegmentedButton<ThemeMode>(
            segments: [
              ButtonSegment(
                value: ThemeMode.system,
                label: Text(t.labels.System),
              ),
              ButtonSegment(
                value: ThemeMode.light,
                label: Text(t.labels.Light),
              ),
              ButtonSegment(value: ThemeMode.dark, label: Text(t.labels.Dark)),
            ],
            selected: {selectedThemeMode},
            onSelectionChanged: (s) => setState(() {
              selectedThemeMode = s.first;
              setTheme(selectedThemeMode);
            }),
          ),
        ),
      ],
    );
  }
}
