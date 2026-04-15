<template>
  <div class="flex flex-col item-center gap-20 pb-12">
          <!-- BACKGROUND -->
      <div class="absolute inset-0 pointer-events-none z-0">

        <div class="absolute w-100 h-100 rounded-full
          bg-[var(--primary)] blur-[80px] opacity-[0.12]
          -top-25 -left-25 animate-drift"></div>

        <div class="absolute w-75 h-75 rounded-full
          bg-[#6eb4e8] blur-[80px] opacity-[0.12]
          -bottom-20 right-25 animate-drift2"></div>

        <div class="absolute inset-0
          bg-[linear-gradient(rgba(200,169,110,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(200,169,110,0.03)_1px,transparent_1px)]
          bg-size-[40px_40px]"></div>

      </div>

    <!-- HERO -->
    <section class="relative
      gap-8 md:gap-12 items-center min-h-120 py-12 overflow-hidden">



      <!-- CODE CARD -->

      <Upload />
    </section>

    <!-- FEATURES -->
    <section
      class="grid gap-4 grid-cols-[repeat(auto-fit,minmax(180px,1fr))] sm:grid-cols-[repeat(auto-fit,minmax(220px,1fr))]">

      <div v-for="(feat, i) in features" :key="feat.title" class="p-6 bg-[var(--bg-card)]
        border border-[rgba(255,255,255,0.06)]
        rounded-[14px]
        flex flex-col gap-2
        hover:border-[var(--primary-border)]
        hover:-translate-y-0.75
        transition animate-slideUp" :style="{ animationDelay: `${i * 0.1}s` }">
        <div class="text-[1.75rem]">{{ feat.icon }}</div>
        <h3 class="font-bold text-[1rem]">{{ feat.title }}</h3>

        <p class="text-[0.85rem] text-[#5a5550] leading-[1.55] flex-1">
          {{ feat.desc }}
        </p>

        <div class="text-[0.72rem] font-mono text-[var(--primary)]
          bg-[var(--primary-soft)]
          px-2 py-[0.2rem] rounded-md w-fit">
          {{ feat.tag }}
        </div>
      </div>

    </section>

    <!-- FLOW -->
    <section class="text-center flex flex-col items-center gap-4">

      <h2 class="text-[1.5rem] font-bold tracking-[-0.03em]">
        Flujo de autenticación
      </h2>

      <p class="text-[#5a5550] text-[0.9rem]">
        Así viaja el usuario por la aplicación
      </p>

      <div class="flex items-center gap-2 flex-wrap justify-center mt-2">

        <div v-for="(step, i) in flowSteps" :key="step.label" class="flex items-center gap-2 relative">

          <div class="w-12 h-12 sm:w-14 sm:h-14 rounded-[14px]
            flex items-center justify-center text-[1.4rem]" :class="{
              'bg-[rgba(200,169,110,0.12)] border border-[rgba(200,169,110,0.25)]': step.color === 'gold',
              'bg-[rgba(110,180,232,0.12)] border border-[rgba(110,180,232,0.25)]': step.color === 'blue',
              'bg-[rgba(111,207,151,0.12)] border border-[rgba(111,207,151,0.25)]': step.color === 'green'
            }">
            {{ step.icon }}
          </div>

          <span class="absolute text-[0.75rem] text-[#5a5550] mt-18">
            {{ step.label }}
          </span>

          <div v-if="i < flowSteps.length - 1" class="text-[#3a3530] text-[1.2rem] mx-1">
            →
          </div>

        </div>

      </div>

    </section>

    <!-- CTA -->
    <section class="text-center px-8 py-12
      bg-[var(--bg-card)]
      border border-[rgba(200,169,110,0.1)]
      rounded-[20px]
      flex flex-col items-center gap-4">

      <h2 class="text-[1.8rem] font-extrabold tracking-[-0.03em]">
        ¿Listo para practicar?
      </h2>

      <p class="text-[#5a5550] text-[0.95rem] max-w-100 leading-[1.6]">
        Inicia sesión con las credenciales de prueba y explora
        la protección de rutas en acción.
      </p>

      <RouterLink v-if="!isAuthenticated" to="/login" class="px-8 py-[0.85rem] rounded-[10px]
        text-[1rem] font-semibold
        bg-linear-to-br from-[var(--primary)] to-[var(--primary-strong)]
        text-[#0d0d0f]
        hover:opacity-90 hover:-translate-y-0.5
        transition">
        Comenzar ahora →
      </RouterLink>

      <RouterLink v-else to="/dashboard" class="px-8 py-[0.85rem] rounded-[10px]
        text-[1rem] font-semibold
        bg-linear-to-br from-[var(--primary)] to-[var(--primary-strong)]
        text-[#0d0d0f]
        hover:opacity-90 hover:-translate-y-0.5
        transition">
        Ver mi Dashboard →
      </RouterLink>

    </section>

  </div>
</template>

<script setup>
import { RouterLink } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import Upload from '@/components/Upload.vue'
import CodeCard from '@/components/CodeCard.vue'

const { isAuthenticated } = useAuth()

const features = [
  {
    icon: '🔐',
    title: 'Login reactivo',
    desc: 'Formulario con validación en tiempo real usando reactive() y computed().',
    tag: 'LoginView.vue',
  },
  {
    icon: '🛡️',
    title: 'Rutas protegidas',
    desc: 'Navigation guards con beforeEach bloquean el acceso sin autenticación.',
    tag: 'router/index.js',
  },
  {
    icon: '⚡',
    title: 'Estado global',
    desc: 'useAuth() comparte el estado de sesión entre todos los componentes.',
    tag: 'useAuth.js',
  },
  {
    icon: '💾',
    title: 'Persistencia',
    desc: 'El token se guarda en localStorage para mantener la sesión activa.',
    tag: 'composable',
  },
]

const flowSteps = [
  { icon: '🏠', label: 'Home', color: 'gold' },
  { icon: '🔑', label: 'Login', color: 'blue' },
  { icon: '✅', label: 'Autenticado', color: 'green' },
  { icon: '📊', label: 'Dashboard', color: 'gold' },
]
</script>

<style scoped>
@keyframes drift {
  from {
    transform: translate(0, 0);
  }

  to {
    transform: translate(30px, 20px);
  }
}

@keyframes drift2 {
  from {
    transform: translate(0, 0);
  }

  to {
    transform: translate(-30px, -20px);
  }
}

@keyframes blink {

  0%,
  100% {
    opacity: 1
  }

  50% {
    opacity: .3
  }
}

@keyframes floatCard {

  0%,
  100% {
    transform: translateY(0)
  }

  50% {
    transform: translateY(-10px)
  }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(16px)
  }

  to {
    opacity: 1;
    transform: translateY(0)
  }
}

.animate-drift {
  animation: drift 8s ease-in-out infinite alternate;
}

.animate-drift2 {
  animation: drift 10s ease-in-out infinite alternate-reverse;
}

.animate-blink {
  animation: blink 2s ease-in-out infinite;
}

.animate-floatCard {
  animation: floatCard 6s ease-in-out infinite;
}

.animate-slideUp {
  animation: slideUp .4s ease both;
}
</style>
