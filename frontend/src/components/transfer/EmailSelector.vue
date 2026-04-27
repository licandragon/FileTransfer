<template>
    <div class="border-b border-gray-200 dark:border-[#2a2a2e] py-2">
        <div class="flex justify-between items-center mb-2">
            <span class="text-[10px] font-bold text-gray-400 uppercase tracking-widest">Enviar a</span>
            <span class="text-[10px] text-gray-400">{{ modelValue.length }}/{{ limit }}</span>
        </div>

        <div class="flex flex-wrap gap-2 mb-2">
            <div v-for="(email, index) in modelValue" :key="index"
                class="bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-300 text-xs px-3 py-1 rounded-full flex items-center gap-2 border border-blue-100 dark:border-blue-800/50">
                {{ email }}
                <button @click="removeEmail(index)" class="hover:text-red-500">✕</button>
            </div>
        </div>

        <input v-model="input" @keydown.enter.prevent="addEmail" @keydown.tab.prevent="addEmail"
            @keydown.delete="handleBackspace" :disabled="modelValue.length >= limit"
            :placeholder="modelValue.length >= limit ? 'Límite alcanzado' : 'Introduce un email...'"
            class="w-full bg-transparent outline-none text-sm text-gray-700 dark:text-gray-200 placeholder:text-gray-500" />
    </div>
</template>

<script setup>
import { ref } from 'vue';
const props = defineProps(['modelValue', 'limit']);
const emit = defineEmits(['update:modelValue']);

const input = ref('');

const addEmail = () => {
    const val = input.value.trim().toLowerCase();
    if (!val || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(val) || props.modelValue.includes(val)) return;
    emit('update:modelValue', [...props.modelValue, val]);
    input.value = '';
};

const removeEmail = (index) => {
    const newList = [...props.modelValue];
    newList.splice(index, 1);
    emit('update:modelValue', newList);
};

const handleBackspace = () => {
    if (input.value === '' && props.modelValue.length > 0) {
        removeEmail(props.modelValue.length - 1);
    }
};
</script>