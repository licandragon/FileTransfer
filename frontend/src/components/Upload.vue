<template>
    <div
        class="w-[420px] bg-[#ffffff] dark:bg-[#111113] rounded-3xl shadow-[0_20px_60px_rgba(0,0,0,.25)] p-5 border border-[rgba(0,0,0,0.05)] dark:border-[rgba(255,255,255,.05)]">

        <!-- Upload Buttons -->
        <div class="grid grid-cols-2 gap-3 mb-5">

            <!-- Add Files -->
            <div class="bg-[#e8edf7] hover:bg-[#dde4f3] dark:bg-[#1a1a1d] dark:hover:bg-[#202024] transition cursor-pointer rounded-2xl p-4 flex flex-col items-center justify-center gap-2 border border-transparent hover:border-[rgba(0,0,0,.05)] dark:hover:border-[rgba(255,255,255,.05)]"
                @click="selectFiles">
                <div
                    class="w-10 h-10 rounded-full bg-blue-600 text-white flex items-center justify-center text-lg font-bold shadow">
                    +
                </div>
                <span class="text-sm font-medium text-gray-700 dark:text-gray-200">
                    Añadir archivos
                </span>
            </div>

            <!-- Add Folders -->
            <div class="bg-[#e8edf7] hover:bg-[#dde4f3] dark:bg-[#1a1a1d] dark:hover:bg-[#202024] transition cursor-pointer rounded-2xl p-4 flex flex-col items-center justify-center gap-2 border border-transparent hover:border-[rgba(0,0,0,.05)] dark:hover:border-[rgba(255,255,255,.05)]"
                @click="selectFolder">
                <div
                    class="w-10 h-10 rounded-full bg-blue-600 text-white flex items-center justify-center text-lg shadow">
                    📁
                </div>
                <span class="text-sm font-medium text-gray-700 dark:text-gray-200">
                    Añadir carpetas
                </span>
            </div>
        </div>

        <!-- Hidden Inputs -->
        <input ref="fileInput" type="file" multiple class="hidden" @change="handleFiles" />

        <input ref="folderInput" type="file" webkitdirectory directory multiple class="hidden" @change="handleFiles" />

        <!-- Info Row -->
        <div class="flex justify-between items-center text-sm mb-5">
            <span class="text-gray-500 dark:text-gray-400">
                Consigue transferencias ilimitadas
            </span>

            <button class="text-purple-600 font-medium hover:opacity-80 transition">
                ⚡ Aumentar límite
            </button>
        </div>

        <!-- Email To -->
        <div class="border-b border-gray-200 dark:border-[#2a2a2e] py-3">
            <div class="flex justify-between text-sm text-gray-600 dark:text-gray-400 mb-2">
                <span>Enviar email a</span>
                <span>{{ emails.length }} de 3</span>
            </div>

            <!-- Email Tags -->
            <div class="flex flex-wrap gap-2 mb-2">
                <div v-for="(email, index) in emails" :key="index"
                    class="bg-blue-100 dark:bg-blue-900/40 text-blue-700 dark:text-blue-300 text-xs px-3 py-1 rounded-full flex items-center gap-2">
                    {{ email }}
                    <button @click="removeEmail(index)" class="text-xs hover:opacity-70">
                        ✕
                    </button>
                </div>
            </div>

            <!-- Input -->
            <input v-model="emailInput" @keydown.enter.prevent="addEmail" @keydown="handleKeydown" type="text"
                placeholder="Agregar correo"
                class="w-full bg-transparent outline-none text-sm text-gray-700 dark:text-gray-200 placeholder:text-gray-400" />
        </div>

        <!-- Your Email -->
        <div class="border-b border-gray-200 dark:border-[#2a2a2e] py-3">
            <input type="email" placeholder="Tu email"
                class="w-full bg-transparent outline-none text-sm text-gray-700 dark:text-gray-200 placeholder:text-gray-400" />
        </div>

        <!-- Title -->
        <div class="border-b border-gray-200 dark:border-[#2a2a2e] py-3">
            <input type="text" placeholder="Título"
                class="w-full bg-transparent outline-none text-sm text-gray-700 dark:text-gray-200 placeholder:text-gray-400" />
        </div>

        <div class="border-b border-gray-200 dark:border-[#2a2a2e] pt-5 mb-5">
            <textarea type="text" name="message" style="resize: none;" placeholder="Mensaje"
                class="w-full bg-transparent outline-none text-sm text-gray-700 dark:text-gray-200 placeholder:text-gray-400" />
        </div>

        <!-- Bottom Options -->
        <div class="flex gap-3 mb-5">

            <div class="relative flex-1">
                <button @click="toggleExpiry"
                    class="w-full border border-gray-200 dark:border-[#2a2a2e] rounded-2xl py-3 px-4 text-sm text-gray-700 dark:text-gray-200 bg-white dark:bg-[#111113] hover:bg-gray-50 dark:hover:bg-[#1a1a1d] transition flex items-center justify-between">
                    <span>📅 {{ selectedExpiry.label }}</span>
                    <span class="text-xs">▾</span>
                </button>

                <!-- Dropdown -->
                <div v-if="showExpiry"
                    class="absolute bottom-14 left-0 w-full bg-white dark:bg-[#111113] border border-gray-200 dark:border-[#2a2a2e] rounded-2xl shadow-xl overflow-hidden z-20">
                    <div v-for="option in expiryOptions" :key="option.value" @click="selectExpiry(option)"
                        class="px-4 py-3 text-sm cursor-pointer hover:bg-gray-50 dark:hover:bg-[#1a1a1d] flex justify-between items-center">
                        <span>{{ option.label }}</span>
                        <span v-if="selectedExpiry.value === option.value">✓</span>
                    </div>
                </div>
            </div>

            <button
                class="w-14 border border-gray-200 dark:border-[#2a2a2e] rounded-2xl py-3 flex items-center justify-center text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-[#1a1a1d] transition">
                ...
            </button>

        </div>

        <!-- Transfer Button -->
        <button
            class="w-full bg-blue-600 hover:bg-blue-700 transition text-white py-4 rounded-2xl font-medium shadow-lg hover:shadow-xl">
            Transferir
        </button>

    </div>
