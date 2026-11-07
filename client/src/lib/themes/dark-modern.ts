import darkPlus from "shiki/themes/dark-plus.mjs";

// VS Code's "Dark Modern" theme inherits the Dark+ token colors and only
// overrides the editor/UI colors. Shiki already ships Dark+, so merge the
// handful of editor colors that matter for code rendering on top of it.
export const darkModernTheme = {
  ...darkPlus,
  name: "dark-modern",
  displayName: "Dark Modern",
  colors: {
    ...darkPlus.colors,
    "editor.background": "#1F1F1F",
    "editor.foreground": "#CCCCCC",
  },
};
