import { useCallback, useEffect, useState } from "react";
import { Button } from "@heroui/react";
import { useTranslation } from "react-i18next";
import { SelectPatchFile, PreviewPatchFile } from "../../wailsjs/go/main/App";

interface PatchEntry {
  path: string;
  count: number | null; // null while loading
}

interface PatchFilePickerProps {
  label: string;
  patchPaths: string[];
  onChange: (paths: string[]) => void;
}

function PatchFilePicker({ label, patchPaths, onChange }: PatchFilePickerProps) {
  const { t } = useTranslation();
  const [entries, setEntries] = useState<PatchEntry[]>([]);

  // Keep entries in sync when parent resets paths (e.g. on clear).
  useEffect(() => {
    setEntries((prev) => {
      const prevPaths = prev.map((e) => e.path);
      if (JSON.stringify(prevPaths) === JSON.stringify(patchPaths)) return prev;
      return patchPaths.map((path) => {
        const existing = prev.find((e) => e.path === path);
        return existing ?? { path, count: null };
      });
    });
  }, [patchPaths]);

  // Fetch entry counts for patches that don't have one yet.
  useEffect(() => {
    entries.forEach((entry) => {
      if (entry.count === null) {
        PreviewPatchFile(entry.path)
          .then((count) => {
            setEntries((prev) =>
              prev.map((e) => (e.path === entry.path ? { ...e, count } : e))
            );
          })
          .catch(() => {
            setEntries((prev) =>
              prev.map((e) => (e.path === entry.path ? { ...e, count: -1 } : e))
            );
          });
      }
    });
  }, [entries]);

  const handleAdd = useCallback(async () => {
    const path = await SelectPatchFile();
    if (!path || patchPaths.includes(path)) return;
    const newPaths = [...patchPaths, path];
    onChange(newPaths);
    setEntries((prev) => [...prev, { path, count: null }]);
  }, [patchPaths, onChange]);

  const handleRemove = useCallback(
    (path: string) => {
      onChange(patchPaths.filter((p) => p !== path));
      setEntries((prev) => prev.filter((e) => e.path !== path));
    },
    [patchPaths, onChange]
  );

  return (
    <div className="mb-2">
      <div className="flex items-center justify-between mb-1">
        <span className="text-sm font-medium text-gray-300">{label}</span>
        <Button size="sm" variant="bordered" onPress={handleAdd} className="text-xs">
          + {t("BUTTON_ADD_PATCH")}
        </Button>
      </div>
      {entries.length === 0 ? (
        <p className="text-xs text-gray-500 italic">{t("HINT_NO_PATCHES")}</p>
      ) : (
        <ul className="space-y-1">
          {entries.map((entry) => {
            const fileName = entry.path.replace(/\\/g, "/").split("/").pop();
            const countText =
              entry.count === null
                ? "…"
                : entry.count === -1
                ? t("HINT_PATCH_ERROR")
                : t("HINT_PATCH_ENTRIES", { count: entry.count });
            return (
              <li
                key={entry.path}
                className="flex items-center justify-between bg-gray-700 rounded px-2 py-1"
              >
                <span className="text-xs text-white truncate flex-1 mr-2" title={entry.path}>
                  {fileName}
                </span>
                <span className="text-xs text-gray-400 mr-2 shrink-0">{countText}</span>
                <button
                  onClick={() => handleRemove(entry.path)}
                  className="text-gray-400 hover:text-red-400 text-xs shrink-0"
                  aria-label="Remove patch"
                >
                  ✕
                </button>
              </li>
            );
          })}
        </ul>
      )}
    </div>
  );
}

export default PatchFilePicker;
