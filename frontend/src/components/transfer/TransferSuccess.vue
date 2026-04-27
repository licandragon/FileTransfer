<template>
    <div
        class="absolute inset-0 z-50 bg-white dark:bg-[#111113] flex flex-col items-center justify-center p-8 text-center animate-in zoom-in duration-500">

        <div
            class="w-20 h-20 bg-green-100 dark:bg-green-900/20 text-green-500 rounded-full flex items-center justify-center mb-6 text-3xl shadow-lg shadow-green-500/10">
            ✓
        </div>

        <h2 class="text-2xl font-bold dark:text-white mb-2">
            {{ mode === 'email' ? '¡Email enviado!' : '¡Listo para compartir!' }}
        </h2>

        <p class="text-sm text-gray-500 mb-8 max-w-[240px]">
            {{ mode === 'email'
                ? 'Tus archivos se han enviado correctamente. Revisa tu bandeja de entrada.'
                : 'Cualquier persona con el enlace podrá descargar tus archivos.'
            }}
        </p>

        <div v-if="mode === 'link'" class="w-full relative group">
            <input readonly :value="link"
                class="w-full bg-gray-50 dark:bg-[#1a1a1d] border border-gray-100 dark:border-white/5 rounded-2xl py-4 px-4 text-xs text-gray-600 dark:text-gray-300 outline-none pr-24" />
            <button @click="copyLink" :class="copied ? 'bg-green-500' : 'bg-blue-600 hover:bg-blue-700'"
                class="absolute right-2 top-2 bottom-2 px-4 rounded-xl text-[10px] font-bold text-white uppercase transition-all duration-300">
                {{ copied ? '¡Copiado!' : 'Copiar' }}
            </button>
        </div>

        <div v-else
            class="w-full p-4 bg-blue-50 dark:bg-blue-900/10 rounded-2xl border border-blue-100 dark:border-blue-900/10">
            <span class="text-[10px] text-blue-500 font-bold uppercase block mb-1">Enviado a</span>
            <div class="text-xs text-gray-600 dark:text-gray-300 truncate">
                {{ emails?.join(', ') || 'Destinatarios' }}
            </div>
        </div>

        <button @click="$emit('reset')"
            class="mt-12 text-[10px] font-bold text-gray-400 uppercase tracking-[0.2em] hover:text-blue-500 transition-colors">
            Enviar otro archivo
        </button>
    </div>
</template>

<script setup>
import { ref } from 'vue';

const props = defineProps({
    mode: String,   // 'email' | 'link'
    link: String,   // URL generada
    emails: Array   // Lista de destinatarios
});

const emit = defineEmits(['reset']);
const copied = ref(false);

const copyLink = async () => {
    try {
        await navigator.clipboard.writeText(props.link);
        copied.value = true;
        setTimeout(() => (copied.value = false), 2000);
    } catch (err) {
        console.error("Error al copiar enlace", err);
    }
};
</script>