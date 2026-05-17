<template>
  <div class="flex justify-between items-center text-sm text-gray-300 bg-[#202024] px-3 py-2 rounded-xl border border-transparent hover:border-blue-500/30 transition-colors">
    <div class="flex flex-col w-[70%] truncate">
      <span class="font-medium text-gray-200">
        {{ item.type === 'folder' ? '📁' : '📄' }} {{ item.name }}
      </span>
      <span class="text-[10px] text-gray-500 uppercase tracking-tight">
        {{ item.type === 'folder' ? `${item.files.length} elementos` : displaySize }}
      </span>
    </div>
    <button @click="$emit('remove')" class="text-gray-500 hover:text-red-400 p-1 transition-colors">✕</button>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  item: {
    type: Object,
    required: true
  }
});
defineEmits(['remove']);

const formatSize = (bytes) => {
  if (!bytes || bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + " " + sizes[i];
};

const displaySize = computed(() => {
  if (props.item.type === 'file' && props.item.file) {
    // Si es un archivo, tomamos el tamaño nativo .size
    return formatSize(props.item.file.size);
  } else if (props.item.type === 'folder' && props.item.files) {
    // (Opcional) Si es carpeta, suma el tamaño de todos sus archivos internos
    const totalBytes = props.item.files.reduce((acc, f) => acc + (f.size || 0), 0);
    return formatSize(totalBytes);
  }
  return '0 B';
});

</script>