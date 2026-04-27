<template>
    <div
        class="absolute inset-0 z-50 bg-white dark:bg-[#111113] flex flex-col items-center justify-center p-8 text-center animate-in fade-in duration-300">

        <div class="relative w-40 h-40 mb-8">
            <svg class="w-full h-full transform -rotate-90">
                <circle cx="80" cy="80" r="70" stroke="currentColor" stroke-width="8" fill="transparent"
                    class="text-gray-100 dark:text-[#1a1a1d]" />
                <circle cx="80" cy="80" r="70" stroke="currentColor" stroke-width="8" fill="transparent"
                    :stroke-dasharray="440" :stroke-dashoffset="440 - (440 * progress) / 100" stroke-linecap="round"
                    class="text-blue-600 transition-all duration-300 ease-out" />
            </svg>

            <div class="absolute inset-0 flex flex-col items-center justify-center">
                <span class="text-3xl font-bold dark:text-white">{{ progress }}%</span>
                <span class="text-[10px] uppercase tracking-[0.2em] text-gray-400 font-semibold mt-1">Subiendo</span>
            </div>
        </div>

        <h2 class="text-xl font-bold dark:text-white mb-2">
            {{ progress < 100 ? 'Enviando transferencia...' : '¡Casi listo!' }} </h2>
                <p class="text-sm text-gray-500 max-w-[220px] leading-relaxed">
                    No cierres la ventana, estamos preparando tus archivos.
                </p>

                <button @click="$emit('cancel')"
                    class="mt-10 text-[10px] font-bold text-red-500 uppercase tracking-widest hover:opacity-70 transition">
                    Cancelar
                </button>
    </div>
</template>

<script setup>
defineProps({
    progress: { type: Number, default: 0 }
});
defineEmits(['cancel']);
</script>