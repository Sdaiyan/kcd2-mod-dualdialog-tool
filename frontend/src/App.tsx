import { useState, useCallback, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { Card, CardBody } from "@heroui/react";
import { Language } from "@/constants/languages";
import Navbar from "@/components/Navbar";
import ExportButton from "@/components/ExportButton";
import LanguageSelect from "@/components/LanguageSelect";
import KingdomComeFolderPicker from "@/components/KingdomComeFolderPicker";
import CategorySelect from "@/components/CategorySelect";
import PatchFilePicker from "@/components/PatchFilePicker";
// import hooks & utilities
import isNil from "lodash/isNil";
import isEmpty from "lodash/isEmpty";
import useExport from "@/hooks/useExport";
import type { CategoryConfig } from "@/types";
import { GetDefaultCategories } from "../wailsjs/go/main/App";
// import css
import "./App.css";
import "./i18n";

function App() {
  const { t } = useTranslation();

  const [isExporting, setIsExporting] = useState(false);
  const [showAdvanced, setShowAdvanced] = useState(false);

  const [folder, setFolder] = useState("");
  const [isFolderError, setIsFolderError] = useState(false);
  const [mainLanguage, setMainLanguage] = useState<Language | undefined>();
  const [subLanguage, setSubLanguage] = useState<Language | undefined>();

  const [categories, setCategories] = useState<CategoryConfig[]>([]);
  const [mainPatchPaths, setMainPatchPaths] = useState<string[]>([]);
  const [subPatchPaths, setSubPatchPaths] = useState<string[]>([]);

  // Load default category configs from the Go backend on startup.
  useEffect(() => {
    GetDefaultCategories().then(setCategories).catch(console.error);
  }, []);

  // disabled main language keys
  const disabledSubLanguageKeys = isNil(mainLanguage) ? [] : [mainLanguage];

  // check if export button is disabled
  const isExportButtonDisabled =
    isNil(mainLanguage) || isNil(subLanguage) || isEmpty(folder);

  const startExport = useExport();

  const handleExportButtonPress = useCallback(() => {
    if (!mainLanguage || !subLanguage) {
      return;
    }

    setIsExporting(true);

    startExport(mainLanguage, subLanguage, categories, mainPatchPaths, subPatchPaths).then(() => {
      setIsExporting(false);
    });
  }, [startExport, mainLanguage, subLanguage, categories, mainPatchPaths, subPatchPaths]);

  const handleMainLanguageSelect = useCallback((language: Language) => {
    setMainLanguage(language);
    setSubLanguage(undefined); // reset sub language
  }, []);

  const handleSubLanguageSelect = useCallback((language: Language) => {
    setSubLanguage(language);
  }, []);

  return (
    <div id="App" className="bg-gray-800 h-screen overflow-y-auto px-[100px] pb-11 pt-4">
      {/* Navbar */}
      <Navbar />

      <KingdomComeFolderPicker
        value={folder}
        isError={isFolderError}
        onSelect={(path) => {
          setFolder(path);
          setIsFolderError(false);
        }}
        onSelectError={() => setIsFolderError(true)}
      />

      <div className="flex w-full flex-wrap md:flex-nowrap gap-4 items-center">
        <LanguageSelect
          label={t("LABEL_GAME_LANGUAGE")}
          value={mainLanguage}
          disabledKeys={mainLanguage ? [mainLanguage] : undefined}
          onSelect={handleMainLanguageSelect}
        />
        <LanguageSelect
          label={t("LABEL_PAIRED_LANGUAGE")}
          value={subLanguage}
          onSelect={handleSubLanguageSelect}
          disabledKeys={disabledSubLanguageKeys}
          hideAsianLanguages
        />
      </div>

      <div className="w-full flex justify-center items-center mt-3">
        <Card>
          <CardBody className="text-[13px]">
            {t("NOTE_1")}
            <ul>
              <li>✅ {t("NOTE_COMB_1")}</li>
              <li>✅ {t("NOTE_COMB_2")}</li>
            </ul>
            <br />
            {t("NOTE_2")}
          </CardBody>
        </Card>
      </div>

      {/* Advanced Settings toggle */}
      <div className="mt-3">
        <button
          className="text-sm text-blue-400 hover:text-blue-300 underline"
          onClick={() => setShowAdvanced((v) => !v)}
        >
          {showAdvanced ? "▾" : "▸"} {t("SECTION_ADVANCED")}
        </button>

        {showAdvanced && (
          <div className="mt-2 bg-gray-700 rounded-lg p-4 space-y-4">
            {/* Category selection */}
            {categories.length > 0 && (
              <CategorySelect categories={categories} onChange={setCategories} />
            )}

            <hr className="border-gray-600" />

            {/* Patch file pickers (main + sub language) */}
            <PatchFilePicker
              label={t("LABEL_PATCHES_MAIN")}
              patchPaths={mainPatchPaths}
              onChange={setMainPatchPaths}
            />
            <PatchFilePicker
              label={t("LABEL_PATCHES_SUB")}
              patchPaths={subPatchPaths}
              onChange={setSubPatchPaths}
            />
          </div>
        )}
      </div>

      <ExportButton
        isLoading={isExporting}
        onPress={handleExportButtonPress}
        isDisabled={isExportButtonDisabled}
      >
        {!isExporting && t("BUTTON_EXPORT")}
      </ExportButton>
    </div>
  );
}

export default App;

