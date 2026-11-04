type LocalStorageKeys = "theme" | "file-explorer-settings";
const LocalStorageKeys = {
  theme: "theme",
  "file-explorer-settings": "file-explorer-settings",
};
export function useLocalStorage() {
  return {
    save: <T>(key: LocalStorageKeys, value: T) => {
      if (typeof window === "undefined") {
        console.error("You cannot use localStorage in backend");
        return;
      }
      const keyValue = LocalStorageKeys[key];

      let valueToSave = JSON.stringify(value);
      if (typeof value === "string") {
        valueToSave = value;
      }
      localStorage.setItem(keyValue, valueToSave);
    },
    get: <T>(key: LocalStorageKeys, d: T) => {
      if (typeof window === "undefined") {
        console.error("You cannot use localStorage in backend");
        return d;
      }
      const keyValue = LocalStorageKeys[key];
      const i = localStorage.getItem(keyValue);
      if (i === null) return d;
      if (typeof d === "string") {
        // no need to parse
        return i as T;
      }
      try {
        return JSON.parse(i) as T;
      } catch (e) {
        console.error("Failed to parse localStorage value", i, e);
        return i as T;
      }
    },
    remove: (key: string) => {
      if (typeof window === "undefined") {
        console.error("You cannot use localStorage in backend");
        return;
      }
      localStorage.removeItem(key);
    },
    exists: (key: LocalStorageKeys): boolean => {
      return useLocalStorage().get(key, false);
    },
  };
}
