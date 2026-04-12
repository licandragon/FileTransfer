<!-- src/views/DashboardView.vue -->
<template>
  <div class="dashboard flex flex-col gap-7">

    <!-- HEADER -->
    <section
      class="flex flex-wrap items-center justify-between gap-4
      p-5 sm:p-7 rounded-2xl
      bg-[rgba(22,22,26,0.8)]
      border border-[rgba(200,169,110,0.12)]"
    >
      <div>
        <p class="text-[0.78rem] uppercase tracking-widest text-[#c8a96e] mb-1">
          Panel de control
        </p>

        <h1 class="text-[1.6rem] sm:text-[1.85rem] font-bold tracking-[-0.03em] text-[#e8e4dc] mb-1">
          Hola,
          <span class="highlight">{{ userName }}</span> 👋
        </h1>

        <p class="text-[#5a5550] text-sm capitalize">
          {{ currentDate }} — Todo en orden por aquí.
        </p>
      </div>

      <div
        class="flex items-center gap-2
        px-4 py-2 rounded-full
        text-[0.82rem] text-[#c8a96e]
        bg-[rgba(200,169,110,0.08)]
        border border-[rgba(200,169,110,0.2)]"
      >
        <span class="text-base">⬡</span>
        <span>Sesión activa</span>
      </div>
    </section>

    <!-- STATS -->
    <section
      class="grid gap-4
      grid-cols-[repeat(auto-fit,minmax(150px,1fr))]
      sm:grid-cols-[repeat(auto-fit,minmax(180px,1fr))]"
    >
      <div
        v-for="(stat, i) in stats"
        :key="stat.label"
        class="stat-card flex items-center gap-4
        p-5 rounded-xl
        bg-[rgba(22,22,26,0.8)]
        border border-[rgba(255,255,255,0.06)]"
        :style="{ animationDelay: `${i * 0.08}s` }"
      >
        <div class="text-2xl">{{ stat.icon }}</div>

        <div class="flex flex-col flex-1">
          <span class="text-[1.3rem] sm:text-[1.4rem] font-bold text-[#e8e4dc] leading-none">
            {{ stat.value }}
          </span>
          <span class="text-[0.78rem] text-[#5a5550] mt-0.5">
            {{ stat.label }}
          </span>
        </div>

        <div
          class="text-[0.78rem] font-semibold px-2 py-0.5 rounded-md"
          :class="stat.trend > 0 ? 'trend-up' : 'trend-down'"
        >
          {{ stat.trend > 0 ? '↑' : '↓' }} {{ Math.abs(stat.trend) }}%
        </div>
      </div>
    </section>

    <!-- BODY -->
    <div class="grid gap-4 lg:grid-cols-[1fr_340px]">

      <!-- ACTIVIDAD -->
      <section class="panel">
        <div class="panel-header">
          <h2 class="panel-title">Actividad reciente</h2>
          <span class="panel-badge">{{ activity.length }} eventos</span>
        </div>

        <ul class="flex flex-col gap-3">
          <li
            v-for="(item, i) in activity"
            :key="i"
            class="activity-item"
            :style="{ animationDelay: `${i * 0.07}s` }"
          >
            <span class="activity-dot" :class="`dot--${item.type}`"></span>

            <div class="flex flex-col gap-0.5">
              <span class="text-[0.875rem] text-[#c8c2ba]">
                {{ item.msg }}
              </span>
              <span class="text-[0.75rem] text-[#4a4540]">
                {{ item.time }}
              </span>
            </div>
          </li>
        </ul>
      </section>

      <!-- SESIÓN -->
      <section class="panel">
        <div class="panel-header">
          <h2 class="panel-title">Tu sesión</h2>
        </div>

        <div class="flex flex-col gap-2 mb-6">
          <div class="session-row">
            <span class="session-key">Usuario</span>
            <span class="session-val">{{ userName }}</span> 
          </div>

          <div class="session-row">
            <span class="session-key">Email</span>
            <span class="session-val">{{ user?.email }}</span>
          </div>

          <div class="session-row">
            <span class="session-key">Token</span>
            <span class="session-val token-val">{{ maskedToken }}</span>
          </div>

          <div class="session-row">
            <span class="session-key">Estado</span>
            <span class="status-active">● Autenticado</span>
          </div>
        </div>

        <button
          class="w-full py-[0.7rem] rounded-lg
          text-[#e07070] font-semibold text-[0.9rem]
          bg-[rgba(224,112,112,0.08)]
          border border-[rgba(224,112,112,0.25)]
          transition hover:bg-[rgba(224,112,112,0.15)]"
          @click="handleLogout"
        >
          Cerrar sesión
        </button>
      </section>

    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const { user, logout } = useAuth()
const userName = computed(() =>
  user.value?.user_metadata?.full_name || 'Usuario'
)

const currentDate = computed(() =>
  new Date().toLocaleDateString('es-MX', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
)

const maskedToken = computed(() => {
  const token = localStorage.getItem('auth_token') || ''
  return token.slice(0, 10) + '••••••••••••'
})

const stats = [
  { icon: '📦', label: 'Proyectos', value: '12', trend: +8 },
  { icon: '✅', label: 'Tareas', value: '48', trend: +15 },
  { icon: '👥', label: 'Usuarios', value: '7', trend: -2 },
  { icon: '📈', label: 'Rendimiento', value: '94%', trend: +3 },
]

const activity = [
  { type: 'success', msg: 'Inicio de sesión exitoso', time: 'Justo ahora' },
  { type: 'info', msg: 'Perfil actualizado', time: 'Hace 2 horas' },
  { type: 'warning', msg: 'Intento de acceso denegado', time: 'Ayer, 11:30pm' },
  { type: 'success', msg: 'Nueva tarea completada', time: 'Hace 2 días' },
  { type: 'info', msg: 'Reporte mensual generado', time: 'Hace 3 días' },
]

async function handleLogout() {
  await logout()
  router.push('/login')
}
</script>

<style scoped>

.dashboard {
  animation: fadeIn 0.4s ease both;
}

@keyframes fadeIn {
  from { opacity:0; transform:translateY(12px);}
  to { opacity:1; transform:translateY(0);}
}

.panel {
  background: rgba(22,22,26,0.8);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 14px;
  padding: 1.25rem;
}

.panel-header{
  display:flex;
  justify-content:space-between;
  align-items:center;
  margin-bottom:1.25rem;
}

.panel-title{
  font-size:.95rem;
  font-weight:600;
  color:#e8e4dc;
}

.panel-badge{
  font-size:.75rem;
  color:#c8a96e;
  background:rgba(200,169,110,0.1);
  padding:.2rem .6rem;
  border-radius:999px;
}

.stat-card{
  animation:slideUp .4s ease both;
  transition:border-color .2s, transform .2s;
}

.stat-card:hover{
  border-color:rgba(200,169,110,0.2);
  transform:translateY(-2px);
}

@keyframes slideUp{
  from{opacity:0; transform:translateY(16px);}
  to{opacity:1; transform:translateY(0);}
}

.activity-item{
  display:flex;
  gap:.85rem;
  animation:slideUp .35s ease both;
}

.activity-dot{
  width:8px;
  height:8px;
  border-radius:50%;
  margin-top:5px;
  flex-shrink:0;
}

.dot--success{ background:#6fcf97; box-shadow:0 0 6px #6fcf9788;}
.dot--info{ background:#6eb4e8; box-shadow:0 0 6px #6eb4e888;}
.dot--warning{ background:#e8c46e; box-shadow:0 0 6px #e8c46e88;}

.session-row{
  display:flex;
  justify-content:space-between;
  padding:.55rem 0;
  border-bottom:1px solid rgba(255,255,255,0.04);
  font-size:.875rem;
}

.session-key{ color:#5a5550;}
.session-val{ color:#c8c2ba; font-weight:500;}

.token-val{
  font-family:monospace;
  font-size:.8rem;
  color:#c8a96e;
  word-break:break-all;
}

.status-active{
  color:#6fcf97;
  font-size:.82rem;
}

.trend-up{
  color:#6fcf97;
  background:rgba(111,207,151,0.1);
}

.trend-down{
  color:#e07070;
  background:rgba(224,112,112,0.1);
}

.highlight{
  background:linear-gradient(135deg,#c8a96e,#e8d5a3);
  -webkit-background-clip:text;
  background-clip:text;
  -webkit-text-fill-color:transparent;
}

</style>