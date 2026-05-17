<template>
  <div
    class="relative w-full max-w-[420px] mx-auto bg-white dark:bg-[#111113] rounded-[2.5rem] shadow-[0_20px_60px_rgba(0,0,0,.25)] p-5 border border-black/5 dark:border-white/5 min-h-[550px] transition-all"
    :class="showOptions ? 'overflow-visible' : 'overflow-hidden'">
    <!-- Toast para mostrar mensajes de error -->
    <Transition name="toast">
      <div v-if="error"
        class="absolute top-4 left-4 right-4 bg-red-500 opacity-90 text-white rounded-2xl z-[90] p-3 text-xs font-medium shadow-xl border border-white/10 flex items-center justify-between gap-2">
        <span class="flex items-center gap-1.5"> ⚠️ {{ error }} </span>
        <!-- Botón para cerrar manualmente -->
        <button @click="clearError" class="hover:bg-white/20 p-1 rounded-lg transition-colors text-sm leading-none">
          ✕
        </button>
      </div>
    </Transition>

    <!-- Caracteristica en desarrollo-->
    <Transition name="toast">
      <div v-if="folderFeatureInDevelopment" class="absolute top-4 left-4 right-4 bg-amber-500 text-white rounded-2xl z-[90] p-3 text-xs font-medium shadow-xl border border-white/10 flex items-center justify-between">
        <span>🚀 La carga de carpetas está en desarrollo. ¡Disponible pronto!</span>
        <button @click="folderFeatureInDevelopment = false" class="hover:bg-white/20 p-1 rounded-lg transition-colors text-sm leading-none">
          ✕
        </button>
      </div>
    </Transition>

    <Transition name="fade">
      <TransferProgress v-if="isUploading" :progress="uploadProgress" class="absolute inset-0 z-[60]" />
    </Transition>

    <Transition name="fade">
      <TransferSuccess v-if="isSuccess" :mode="transferMode" :link="generatedLink" @reset="resetAll"
        class="absolute inset-0 z-[70]" />
    </Transition>

    <div :class="{
      'opacity-20 pointer-events-none grayscale-[50%]': isUploading || isSuccess,
    }" class="transition-all duration-500">
      <div class="mb-5">
        <div v-if="items.length === 0" class="grid grid-cols-2 gap-3">
          <div @click="selectFiles"
            class="bg-[#e8edf7] hover:bg-[#dde4f3] dark:bg-[#1a1a1d] dark:hover:bg-[#202024] transition cursor-pointer rounded-2xl p-4 flex flex-col items-center justify-center gap-2">
            <div
              class="w-10 h-10 rounded-full bg-blue-600 text-white flex items-center justify-center text-xl font-bold">
              +
            </div>
            <span class="text-sm text-gray-700 dark:text-gray-200">Añadir archivos</span>
          </div>
          <div @click="selectFolder"
            class="bg-[#e8edf7] hover:bg-[#dde4f3] dark:bg-[#1a1a1d] dark:hover:bg-[#202024] transition cursor-pointer rounded-2xl p-4 flex flex-col items-center justify-center gap-2">
            <div class="w-10 h-10 rounded-full bg-blue-600 text-white flex items-center justify-center text-xl">
              📁
            </div>
            <span class="text-sm text-gray-700 dark:text-gray-200">Añadir carpetas</span>
          </div>
        </div>

        <div v-else class="space-y-3">
          <div
            class="bg-gray-50 dark:bg-[#1a1a1d] rounded-2xl p-3 space-y-2 max-h-[180px] overflow-y-auto custom-scrollbar">
            <FileItem v-for="(item, index) in items" :key="index" :item="item" @remove="removeItem(index)" />
          </div>

          <div class="flex justify-between items-center text-[11px] font-bold uppercase text-gray-400">
            <span>{{ items.length }} elementos</span>
            <div class="relative add-menu-container">
              <button @click.stop="showAddMenu = !showAddMenu"
                class="text-blue-500 flex items-center gap-1 font-bold tracking-tight">
                <span class="text-lg">+</span> Añadir más
              </button>
              <div v-if="showAddMenu"
                class="absolute right-0 mt-2 w-40 bg-white dark:bg-[#1a1a1d] border border-gray-100 dark:border-white/5 rounded-xl shadow-xl z-30">
                <div @click="
                  selectFiles();
                showAddMenu = false;
                "
                  class="px-4 py-3 text-sm hover:bg-gray-50 dark:hover:bg-white/5 cursor-pointer border-b dark:border-white/5">
                  📄 Archivo
                </div>
                <div @click="
                  selectFolder();
                showAddMenu = false;
                " class="px-4 py-3 text-sm hover:bg-gray-50 dark:hover:bg-white/5 cursor-pointer">
                  📁 Carpeta
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="space-y-0 border-t border-gray-100 dark:border-white/5 mt-4">
        <EmailSelector v-if="transferMode === 'email'" v-model="recipientEmails" />

        <input v-model="senderEmail" type="email" placeholder="Tu email"
          class="w-full bg-transparent border-b border-gray-100 dark:border-white/5 py-4 outline-none text-sm focus:border-blue-500 transition-colors" />

        <input v-model="subject" type="text" placeholder="Asunto"
          class="w-full bg-transparent border-b border-gray-100 dark:border-white/5 py-4 outline-none text-sm focus:border-blue-500 transition-colors" />

        <textarea v-model="message" placeholder="Mensaje"
          class="w-full bg-transparent border-b border-gray-100 dark:border-white/5 py-4 outline-none text-sm resize-none h-20 pt-4" />
      </div>

      <div class="mt-6 flex gap-3">
        <div class="relative flex-1 dropdown-expiry-container">
          <button @click.stop="showExpiry = !showExpiry"
            class="w-full border border-gray-200 dark:border-[#2a2a2e] rounded-2xl py-3 px-4 text-sm text-left flex justify-between bg-white dark:bg-transparent">
            📅 {{ selectedExpiry.label }}
          </button>
          <div v-if="showExpiry"
            class="absolute bottom-full mb-2 left-0 w-full bg-white dark:bg-[#1a1a1d] border border-gray-200 dark:border-[#2a2a2e] rounded-2xl shadow-xl z-50">
            <div v-for="opt in reversedExpiryOptions" :key="opt.value" @click="selectExpiry(opt)"
              class="px-4 py-3 text-sm hover:bg-gray-50 dark:hover:bg-[#202024] cursor-pointer">
              {{ opt.label }}
            </div>
          </div>
        </div>

        <button @click.stop="showOptions = !showOptions"
          class="w-14 border border-gray-200 dark:border-[#2a2a2e] rounded-2xl flex items-center justify-center hover:bg-gray-50 dark:hover:bg-white/5 transition-colors">
          ...
        </button>
      </div>

      <button @click="handleStartUpload"
        class="w-full bg-blue-600 hover:bg-blue-700 text-white py-4 rounded-2xl font-bold mt-4 shadow-lg shadow-blue-500/20 transition-all">
        Transferir
      </button>
    </div>

    <TransferOptions v-model:show="showOptions" v-model:mode="transferMode" />

    <input ref="fileInput" type="file" multiple class="hidden" @change="onFilesSelected" />
    <input ref="folderInput" type="file" webkitdirectory directory multiple class="hidden" @change="onFilesSelected" />
  </div>