</template>

<script setup>
import { ref } from 'vue'

const fileInput = ref(null)
const folderInput = ref(null)

const files = ref([])

const emails = ref([])
const emailInput = ref('')

const addEmail = () => {
    const email = emailInput.value.trim()

    if (!email) return
    if (emails.value.length >= 3) return
    if (emails.value.includes(email)) return

    emails.value.push(email)
    emailInput.value = ''
}

const removeEmail = (index) => {
    emails.value.splice(index, 1)
}

const handleKeydown = (event) => {
    if (event.key === ',' || event.key === 'Tab') {
        event.preventDefault()
        addEmail()
    }
}

const showExpiry = ref(false)

const expiryOptions = [
    { value: 1, label: '1 día (Gratis)' },
    { value: 3, label: '3 días (Gratis)' },
    { value: 7, label: '7 días' },
    { value: 30, label: '30 días' },
    { value: 60, label: '60 días' },
    { value: 365, label: '1 año' }
]

const selectedExpiry = ref(expiryOptions[1])

const toggleExpiry = () => {
    showExpiry.value = !showExpiry.value
}

const selectExpiry = (option) => {
    selectedExpiry.value = option
    showExpiry.value = false
}

const selectFiles = () => {
    fileInput.value.click()
}

const selectFolder = () => {
    folderInput.value.click()
}

const handleFiles = (event) => {
    const selected = Array.from(event.target.files)
    files.value.push(...selected)

    console.log('Files:', files.value)
}
</script>

<style scoped></style>
