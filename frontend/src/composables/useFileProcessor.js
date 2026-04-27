import { ref } from "vue";

export function useFileProcessor() {
  const files = ref([]); // Archivos nativos
  const items = ref([]); // Items para la UI

  const formatSize = (bytes) => {
    if (bytes === 0) return "0 B";
    const k = 1024;
    const sizes = ["B", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + " " + sizes[i];
  };

  const processFiles = (fileList) => {
    const map = {};
    fileList.forEach((file) => {
      const path = file.webkitRelativePath || file.name;
      const parts = path.split("/");

      if (parts.length > 1) {
        const folderName = parts[0];
        if (!map[folderName]) {
          map[folderName] = { type: "folder", name: folderName, files: [] };
        }
        map[folderName].files.push(file);
      } else {
        map[file.name] = { type: "file", name: file.name, file };
      }
    });
    items.value = Object.values(map);
  };

  const addValidFiles = (newFiles) => {
    const filtered = newFiles.filter((file) => {
      const rootName = file.webkitRelativePath
        ? file.webkitRelativePath.split("/")[0]
        : file.name;
      const type = file.webkitRelativePath ? "folder" : "file";
      return !items.value.some(
        (item) => item.name === rootName && item.type === type,
      );
    });

    if (filtered.length > 0) {
      files.value.push(...filtered);
      processFiles(files.value);
    }
  };

  const removeItem = (index) => {
    const item = items.value[index];
    if (item.type === "folder") {
      files.value = files.value.filter(
        (f) => !f.webkitRelativePath.startsWith(item.name + "/"),
      );
    } else {
      files.value = files.value.filter((f) => f !== item.file);
    }
    processFiles(files.value);
  };

  return { items, files, addValidFiles, removeItem, formatSize };
}