</template>

<script setup>
import { ref, watch, computed, onMounted, onUnmounted } from "vue";
import { useFileProcessor } from "@/composables/useFileProcessor";
import { useTransferApi } from "@/composables/useTransferApi";

import FileItem from "./FileItem.vue";
import EmailSelector from "./EmailSelector.vue";
import TransferOptions from "./TransferOptions.vue";
import TransferProgress from "./TransferProgress.vue";
import TransferSuccess from "./TransferSuccess.vue";

const { items, files, addValidFiles, removeItem } = useFileProcessor();
const { startTransfer, isUploading, uploadProgress, error } = useTransferApi();

// Estados UI
const fileInput = ref(null);
const folderInput = ref(null);
const isSuccess = ref(false);
const showAddMenu = ref(false);
const showExpiry = ref(false);
const showOptions = ref(false);
const transferMode = ref("email"); // Controlado por TransferOptions
const generatedLink = ref("");

// Form
const recipientEmails = ref([]);
const senderEmail = ref("");
const subject = ref("");
const message = ref("");

const expiryOptions = [
  { value: 1, label: "1 día (Gratis)" },
  { value: 7, label: "7 días" },
  { value: 30, label: "30 días" },
  { value: 60, label: "60 días" },
  { value: 365, label: "1 año" },
];

