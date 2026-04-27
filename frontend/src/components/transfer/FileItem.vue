<template>
  <div class="flex justify-between items-center text-sm text-gray-300 bg-[#202024] px-3 py-2 rounded-xl border border-transparent hover:border-blue-500/30 transition-colors">
    <div class="flex flex-col w-[70%] truncate">
      <span class="font-medium text-gray-200">
        {{ item.type === 'folder' ? '📁' : '📄' }} {{ item.name }}
      </span>
      <span class="text-[10px] text-gray-500 uppercase tracking-tight">
        {{ item.type === 'folder' ? `${item.files.length} elementos` : sizeLabel }}
      </span>
    </div>
    <button @click="$emit('remove')" class="text-gray-500 hover:text-red-400 p-1 transition-colors">✕</button>
  </div>
</template>

<script setup>
const props = defineProps(['item', 'sizeLabel']);
defineEmits(['remove']);
</script>