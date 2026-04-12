<!-- src/views/RegisterView.vue -->
<template>
    <div class="flex items-center justify-center min-h-[calc(100vh-65px)] px-4 py-8">

        <div class="w-full max-w-md
      bg-[rgba(22,22,26,0.85)]
      border border-[rgba(255,255,255,0.06)]
      rounded-2xl
      p-7
      backdrop-blur
      animate-fade">

            <!-- HEADER -->
            <div class="text-center mb-6">
                <h1 class="text-[1.7rem] font-bold text-[#e8e4dc] tracking-[-0.02em]">
                    Crear cuenta
                </h1>

                <p class="text-[#6b6560] text-sm mt-1">
                    Regístrate para comenzar
                </p>
            </div>

            <div v-if="errorMsg" class="text-red-400 text-sm mb-3">
                {{ errorMsg }}
            </div>

            <div v-if="successMsg" class="text-green-400 text-sm mb-3">
                {{ successMsg }}
            </div>

            <form class="flex flex-col gap-4" @submit.prevent="handleRegister">

                <!-- NAME -->
                <div class="flex flex-col gap-1">
                    <label class="text-[0.8rem] text-[#9a938c]">
                        Nombre
                    </label>

                    <input v-model="name" type="text" placeholder="Tu nombre" required class="input w-full" />
                </div>

                <!-- EMAIL -->
                <div class="flex flex-col gap-1">
                    <label class="text-[0.8rem] text-[#9a938c]">
                        Email
                    </label>

                    <input v-model="email" type="email" placeholder="correo@email.com" required class="input w-full" />
                </div>

                <!-- PASSWORD -->
                <div class="flex flex-col gap-2">

                    <label class="text-[0.8rem] text-[#9a938c]">
                        Contraseña
                    </label>

                    <div class="relative">

                        <input v-model="password" :type="showPassword ? 'text' : 'password'" placeholder="••••••••"
                            required class="input w-full pr-10" />

                        <!-- toggle -->
                        <button type="button" @click="showPassword = !showPassword" class="absolute right-3 top-1/2 -translate-y-1/2
              text-[#6b6560] hover:text-[#c8a96e] transition">
                            {{ showPassword ? "🙈" : "👁️" }}
                        </button>

                    </div>

                    <!-- strength bar -->
                    <div class="h-1 w-full bg-[#1b1b1e] rounded overflow-hidden">
                        <div class="h-full transition-all duration-300" :class="strengthColor"
                            :style="{ width: strengthWidth }"></div>
                    </div>

                    <!-- checklist -->
                    <ul class="text-[0.75rem] flex flex-col gap-1 mt-1">

                        <li :class="ruleLength ? 'text-green-400' : 'text-[#6b6560]'">
                            ✔ mínimo 8 caracteres
                        </li>

                        <li :class="ruleUpper ? 'text-green-400' : 'text-[#6b6560]'">
                            ✔ una letra mayúscula
                        </li>

                        <li :class="ruleNumber ? 'text-green-400' : 'text-[#6b6560]'">
                            ✔ al menos un número
                        </li>

                        <li :class="ruleSymbol ? 'text-green-400' : 'text-[#6b6560]'">
                            ✔ un símbolo (!@#$)
                        </li>

                    </ul>

                </div>

                <!-- CONFIRM PASSWORD -->
                <div class="flex flex-col gap-1">

                    <label class="text-[0.8rem] text-[#9a938c]">
                        Confirmar contraseña
                    </label>

                    <input v-model="confirmPassword" :type="showPassword ? 'text' : 'password'" placeholder="••••••••"
                        required class="input w-full" />

                    <span v-if="confirmPassword && password !== confirmPassword" class="text-[0.75rem] text-red-400">
                        Las contraseñas no coinciden
                    </span>

                </div>

                <!-- SUBMIT -->
                <button :disabled="!formValid || isLoading" type="submit" class="mt-2
          py-[0.65rem]
          rounded-lg
          font-semibold
          text-[0.9rem]
          transition
          bg-linear-to-br from-[#c8a96e] to-[#a0834e]
          text-[#0d0d0f]
          hover:opacity-90
          disabled:opacity-40">
                    Crear cuenta
                </button>

            </form>

            <!-- LOGIN LINK -->
            <p class="text-center text-sm text-[#6b6560] mt-6">
                ¿Ya tienes cuenta?

                <RouterLink to="/login" class="text-[#c8a96e] hover:underline">
                    Inicia sesión
                </RouterLink>
            </p>

        </div>

    </div>
</template>

<script setup>
import { ref, computed } from "vue"
import { RouterLink, useRouter } from "vue-router"
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const { register } = useAuth()

const name = ref("")
const email = ref("")
const password = ref("")
const confirmPassword = ref("")
const showPassword = ref(false)

const isLoading = ref(false)
const errorMsg = ref('')
const successMsg = ref('')

/* reglas de contraseña */

const ruleLength = computed(() => password.value.length >= 8)
const ruleUpper = computed(() => /[A-Z]/.test(password.value))
const ruleNumber = computed(() => /[0-9]/.test(password.value))
const ruleSymbol = computed(() => /[^A-Za-z0-9]/.test(password.value))

/* fortaleza */

const strength = computed(() => {
    let score = 0
    if (ruleLength.value) score++
    if (ruleUpper.value) score++
    if (ruleNumber.value) score++
    if (ruleSymbol.value) score++
    return score
})

const strengthWidth = computed(() => `${strength.value * 25}%`)

const strengthColor = computed(() => {
    const colors = [
        "bg-red-500",
        "bg-orange-500",
        "bg-yellow-500",
        "bg-green-500",
        "bg-emerald-400"
    ]
    return colors[strength.value]
})

/* validación del formulario */

const formValid = computed(() => {

    return (
        name.value &&
        email.value &&
        ruleLength.value &&
        ruleUpper.value &&
        ruleNumber.value &&
        ruleSymbol.value &&
        password.value === confirmPassword.value
    )

})

async function handleRegister() {

    errorMsg.value = ''
    successMsg.value = ''

    if (!formValid.value) return

    isLoading.value = true

    try {

        await register(
            name.value,
            email.value,
            password.value
        )

        successMsg.value = 'Cuenta creada. Revisa tu correo para confirmar.'

        // opcional: redirigir después
        setTimeout(() => {
            router.push('/login')
        }, 1500)

    } catch (error) {

        console.error('Error:', error)

        if (error.message?.includes('User already registered')) {
            errorMsg.value = 'Este correo ya está registrado'
        } else if (error.message?.includes('Password')) {
            errorMsg.value = 'La contraseña no es válida'
        } else {
            errorMsg.value = 'Error al crear la cuenta'
        }

    } finally {
        isLoading.value = false
    }
}
</script>

<style scoped>
.input {
    padding: 0.55rem 0.75rem;
    border-radius: 0.6rem;
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid rgba(255, 255, 255, 0.08);
    color: #e8e4dc;
    font-size: 0.9rem;
    outline: none;
    transition: all .2s;
}

.input:focus {
    border-color: #c8a96e;
    box-shadow: 0 0 0 1px rgba(200, 169, 110, 0.2);
}

.animate-fade {
    animation: fadeIn .4s ease both;
}

@keyframes fadeIn {
    from {
        opacity: 0;
        transform: translateY(10px);
    }

    to {
        opacity: 1;
        transform: translateY(0);
    }
}
</style>