// Crear la versión invertida para el dropdown, mejor experiencia para usuario.
const reversedExpiryOptions = computed(() => [...expiryOptions].reverse());

const selectedExpiry = ref(expiryOptions[0]);

// Handlers
const selectFiles = () => fileInput.value.click();

//Se deactiva temporalmente en lo que se desarrolla el manejo de estructura por carpetas
const folderFeatureInDevelopment = ref(false)
const selectFolder = () => {
    folderFeatureInDevelopment.value = true
    setTimeout(() => {
      folderFeatureInDevelopment.value = false
    }, 3000)
}
const onFilesSelected = (e) => {
  addValidFiles(Array.from(e.target.files));
  e.target.value = null;
};
const selectExpiry = (opt) => {
  selectedExpiry.value = opt;
  showExpiry.value = false;
};

//Testing
/*
const handleStartUpload = () => {
  if (items.value.length === 0) return alert('Añade archivos primero')
  isUploading.value = true
  uploadProgress.value = 0
  
  const interval = setInterval(() => {
    uploadProgress.value += 10
    if (uploadProgress.value >= 100) {
      clearInterval(interval)
      setTimeout(() => {
        isUploading.value = false
        isSuccess.value = true
        generatedLink.value = "https://tuapp.com/share/" + Math.random().toString(36).substring(7)
      }, 500)
    }
  }, 200)
}
*/

const handleStartUpload = async () => {
  if (files.value.length === 0) return alert("Añade archivos primero");
  if (!senderEmail.value) return alert("El email del remitente es obligatorio");
  if (transferMode.value === "email" && recipientEmails.value.length === 0) {
    return alert("Añade al menos un destinatario");
  }

  try {
    const metadata = {
      senderEmail: senderEmail.value,
      subject: subject.value,
      message: message.value,
      recipients: recipientEmails.value,
      totalFiles: files.value.length,
      expiresAt: selectedExpiry.value.value,
    };

    console.log(items)
    // startTransfer ahora es la función de Axios
    // Pasamos items.value que contiene los archivos reales del useFileProcessor
    const downloadToken = await startTransfer(metadata, files.value);

    // Si llegamos aquí, la subida fue exitosa
    generatedLink.value = `${window.location.origin}/download/${downloadToken}`;
    isSuccess.value = true;
  } catch (err) {
    // El error ya se guarda en la variable 'error' del composable
    alert("Error en la transferencia: " + err);
  }
};

const resetAll = () => {
  isSuccess.value = false;
  items.value = [];
  recipientEmails.value = [];
  senderEmail.value = "";
  subject.value = "";
  message.value = "";
};

// Función para limpiar el error manualmente
const clearError = () => {
  if (error) error.value = null
}


// Click Outside global
const handleClickOutside = (e) => {
  if (!e.target.closest(".add-menu-container")) showAddMenu.value = false;
  if (!e.target.closest(".dropdown-expiry-container")) showExpiry.value = false;
};

watch(error, (newError) => {
  if (newError) {
    setTimeout(() => {
      clearError()
    }, 4000)
  }
})
onMounted(() => window.addEventListener("click", handleClickOutside));
onUnmounted(() => window.removeEventListener("click", handleClickOutside));
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #2a2a2e;
  border-radius: 10px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.4s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}

.toast-enter-from {
  opacity: 0;
  transform: translateY(-20px);
}

.toast-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
