<template>
    <div class="w-full max-w-[420px] mx-auto bg-[#ffffff] dark:bg-[#111113] rounded-3xl shadow-[0_20px_60px_rgba(0,0,0,.25)] p-5 border border-[rgba(0,0,0,0.05)] dark:border-[rgba(255,255,255,.05)]"
        @dragover.prevent="isDragging = true" @dragleave="isDragging = false" @drop.prevent="handleDrop"
        :class="isDragging ? 'ring-2 ring-blue-500' : ''">

        <!-- Upload Section -->
        <div class="mb-5">

            <!-- SIN archivos -->
            <div v-if="items.length === 0" class="grid grid-cols-2 gap-3">
                <div class="bg-[#e8edf7] hover:bg-[#dde4f3] dark:bg-[#1a1a1d] dark:hover:bg-[#202024] transition cursor-pointer rounded-2xl p-4 flex flex-col items-center justify-center gap-2"
                    @click="selectFiles">
                    <div
                        class="w-10 h-10 rounded-full bg-blue-600 text-white flex items-center justify-center text-lg font-bold">
                        +
                    </div>
                    <span class="text-sm text-gray-700 dark:text-gray-200">
                        Añadir archivos
                    </span>
                </div>

                <div class="bg-[#e8edf7] hover:bg-[#dde4f3] dark:bg-[#1a1a1d] dark:hover:bg-[#202024] transition cursor-pointer rounded-2xl p-4 flex flex-col items-center justify-center gap-2"
                    @click="selectFolder">
                    <div class="w-10 h-10 rounded-full bg-blue-600 text-white flex items-center justify-center text-lg">
                        📁
                    </div>
                    <span class="text-sm text-gray-700 dark:text-gray-200">
                        Añadir carpetas
                    </span>
                </div>
            </div>

            <!-- CON archivos -->
            <div v-else class="space-y-3">

                <!-- Lista -->
                <div class="bg-[#1a1a1d] rounded-2xl p-3 space-y-2 max-h-[193px] overflow-y-auto custom-scrollbar">

                    <div v-for="(item, index) in items" :key="index"
                        class="flex justify-between items-center text-sm text-gray-300 bg-[#202024] px-3 py-2 rounded-xl">

                        <!-- 📁 Carpeta -->
                        <div v-if="item.type === 'folder'" class="flex flex-col w-[70%] truncate">
                            <span>📁 {{ item.name }}</span>
                            <span class="text-xs text-gray-400">
                                {{ item.files.length }} elementos
                            </span>
                        </div>

                        <!-- 📄 Archivo -->
                        <div v-else class="flex flex-col w-[70%] truncate">
                            <span>📄 {{ item.file.name }}</span>
                            <span class="text-xs text-gray-400">
                                {{ formatSize(item.file.size) }}
                            </span>
                        </div>

                        <button @click="removeItem(index)" class="text-xs hover:text-red-400">
                            ✕
                        </button>
                    </div>

                </div>

                <!-- Footer -->
                <div class="flex justify-between items-center text-sm text-gray-400">
                    <span>{{ items.length }} elementos</span>

                    <!-- 🔥 Dropdown Añadir más -->
                    <div class="relative">
                        <button @click.stop="showAddMenu = !showAddMenu"
                            class="flex items-center gap-2 text-blue-500 hover:opacity-80">
                            ➕ Añadir más
                        </button>

                        <div v-if="showAddMenu"
                            class="absolute right-0 mt-2 w-40 bg-white dark:bg-[#111113] border border-gray-200 dark:border-[#2a2a2e] rounded-xl shadow-lg z-30">

                            <div @click="selectFiles(); showAddMenu = false"
                                class="px-4 py-2 text-sm hover:bg-gray-100 dark:hover:bg-[#1a1a1d] cursor-pointer">
                                📄 Archivo
                            </div>

                            <div @click="selectFolder(); showAddMenu = false"
                                class="px-4 py-2 text-sm hover:bg-gray-100 dark:hover:bg-[#1a1a1d] cursor-pointer">
                                📁 Carpeta
                            </div>

                        </div>
                    </div>

                </div>

            </div>
        </div>

        <!-- Inputs -->
        <input ref="fileInput" type="file" multiple class="hidden" @change="handleFiles" />
        <input ref="folderInput" type="file" webkitdirectory directory multiple class="hidden" @change="handleFiles" />
        <!-- Email To -->
        <div class="border-b border-gray-200 dark:border-[#2a2a2e]">
            <div class="flex justify-end text-sm text-gray-600 dark:text-gray-400">
                <span>{{ recipientEmails.length }} de 3</span>
            </div>

            <div class="flex flex-wrap gap-2 mb-2">
                <div v-for="(email, index) in recipientEmails" :key="index"
                    class="bg-blue-100 dark:bg-blue-900/40 text-blue-700 dark:text-blue-300 text-xs px-3 py-1 rounded-full flex items-center gap-2">
                    {{ email }}
                    <button @click="removeEmail(index)" class="text-xs hover:opacity-70">
                        ✕
                    </button>
                </div>
            </div>
            <div class="border-b transition-colors duration-300 py-3"
                :class="recipientEmailError ? 'border-red-500' : 'border-gray-200 dark:border-[#2a2a2e]'">
                <input v-model="recipientEmailInput"
                    @keydown.enter.prevent="addRecipientEmail"
                    @keydown="handleKeydown"
                    type="text" :disabled="recipientEmails.length >= 3"
                    :placeholder="recipientEmails.length >= 3 ? 'Límite de correos alcanzado' : 'Enviar email a'"
                    class="w-full bg-transparent outline-none text-sm text-gray-700 dark:text-gray-200 placeholder:text-gray-400" />
            </div>
            <Transition name="fade">
                <p v-if="recipientEmailError" class="text-[11px] text-red-500 mt-1 font-medium">
                    {{ recipientEmailError }}
                </p>
            </Transition>
        </div>

        <!-- Inputs extra -->
        <div class="border-b border-gray-200 dark:border-[#2a2a2e] py-3">
            <input v-model="senderEmail" type="email" placeholder="Tu email"
                class="w-full bg-transparent outline-none text-sm text-gray-700 dark:text-gray-200 placeholder:text-gray-400" />
        </div>

        <div class="border-b border-gray-200 dark:border-[#2a2a2e] py-3">
            <input type="text" placeholder="Asunto"
                class="w-full bg-transparent outline-none text-sm text-gray-700 dark:text-gray-200 placeholder:text-gray-400" />
        </div>

        <div class="border-b border-gray-200 dark:border-[#2a2a2e] pt-5 mb-5">
            <textarea placeholder="Mensaje"
                class="w-full bg-transparent outline-none text-sm text-gray-700 dark:text-gray-200 placeholder:text-gray-400 resize-none overflow-y-auto custom-scrollbar" />
        </div>

        <!-- Opciones -->
        <div class="flex gap-3 mb-5">

            <div class="relative flex-1">
                <button @click.stop="toggleExpiry"
                    class="w-full border border-gray-200 dark:border-[#2a2a2e] rounded-2xl py-3 px-4 text-sm text-gray-700 dark:text-gray-200 bg-white dark:bg-[#111113] flex justify-between">
                    📅 {{ selectedExpiry.label }}
                </button>

                <!-- Dropdown -->
                <div v-if="showExpiry"
                    class="absolute bottom-14 left-0 w-full bg-white dark:bg-[#111113] border border-gray-200 dark:border-[#2a2a2e] rounded-2xl shadow-xl overflow-hidden z-20">
                    <div v-for="option in reversedExpiryOptions" :key="option.value" @click="selectExpiry(option)"
                        class="px-4 py-3 text-sm cursor-pointer hover:bg-gray-50 dark:hover:bg-[#1a1a1d] flex justify-between items-center">
                        <span>{{ option.label }}</span>
                        <span v-if="selectedExpiry.value === option.value" class="text-blue-500">✓</span>
                    </div>
                </div>
            </div>

            <button @click="toggleMenu" class="w-14 border rounded-2xl py-3">
                ...
            </button>

        </div>

        <!-- Botón -->
        <button @click="startUpload" class="w-full bg-blue-600 hover:bg-blue-700 text-white py-4 rounded-2xl">
            <span v-if="!isUploading">Transferir</span>
            <span v-else>Subiendo... {{ uploadProgress }}%</span>
        </button>

        <!-- Barra -->
        <div v-if="isUploading" class="mt-3 w-full bg-gray-700 rounded-full h-2">
            <div class="bg-blue-500 h-2 rounded-full" :style="{ width: uploadProgress + '%' }"></div>
        </div>

    </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'

