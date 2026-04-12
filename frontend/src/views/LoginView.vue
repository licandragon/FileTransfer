<!-- src/views/LoginView.vue -->
<template>
  <div class="login-page">

    <!-- Fondo decorativo -->
    <div class="bg-grid" aria-hidden="true"></div>
    <div class="bg-glow" aria-hidden="true"></div>

    <!-- CARD -->
    <div class="login-card">

      <!-- HEADER -->
      <div class="text-center mb-7">

        <span class="card-icon">⬡</span>

        <h1
          class="text-[1.6rem] font-bold tracking-[-0.03em]
          bg-linear-to-br from-[#e8e4dc] to-[#c8a96e]
          bg-clip-text text-transparent mb-1"
        >
          Bienvenido
        </h1>

        <p class="text-[#6b6560] text-[0.88rem]">
          Ingresa tus credenciales para continuar
        </p>

      </div>

      <!-- ERROR -->
      <Transition name="alert">
        <div v-if="errorMsg" class="alert-error" role="alert">
          <span>⚠</span>
          {{ errorMsg }}
        </div>
      </Transition>

      <!-- FORM -->
      <form
        class="flex flex-col gap-4"
        @submit.prevent="handleLogin"
        novalidate
      >

        <!-- EMAIL -->
        <div class="flex flex-col gap-1.5">

          <label
            for="email"
            class="text-[0.82rem] font-semibold text-[#9a938c]
            uppercase tracking-[0.04em]"
          >
            Correo electrónico
          </label>

          <div class="field-input-wrap">

            <span class="field-icon">✉</span>

            <input
              id="email"
              v-model.trim="form.email"
              type="email"
              class="field-input"
              placeholder="usuario@ejemplo.com"
              autocomplete="email"
              @blur="validateEmail"
            />

          </div>

          <span v-if="errors.email" class="field-error">
            {{ errors.email }}
          </span>

        </div>

        <!-- PASSWORD -->
        <div class="flex flex-col gap-1.5">

          <label
            for="password"
            class="text-[0.82rem] font-semibold text-[#9a938c]
            uppercase tracking-[0.04em]"
          >
            Contraseña
          </label>

          <div class="field-input-wrap">

            <span class="field-icon">🔑</span>

            <input
              id="password"
              v-model="form.password"
              :type="showPassword ? 'text' : 'password'"
              class="field-input"
              placeholder="••••••••"
              autocomplete="current-password"
              @blur="validatePassword"
            />

            <button
              type="button"
              class="field-toggle"
              @click="showPassword = !showPassword"
              :aria-label="showPassword ? 'Ocultar contraseña' : 'Ver contraseña'"
            >
              {{ showPassword ? '🙈' : '👁' }}
            </button>

          </div>

          <span v-if="errors.password" class="field-error">
            {{ errors.password }}
          </span>

        </div>

        <!-- HINT -->
        <div class="hint-box">

          <span class="hint-label">
            Credenciales de prueba:
          </span>

          <code>admin@vue.com</code>
          /
          <code>vue1234</code>

        </div>

        <!-- SUBMIT -->
        <button
          type="submit"
          class="btn-submit"
          :disabled="isLoading"
        >

          <span v-if="!isLoading">
            Iniciar sesión →
          </span>

          <span v-else class="spinner"></span>

        </button>

      </form>
      <DiscordLoginButton
      class="btn-submit"
      />

    </div>

  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import DiscordLoginButton from '@/components/DiscordLoginButton.vue'

const router = useRouter()
const { login } = useAuth()

const form = reactive({
  email: '',
  password: '',
})

const errors = reactive({
  email: '',
  password: '',
})

const showPassword = ref(false)
const isLoading = ref(false)
const errorMsg = ref('')

function validateEmail() {
  if (!form.email) {
    errors.email = 'El correo es obligatorio.'
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
    errors.email = 'Ingresa un correo válido.'
  } else {
    errors.email = ''
  }
}

function validatePassword() {
  if (!form.password) {
    errors.password = 'La contraseña es obligatoria.'
  } else if (form.password.length < 6) {
    errors.password = 'Mínimo 6 caracteres.'
  } else {
    errors.password = ''
  }
}

function isFormValid() {
  validateEmail()
  validatePassword()
  return !errors.email && !errors.password
}

async function handleLogin() {
  errorMsg.value = ''

  if (!isFormValid()) return

  isLoading.value = true

  try {

    await login(form.email, form.password)

    await router.push('/dashboard')

  } catch (error) {

    console.error('Error Completo:', error)

    if (error.message?.includes('Invalid login')) {
      errorMsg.value =
        'Credenciales incorrectas. Intenta de nuevo.'
    } else if (error.message?.includes('Email not confirmed')) {
      errorMsg.value = 'Debes confirmar tu correo'
    } else {
      errorMsg.value = 'Error inesperado'
    }

  } finally {
    isLoading.value = false
  }
}

</script>

<style scoped>

/* layout */
.login-page{
  min-height: calc(100vh - 65px);
  display:flex;
  align-items:center;
  justify-content:center;
  position:relative;
  overflow:hidden;
  padding:2rem 1rem;
}

/* grid background */
.bg-grid{
  position:absolute;
  inset:0;
  background-image:
  linear-gradient(rgba(200,169,110,.04) 1px,transparent 1px),
  linear-gradient(90deg,rgba(200,169,110,.04) 1px,transparent 1px);
  background-size:48px 48px;
  pointer-events:none;
}

/* glow */
.bg-glow{
  position:absolute;
  width:500px;
  height:500px;
  background:radial-gradient(circle,rgba(200,169,110,.08)0%,transparent 70%);
  top:50%;
  left:50%;
  transform:translate(-50%,-50%);
  filter:blur(20px);
}

/* card */
.login-card{
  width:100%;
  max-width:420px;
  background:rgba(22,22,26,.95);
  border:1px solid rgba(200,169,110,.15);
  border-radius:16px;
  padding:2.5rem 2rem;
  box-shadow:
  0 0 0 1px rgba(255,255,255,.03),
  0 24px 48px rgba(0,0,0,.5);
  animation:cardIn .4s ease both;
}

/* icon */
.card-icon{
  display:block;
  font-size:2rem;
  color:#c8a96e;
  margin-bottom:.75rem;
  animation:pulse 3s ease-in-out infinite;
}

/* alert */
.alert-error{
  background:rgba(224,112,112,.1);
  border:1px solid rgba(224,112,112,.3);
  color:#e07070;
  border-radius:8px;
  padding:.65rem 1rem;
  font-size:.88rem;
  margin-bottom:1.25rem;
  display:flex;
  align-items:center;
  gap:.5rem;
}

/* input wrapper */
.field-input-wrap{
  position:relative;
  display:flex;
  align-items:center;
}

/* icon */
.field-icon{
  position:absolute;
  left:.85rem;
  opacity:.5;
}

/* input */
.field-input{
  width:100%;
  background:rgba(255,255,255,.04);
  border:1px solid rgba(255,255,255,.08);
  border-radius:8px;
  padding:.65rem 2.75rem .65rem 2.5rem;
  color:#e8e4dc;
  outline:none;
  transition:border-color .2s, background .2s;
}

.field-input::placeholder{
  color:#3e3a36;
}

.field-input:focus{
  border-color:rgba(200,169,110,.5);
  background:rgba(200,169,110,.04);
}

.field-input:focus-visible{
  box-shadow:0 0 0 2px rgba(200,169,110,.35);
}

/* toggle */
.field-toggle{
  position:absolute;
  right:.75rem;
  background:none;
  border:none;
  cursor:pointer;
  opacity:.6;
}

.field-toggle:hover{
  opacity:1;
}

/* error */
.field-error{
  font-size:.8rem;
  color:#e07070;
}

/* hint */
.hint-box{
  background:rgba(200,169,110,.06);
  border:1px dashed rgba(200,169,110,.2);
  border-radius:8px;
  padding:.65rem .9rem;
  font-size:.8rem;
  display:flex;
  align-items:center;
  gap:.5rem;
  flex-wrap:wrap;
}

.hint-label{
  color:#9a938c;
}

code{
  background:rgba(200,169,110,.12);
  color:#c8a96e;
  padding:.1rem .4rem;
  border-radius:4px;
}

/* button */
.btn-submit{
  margin-top:.5rem;
  padding:.8rem;
  background:linear-gradient(135deg,#c8a96e,#a0834e);
  color:#0d0d0f;
  border:none;
  border-radius:10px;
  font-weight:700;
  cursor:pointer;
  transition:opacity .2s, transform .15s;
}

.btn-submit:hover:not(:disabled){
  opacity:.9;
  transform:translateY(-1px);
}

.btn-submit:disabled{
  opacity:.6;
  cursor:not-allowed;
}

/* spinner */
.spinner{
  display:inline-block;
  width:18px;
  height:18px;
  border:2px solid rgba(13,13,15,.3);
  border-top-color:#0d0d0f;
  border-radius:50%;
  animation:spin .7s linear infinite;
}

/* animations */
@keyframes cardIn{
  from{opacity:0;transform:translateY(20px) scale(.98);}
  to{opacity:1;transform:translateY(0) scale(1);}
}

@keyframes pulse{
  0%,100%{opacity:1;}
  50%{opacity:.6;}
}

@keyframes spin{
  to{transform:rotate(360deg);}
}

/* transition */
.alert-enter-active,
.alert-leave-active{
  transition:all .3s ease;
}

.alert-enter-from,
.alert-leave-to{
  opacity:0;
  transform:translateY(-8px);
}

/* ---------------- MOBILE FIXES ---------------- */

@media (max-width:480px){

.login-card{
  padding:2rem 1.4rem;
}

.card-icon{
  font-size:1.7rem;
}

.bg-glow{
  width:320px;
  height:320px;
}

.field-input{
  padding:.65rem 2.4rem .65rem 2.3rem;
}

.hint-box{
  font-size:.75rem;
}

}

</style>