<template>
    <div
        class="absolute inset-0 z-50 bg-white dark:bg-[#111113] flex flex-col items-center justify-center p-8 text-center animate-in fade-in duration-300">

        <!-- Círculo de Progreso / Error -->
        <div class="relative w-40 h-40 mb-8">
            <svg class="w-full h-full transform -rotate-90">
                <!-- Círculo de fondo estático -->
                <circle cx="80" cy="80" r="70" stroke="currentColor" stroke-width="8" fill="transparent"
                    class="text-gray-100 dark:text-[#1a1a1d]" />

                <!-- Círculo dinámico con cambio de color si hay error -->
                <circle cx="80" cy="80" r="70" stroke="currentColor" stroke-width="8" fill="transparent"
                    :stroke-dasharray="440" :stroke-dashoffset="smoothDashoffset" stroke-linecap="round" :class="[
                        'circle-progress',
                        isError ? 'text-red-500' : 'text-blue-600'
                    ]" />
            </svg>

            <!-- Contenido central del círculo -->
            <div class="absolute inset-0 flex flex-col items-center justify-center">
                <template v-if="isError">
                    <span class="text-4xl">⚠️</span>
                    <span class="text-[9px] uppercase tracking-[0.15em] text-red-500 font-bold mt-1">Error</span>
                </template>
                <template v-else>
                    <span class="text-3xl font-bold dark:text-white">{{ progress }}%</span>
                    <span class="text-[10px] tracking-[0.2em] text-gray-400 font-semibold mt-1">
                        {{ progress < 100 ? 'Subiendo' : 'Listo' }} </span>
                </template>
            </div>
        </div>

        <!-- Textos Dinámicos de Estado -->
        <h2 class="text-xl font-bold dark:text-white mb-2">
            {{ statusText }}
        </h2>
        <p class="text-sm text-gray-500 max-w-[240px] leading-relaxed min-h-[40px]">
            {{ descriptionText }}
        </p>

        <!-- Botonera Dinámica Interactiva -->
        <div class="flex flex-col items-center gap-4 mt-8 w-full">
            <!-- Botón de Acción Principal (Solo aparece en caso de Error) -->
            <button v-if="isError" @click="$emit('retry')"
                class="w-full max-w-[200px] bg-blue-600 hover:bg-blue-700 text-white text-xs font-bold py-3.5 px-6 rounded-2xl shadow-md shadow-blue-500/10 transition-all tracking-wider">
                🔄 Reintentar subida
            </button>

            <!-- Botón secundario para Cancelar / Volver -->
            <button @click="$emit('cancel')" :class="[
                'text-[10px] font-bold tracking-widest hover:opacity-70 transition',
                isError ? 'text-gray-100 dark:text-gray-400' : 'w-full max-w-[200px] bg-red-800 hover:bg-red-700 text-white text-xs font-bold py-3.5 px-6 rounded-2xl shadow-md shadow-blue-500/10 transition-all  tracking-wider'
            ]">
                {{ isError ? 'Cancelar transferencia' : 'Cancelar' }}
            </button>
        </div>
    </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
    progress: { type: Number, default: 0 },
    // Pasamos el error reactivo directamente al componente de progreso
    error: { type: String, default: null }
});

defineEmits(['cancel', 'retry']);

// Computed para determinar si estamos en estado de error
const isError = computed(() => !!props.error);

// Títulos dinámicos de la interfaz
const statusText = computed(() => {
    if (isError.value) return 'Transferencia pausada';
    return props.progress < 100 ? 'Enviando transferencia...' : '¡Casi listo!';
});

const smoothDashoffset = computed(() => {
  const targetProgress = Math.max(0, Math.min(props.progress, 100));
  return 440 - (440 * targetProgress) / 100;
});

// Descripciones dinámicas
const descriptionText = computed(() => {
    if (isError.value) {
        // Si Go o Axios mandaron un string de error, lo mostramos de manera amigable
        return props.error || 'Hubo un problema de conexión con el servidor.';
    }
    return 'No cierres la ventana, estamos preparando tus archivos.';
});
</script>

<style scoped>
.circle-progress {
    transition: stroke-dashoffset 0.4s cubic-bezier(0.25, 1, 0.5, 1), stroke 0.3s ease;
    will-change: stroke-dashoffset;
}
</style>