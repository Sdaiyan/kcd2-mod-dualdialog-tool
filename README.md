# Dual Dialog

**Dual Dialog** is a mod tool for _Kingdom Come: Deliverance 2_ that generates a mod displaying in-game text in two languages simultaneously. Whether you want to play in your native language while reinforcing a second language, or simply prefer a bilingual experience, this tool has you covered.

_Kingdom Come: Deliverance 2_ (KCD2) is a medieval simulation game renowned for its rich dialogue and exceptional voice acting. This mod leverages that strength, turning your playthrough into an opportunity to improve your language skills while diving deep into the story.

This project is available on [GitHub](https://github.com/SDxBacon/kcd2-mod-dualdialog-tool).

> 📖 [中文说明请点这里](./README_zh.md)

## Features

- **Dual-Language Display**: In-game text appears in both your primary language and a paired language side by side.
- **Wide Content Coverage** — 8 bilingual categories supported:
  - ✅ **Dialog** — Overhead, speech, and dialog sequence subtitles
  - ✅ **Quest Journal** — Quest names and descriptions
  - ✅ **Alchemy & Items** — Item names, descriptions, and recipe text
  - ✅ **Skills & Buffs** — Skill names and perk descriptions
  - ✅ **Achievements** — Achievement names and descriptions
  - ✅ **Estate Minigame** — Minigame-related UI text
  - ⚠️ **Menus & Codex** — Menu labels and encyclopedia entries (may have display artefacts)
  - ⚠️ **Tutorials** — Tutorial text with UI markup (may have display artefacts)
- **Per-Category Separator**: Choose `\n` (in-game line break) or ` / ` (inline slash) independently for each category.
- **Community Patch Support**: Load one or more community-fix XML files (e.g. ChineseFix) for either language before generating — corrected entries take priority.
- **No Full Game Install Required**: Only the `Localization/` folder containing the language `.pak` files is needed.

<!-- TODO: screenshot — tool main UI -->

## Download

Get the latest release from the [GitHub Releases page](https://github.com/SDxBacon/kcd2-mod-dualdialog-tool/releases).

## Usage

### Basic

1. **Launch the Tool** — Run `Dual Dialog Tool.exe`.
2. **Select Game Folder** — Choose the folder that contains the `Localization/` subfolder with the language `.pak` files. This does not need to be the full game directory.
3. **Choose Languages** — Select your primary game language and the paired language you want shown alongside.
4. **Export** — Click **Export**, choose an output folder. The tool generates `Dual Dialog.zip`.
5. **Install the Mod** — Extract `Dual Dialog.zip` into your KCD2 `Mods` folder (e.g. `…\KingdomComeDeliverance2\Mods`). Create the folder if it does not exist.
6. **Launch the Game** — Enjoy bilingual text in-game!

<!-- TODO: GIF — basic usage workflow -->

### Advanced Settings

Click **▸ Advanced Settings** in the tool to access:

- **Content Categories** — Toggle each of the 8 categories on/off and set a custom separator per category.
- **Main / Secondary Language Patches** — Add community patch XML files for either language. The tool previews how many entries each patch contains. Patched entries override the original pak data before merging.

<!-- TODO: screenshot — Advanced Settings panel -->

## Suggested Language Combinations

KCD2 language packs bundle only the fonts needed for that script. Loading a script outside the active pack results in missing characters.

| Primary Language | Paired Language | Result |
|---|---|---|
| Latin-based (e.g. English) | Latin-based (e.g. French) | ✅ Works fine |
| Asian (e.g. Chinese) | Latin-based (e.g. English) | ✅ Works fine |
| Asian (e.g. Chinese) | Russian | ⚠ Some Cyrillic characters missing |
| Latin-based | Asian (e.g. Chinese) | ❌ Most Asian characters missing |
| Simplified Chinese | Japanese | ❌ Fonts not shared between Asian packs |

## Feedback & Support

If you encounter any issues or have suggestions, feel free to open an Issue on GitHub.
All feedback is welcome!

## Enjoy your bilingual adventure with Dual Dialog! 🎮🌏

## Disclaimer

- The author only understands Chinese and English — if there are issues with other language combinations, please report them.
- Categories marked ⚠️ (Menus & Codex, Tutorials) contain HTML-like markup. Results may vary; enable them only if needed.
- This tool is provided as-is. Always back up your `Mods` folder before installing.
