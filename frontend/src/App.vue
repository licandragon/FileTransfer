<template>
  <div id="app" class="min-h-screen bg-[#0d0d0f] text-[#e8e4dc] font-[DM Sans]">

    <!-- NAVBAR -->
    <nav class="flex items-center justify-between
  px-4 sm:px-6 lg:px-10
  py-3 sm:py-4
  bg-[rgba(13,13,15,0.85)] backdrop-blur-md
  border-b border-[rgba(255,255,255,0.06)]
  sticky top-0 z-100
  flex-wrap sm:flex-nowrap
">
      <!-- BRAND -->
      <div>
        <RouterLink to="/" class="flex items-center gap-2 font-bold text-[1.2rem] tracking-[-0.02em] no-underline">
          <span class="text-[#c8a96e] text-[1.4rem]">⬡</span>

          <span class="bg-linear-to-br from-[#e8e4dc] to-[#c8a96e]
            bg-clip-text text-transparent">
            VueAuth
          </span>
        </RouterLink>
      </div>

      <!-- LINKS -->
      <div class="flex items-center gap-1 flex-wrap sm:flex-nowrap">


        <RouterLink to="/" class="nav-link px-4 py-[0.45rem] rounded-md text-[0.9rem] font-medium
          text-[#a09a91] hover:text-[#e8e4dc]
          hover:bg-[rgba(255,255,255,0.07)]
          transition">
          Inicio
        </RouterLink>

        <!-- NO autenticado -->
        <template v-if="!isAuthenticated">
          <RouterLink to="/login" class="nav-link px-5 py-[0.45rem] rounded-md text-[0.9rem] font-semibold
            text-[#0d0d0f]
            bg-linear-to-br from-[#c8a96e] to-[#a0834e]
            hover:opacity-90 transition">
            Iniciar sesión
          </RouterLink>
          <RouterLink to="/register" class="nav-link px-5 py-[0.45rem] rounded-md text-[0.9rem] font-semibold
            text-[#0d0d0f]
            bg-linear-to-br from-[#c8a96e] to-[#a0834e]
            hover:opacity-90 transition">
            Registrarse
          </RouterLink>
        </template>

        <!-- autenticado -->
        <template v-else>

          <RouterLink to="/dashboard" class="nav-link px-4 py-[0.45rem] rounded-md text-[0.9rem] font-medium
            text-[#a09a91]
            hover:text-[#e8e4dc]
            hover:bg-[rgba(255,255,255,0.07)]
            transition">
            Dashboard
          </RouterLink>

          <span class="hidden sm:inline
  text-[0.85rem] text-[#c8a96e]
  px-3 border-l border-r
  border-[rgba(200,169,110,0.3)]">

            {{ user?.name }}
          </span>

          <button @click="handleLogout" class="px-4 py-[0.45rem] rounded-md text-[0.9rem] font-medium
            text-[#e07070]
            hover:bg-[rgba(224,112,112,0.1)]
            transition">
            Cerrar sesión
          </button>

        </template>

      </div>
    </nav>

    <!-- MAIN -->
    <main class="max-w-275 mx-auto px-4 sm:px-6 py-6 sm:py-8">

      <RouterView v-slot="{ Component }">
        <Transition name="page" mode="out-in">
          <component :is="Component" />
        </Transition>
      </RouterView>
    </main>

  </div>
</template>

<script setup>
import { RouterLink, RouterView, useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import { watch } from 'vue'



const router = useRouter()
const { isAuthenticated, user, logout } = useAuth()

async function handleLogout() {
  await logout()
  await router.replace('/login')
}

watch(isAuthenticated, (val) => {
  if (!val) {
    router.replace('/login')
  }
})
</script>

<style>
/* ACTIVE ROUTE (igual que el original) */

.nav-link.router-link-active,
.nav-link.router-link-exact-active {
  color: #e8e4dc;
  background: rgba(255, 255, 255, 0.07);
}

/* TRANSICIONES */

.page-enter-active,
.page-leave-active {
  transition: all 0.25s ease;
}

.page-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.page-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
