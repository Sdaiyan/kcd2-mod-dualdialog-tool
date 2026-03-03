import { useCallback } from "react";
import { Checkbox, Input, Tooltip } from "@heroui/react";
import { useTranslation } from "react-i18next";
import type { CategoryConfig } from "@/types";

// Risk level for each category ID.
// "warn" categories contain HTML markup or multi-line entries that may cause
// unexpected display artefacts in certain game interfaces.
const RISKY_CATEGORIES = new Set(["menus", "tutorials"]);

interface CategoryRowProps {
  cat: CategoryConfig;
  nameLabel: string;
  onToggle: (id: string, enabled: boolean) => void;
  onSeparatorChange: (id: string, sep: string) => void;
}

function CategoryRow({ cat, nameLabel, onToggle, onSeparatorChange }: CategoryRowProps) {
  const isRisky = RISKY_CATEGORIES.has(cat.id);
  return (
    <div className="flex items-center gap-3 py-1">
      <Checkbox
        isSelected={cat.enabled}
        onValueChange={(v) => onToggle(cat.id, v)}
        size="sm"
        className="min-w-0"
      />
      <span className="flex-1 text-sm text-white">
        {nameLabel}
        {isRisky && (
          <Tooltip content="This category contains HTML markup. Results may vary — enable only if needed.">
            <span className="ml-1 cursor-help text-yellow-400">⚠</span>
          </Tooltip>
        )}
      </span>
      <Tooltip content='Separator inserted between the two languages. Use \n for an in-game line break, or e.g. " / " for inline display.'>
        <Input
          value={cat.separator}
          onChange={(e) => onSeparatorChange(cat.id, e.target.value)}
          size="sm"
          className="w-20 shrink-0"
          aria-label="Separator"
        />
      </Tooltip>
    </div>
  );
}

interface CategorySelectProps {
  categories: CategoryConfig[];
  onChange: (categories: CategoryConfig[]) => void;
}

function CategorySelect({ categories, onChange }: CategorySelectProps) {
  const { t } = useTranslation();

  const handleToggle = useCallback(
    (id: string, enabled: boolean) => {
      onChange(categories.map((c) => (c.id === id ? { ...c, enabled } : c)));
    },
    [categories, onChange]
  );

  const handleSeparatorChange = useCallback(
    (id: string, separator: string) => {
      onChange(categories.map((c) => (c.id === id ? { ...c, separator } : c)));
    },
    [categories, onChange]
  );

  return (
    <div>
      <div className="flex items-center justify-between mb-1">
        <span className="text-sm font-medium text-gray-300">{t("LABEL_CATEGORIES")}</span>
        <span className="text-xs text-gray-500">{t("LABEL_SEPARATOR")}</span>
      </div>
      <div className="divide-y divide-gray-700">
        {categories.map((cat) => (
          <CategoryRow
            key={cat.id}
            cat={cat}
            nameLabel={t(`CATEGORY_${cat.id.toUpperCase()}`)}
            onToggle={handleToggle}
            onSeparatorChange={handleSeparatorChange}
          />
        ))}
      </div>
    </div>
  );
}

export default CategorySelect;