const fileInput = ref(null)
const folderInput = ref(null)

const files = ref([])
const items = ref([])

const recipientEmails = ref([])
const recipientEmailInput = ref('')
const senderEmail = ref('')

const isDragging = ref(false)
const showAddMenu = ref(false)

// ================= FILES =================
const selectFiles = () => fileInput.value.click()
const selectFolder = () => folderInput.value.click()

const handleFiles = (event) => {
    const selected = Array.from(event.target.files)
    addValidFiles(selected)
    event.target.value = null
}

const handleDrop = (event) => {
    isDragging.value = false
    const dropped = Array.from(event.dataTransfer.files)
    addValidFiles(dropped)
}

const addValidFiles = (newFiles) => {
    const filtered = newFiles.filter(file => {
        const rootName = file.webkitRelativePath
            ? file.webkitRelativePath.split('/')[0]
            : file.name

        const type = file.webkitRelativePath ? 'folder' : 'file'

        const exists = items.value.some(item => {
            return item.name === rootName && item.type === type
        })

        if (exists) {
            console.warn(`${type === 'folder' ? 'Carpeta' : 'Archivo'} duplicado: ${rootName}`)
            return false
        }
        return true
    })

    if (filtered.length > 0) {
        files.value.push(...filtered)
        processFiles(files.value)
    }
}
// ================= GROUP =================
const processFiles = (fileList) => {
    const map = {}

    fileList.forEach(file => {
        const path = file.webkitRelativePath || file.name
        const parts = path.split('/')

        if (parts.length > 1) {
            const folderName = parts[0]

            if (!map[folderName]) {
                map[folderName] = {
                    type: 'folder',
                    name: folderName,
                    files: []
                }
            }

            map[folderName].files.push(file)
        } else {
            map[file.name] = {
                type: 'file',
                name: file.name,
                file
            }
        }
    })

    items.value = Object.values(map)
}

