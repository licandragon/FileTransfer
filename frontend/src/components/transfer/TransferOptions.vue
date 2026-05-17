<template>
  <Transition name="panel">
    <div v-if="show"
      class="absolute z-[80] inset-0 rounded-[2.5rem] md:inset-auto md:top-0 md:left-[105%] md:w-72 md:h-full bg-white dark:bg-[#111113] p-6 shadow-2xl border border-black/5 dark:border-white/5 flex flex-col">

      <div class="flex justify-between items-center mb-6">
        <h3 class="text-[10px] font-bold text-gray-400 uppercase tracking-widest">Ajustes de transferencia</h3>
        <button @click="show = false"
          class="text-gray-400 hover:text-black dark:hover:text-white transition-colors">✕</button>
      </div>

      <div class="space-y-6 overflow-y-auto custom-scrollbar pr-2 flex-1">

        <div class="space-y-3">
          <label class="text-[11px] font-bold text-gray-500 uppercase">Enviar como</label>
          <div class="flex bg-gray-100 dark:bg-[#1a1a1d] p-1 rounded-2xl">
            <button @click="mode = 'email'"
              :class="mode === 'email' ? 'bg-white dark:bg-[#2a2a2e] shadow-sm text-blue-600' : 'text-gray-500'"
              class="flex-1 py-2 text-[10px] rounded-xl transition-all font-bold uppercase">Email</button>
            <button @click="mode = 'link'"
              :class="mode === 'link' ? 'bg-white dark:bg-[#2a2a2e] shadow-sm text-blue-600' : 'text-gray-500'"
              class="flex-1 py-2 text-[10px] rounded-xl transition-all font-bold uppercase">Enlace</button>
          </div>
        </div>

        <hr class="border-gray-100 dark:border-white/5" />

        <div class="space-y-4">
          <label class="text-[11px] font-bold text-gray-500 uppercase">Seguridad</label>

          <div class="flex items-center justify-between">
            <span class="text-xs dark:text-gray-300">Proteger con contraseña</span>
            <input type="checkbox" v-model="extraOptions.hasPassword" class="accent-blue-600" />
          </div>
          <input v-if="extraOptions.hasPassword" type="password" placeholder="Escribe la clave..."
            class="w-full bg-gray-100 dark:bg-[#1a1a1d] border-none rounded-xl p-3 text-xs outline-none focus:ring-1 ring-blue-500" />

          <div class="flex items-center justify-between">
            <span class="text-xs dark:text-gray-300">Avisar al descargar</span>
            <label class="relative inline-flex items-center cursor-pointer">
              <input type="checkbox" v-model="extraOptions.notifyOnDownload" class="sr-only peer">
              <div
                class="w-8 h-4 bg-gray-200 peer-focus:outline-none rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-3 after:w-3 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600">
              </div>
            </label>
          </div>
        </div>

        <div class="space-y-4 pt-2">
          <label class="text-[11px] font-bold text-gray-500 uppercase">Gestión</label>
          <div class="flex items-center justify-between">
            <span class="text-xs dark:text-gray-300">Deshabilitar vista previa</span>
            <input type="checkbox" v-model="extraOptions.disablePreview" class="accent-blue-600" />
          </div>
        </div>

      </div>

      <button @click="show = false"
        class="mt-6 w-full bg-blue-600 text-white py-4 rounded-2xl font-bold hover:bg-blue-700 transition-all">
        Guardar cambios
      </button>
    </div>
  </Transition>
</template>

<script setup>
import { reactive } from 'vue';

const mode = defineModel('mode', { default: 'email' });
const show = defineModel('show', { default: false });
// Estado local para las nuevas opciones
const extraOptions = reactive({
  hasPassword: false,
  password: '',
  notifyOnDownload: true,
  disablePreview: false
});
</script>