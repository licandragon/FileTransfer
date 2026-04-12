<!-- src/views/NotFoundView.vue -->
<template>
  <div
    class="relative min-h-[calc(100vh-65px)] flex items-center justify-center overflow-hidden px-4 py-8"
  >

    <!-- Fondo -->
    <div class="absolute inset-0 pointer-events-none">
      <div class="bg-orb"></div>
      <div class="bg-grid"></div>
    </div>

    <!-- 404 gigante -->
    <span
      aria-hidden="true"
      class="ghost-number"
    >
      404
    </span>

    <!-- Contenido -->
    <div
      class="relative z-10 flex flex-col items-center text-center gap-5 fade-up"
    >

      <!-- icono -->
      <div
        class="w-16 h-16 sm:w-18 sm:h-18 rounded-[20px]
        bg-[rgba(200,169,110,0.08)]
        border border-[rgba(200,169,110,0.2)]
        flex items-center justify-center icon-float"
      >
        <span class="text-[1.8rem] text-[#c8a96e]">⬡</span>
      </div>

      <!-- titulo -->
      <h1
        class="text-[clamp(1.5rem,3vw,2.2rem)]
        font-extrabold tracking-[-0.03em]"
      >
        Página no encontrada
      </h1>

      <!-- descripción -->
      <p
        class="text-[#5a5550] text-[0.95rem]
        flex items-center gap-2 flex-wrap justify-center
        max-w-md"
      >
        La ruta

        <code
          class="bg-[rgba(200,169,110,0.1)]
          border border-[rgba(200,169,110,0.2)]
          text-[#c8a96e]
          px-2 py-0.5
          rounded text-[0.85rem] font-mono"
        >
          {{ currentPath }}
        </code>

        no existe en esta aplicación.
      </p>

      <!-- botones -->
      <div class="flex gap-3 flex-wrap justify-center">

        <RouterLink
          to="/"
          class="inline-flex items-center
          px-6 py-[0.65rem]
          rounded-[10px]
          text-[0.9rem] font-semibold
          text-[#0d0d0f]
          bg-linear-to-br from-[#c8a96e] to-[#a0834e]
          hover:opacity-90 hover:-translate-y-0.5
          transition"
        >
          ← Volver al inicio
        </RouterLink>

        <RouterLink
          v-if="isAuthenticated"
          to="/dashboard"
          class="inline-flex items-center
          px-6 py-[0.65rem]
          rounded-[10px]
          text-[0.9rem] font-semibold
          bg-[rgba(255,255,255,0.04)]
          border border-[rgba(255,255,255,0.08)]
          text-[#9a938c]
          hover:bg-[rgba(255,255,255,0.08)]
          hover:text-[#e8e4dc]
          transition"
        >
          Ir al Dashboard
        </RouterLink>

        <RouterLink
          v-else
          to="/login"
          class="inline-flex items-center
          px-6 py-[0.65rem]
          rounded-[10px]
          text-[0.9rem] font-semibold
          bg-[rgba(255,255,255,0.04)]
          border border-[rgba(255,255,255,0.08)]
          text-[#9a938c]
          hover:bg-[rgba(255,255,255,0.08)]
          hover:text-[#e8e4dc]
          transition"
        >
          Iniciar sesión
        </RouterLink>

      </div>

      <!-- rutas -->
      <div class="flex flex-col items-center gap-3 mt-1 max-w-lg">

        <p
          class="text-[0.78rem]
          text-[#3a3530]
          uppercase tracking-[0.08em]"
        >
          Rutas disponibles:
        </p>

        <div class="flex gap-2 flex-wrap justify-center">

          <RouterLink
            v-for="route in availableRoutes"
            :key="route.path"
            :to="route.path"
            class="inline-flex items-center gap-1
            px-3 py-[0.35rem]
            rounded-full
            bg-[rgba(255,255,255,0.03)]
            border border-[rgba(255,255,255,0.07)]
            text-[0.8rem]
            font-mono
            text-[#6b6560]
            hover:border-[rgba(200,169,110,0.25)]
            hover:text-[#c8a96e]
            hover:bg-[rgba(200,169,110,0.06)]
            transition"
          >
            <span>{{ route.icon }}</span>
            {{ route.path }}
          </RouterLink>

        </div>

      </div>

    </div>

  </div>
</template>

<script setup>
import { computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

const route = useRoute()
const { isAuthenticated } = useAuth()

const currentPath = computed(() => route.path)

const availableRoutes = computed(() => {
  const base = [
    { path: '/', icon: '🏠' },
    { path: '/login', icon: '🔑' }
  ]

  if (isAuthenticated.value) {
    base.push({ path: '/dashboard', icon: '📊' })
  }

  return base
})
</script>

<style scoped>
.bg-orb{
position:absolute;
width:600px;
height:600px;
background:radial-gradient(circle,rgba(200,169,110,0.06)0%,transparent 70%);
top:50%;
left:50%;
transform:translate(-50%,-50%);
}

.bg-grid{
position:absolute;
inset:0;
background-image:
linear-gradient(rgba(200,169,110,0.03)1px,transparent 1px),
linear-gradient(90deg,rgba(200,169,110,0.03)1px,transparent 1px);
background-size:40px 40px;
}

.ghost-number{
position:absolute;
font-size:clamp(160px,28vw,320px);
font-weight:900;
letter-spacing:-0.05em;
color:transparent;
-webkit-text-stroke:1px rgba(200,169,110,0.06);
top:50%;
left:50%;
transform:translate(-50%,-50%);
user-select:none;
white-space:nowrap;
animation:ghostPulse 4s ease-in-out infinite;
}

@keyframes ghostPulse{
0%,100%{opacity:1}
50%{opacity:0.4}
}

.icon-float{
animation:iconFloat 3s ease-in-out infinite;
}

@keyframes iconFloat{
0%,100%{transform:translateY(0)}
50%{transform:translateY(-6px)}
}

.fade-up{
animation:fadeUp .5s ease both;
}

@keyframes fadeUp{
from{opacity:0;transform:translateY(20px)}
to{opacity:1;transform:translateY(0)}
}
</style>