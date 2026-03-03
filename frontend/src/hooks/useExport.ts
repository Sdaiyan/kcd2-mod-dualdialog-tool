import { addToast } from "@heroui/react";
import { useCallback } from "react";
import { useTranslation } from "react-i18next";
import { Language } from "../constants/languages";
import type { CategoryConfig } from "@/types";
import { CreateModZip } from "../../wailsjs/go/main/App";

function useExport() {
  const { t } = useTranslation();

  return useCallback(
    (
      main: Language,
      sub: Language,
      categories: CategoryConfig[],
      mainPatchPaths: string[],
      subPatchPaths: string[]
    ) => {
      return CreateModZip(main, sub, categories, mainPatchPaths, subPatchPaths)
        .then(() => {
          addToast({
            color: "success",
            title: t("TOAST_EXPORT_SUCCESS"),
          });
        })
        .catch(() => {
          addToast({
            title: t("TOAST_EXPORT_FAILED"),
            color: "danger",
          });
        });
    },
    []
  );
}

export default useExport;