// ================= REMOVE =================
const removeItem = (index) => {
    const item = items.value[index]

    if (item.type === 'folder') {
        files.value = files.value.filter(f =>
            !f.webkitRelativePath.startsWith(item.name + '/')
        )
    } else {
        files.value = files.value.filter(f => f !== item.file)
    }

    processFiles(files.value)
}

// ================= SIZE =================
const formatSize = (bytes) => {
    if (bytes < 1024) return bytes + ' B'
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
    if (bytes < 1024 * 1024 * 1024) return (bytes / 1024 / 1024).toFixed(1) + ' MB'
    return (bytes / 1024 / 1024 / 1024).toFixed(1) + ' GB'
}

// ================= EMAIL =================
const recipientEmailError = ref('')
const addRecipientEmail = () => {
    const email = recipientEmailInput.value.trim()
    recipientEmailError.value = ''
    if (!email) return
    if (recipientEmails.value.length >= 3) {
        recipientEmailError.value = 'Máximo 3 destinatarios'
        return
    }
    if (recipientEmails.value.includes(email)) {
        recipientEmailError.value = 'Este correo ya fue añadido'
        return
    }
    if (!isValidEmail(email)) {
        recipientEmailError.value = 'Formato de email inválido'
        console.warn("Formato de email inválido")
        return
    }
    recipientEmails.value.push(email)
    recipientEmailInput.value = ''
}

watch(recipientEmailInput, () => {
    if (recipientEmailError.value) recipientEmailError.value = ''
})

const removeEmail = (index) => recipientEmails.value.splice(index, 1)

const handleKeydown = (event) => {
    if (event.key === ',' || event.key === 'Tab') {
        event.preventDefault()
        addRecipientEmail()
    }
}

// ================= DROPDOWN =================
const showExpiry = ref(false)
const expiryOptions = [
    { value: 1, label: '1 día (Gratis)' },
    { value: 3, label: '3 días (Gratis)' },
    { value: 7, label: '7 días' },
    { value: 30, label: '30 días' },
    { value: 60, label: '60 días' },
    { value: 365, label: '1 año' }
]
// Crear la versión invertida para el dropdown
const reversedExpiryOptions = computed(() => [...expiryOptions].reverse())

const selectedExpiry = ref(expiryOptions[0])

const toggleExpiry = () => showExpiry.value = !showExpiry.value
const selectExpiry = (option) => {
    selectedExpiry.value = option
    showExpiry.value = false
}

const handleClickOutside = (e) => {
    if (!e.target.closest('.dropdown-expiry-container')) showExpiry.value = false
    if (!e.target.closest('.add-menu-container')) showAddMenu.value = false
    if (!e.target.closest('.main-menu-container')) showMenu.value = false
}

const handleEscKey = (e) => {
    if (e.key === 'Escape') {
        showExpiry.value = false
        showAddMenu.value = false
        showMenu.value = false
    }
}

onMounted(() => {
    window.addEventListener('click', handleClickOutside)
    window.addEventListener('keydown', handleEscKey)
})

onUnmounted(() => {
    window.removeEventListener('click', handleClickOutside)
    window.removeEventListener('keydown', handleEscKey)
})

// ================= MENU =================
const showMenu = ref(false)
const toggleMenu = () => showMenu.value = !showMenu.value

// ================= UPLOAD =================
const uploadProgress = ref(0)
const isUploading = ref(false)

const startUpload = () => {
    isUploading.value = true
    uploadProgress.value = 0

    const interval = setInterval(() => {
        uploadProgress.value += 10
        if (uploadProgress.value >= 100) {
            clearInterval(interval)
            isUploading.value = false
        }
    }, 300)
}
// ================= VALIDACION DE CAMPOS =================
const isValidEmail = (email) => {
    const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return regex.test(email)
}

</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
    width: 4px;
}

.custom-scrollbar::-webkit-scrollbar-track {
    background: transparent;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
    background: #2a2a2e;
    border-radius: 10px;
}

.custom-scrollbar::-webkit-scrollbar-thumb:hover {
    background: #3b82f6;
    /* Cambia a azul al pasar el mouse */
}

/* Transición suave para el mensaje de error */
.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
    transform: translateY(-5px);
}
</style>
