// src/composables/useAuth.js
/* import { ref, computed } from "vue";

// Estado GLOBAL compartido (fuera del composable)
const _user = ref(null)
const _token = ref(localStorage.getItem('auth_token') || null)

export function useAuth() {
  const isAuthenticated = computed(() => !!_token.value)

  const user = computed(() => _user.value)

  function login(userData, token) {
    _user.value = userData
    _token.value = token
    localStorage.setItem('auth_token', token)
  }

  function logout() {
    _user.value = null
    _token.value = null
    localStorage.removeItem('auth_token')
  }

  return {
    user,
    isAuthenticated,
    login,
    logout,
  }
} */

import { ref, computed } from "vue";
import { supabase } from "@/lib/supabaseClient";

const _user = ref(null);
const _loading = ref(true);

let _initialized = false;

function setupAuthListener() {
  if (_initialized) return;
  _initialized = true;

  supabase.auth.onAuthStateChange((_event, session) => {
    _user.value = session?.user ?? null;
    _loading.value = false;
  });
}

setupAuthListener();

async function initAuth() {
  const { data } = await supabase.auth.getSession();
  _user.value = data.session?.user ?? null;
  _loading.value = false;
}

export function useAuth() {
  const user = computed(() => _user.value);
  const isAuthenticated = computed(() => !!_user.value);
  const loading = computed(() => _loading.value);

  async function login(email, password) {
    const { data, error } = await supabase.auth.signInWithPassword({
      email,
      password,
    });

    if (error) throw error;
    return data;
  }

  async function register(name, email, password) {
    const { data, error } = await supabase.auth.signUp({
      email,
      password,
      options: {
        data: { full_name: name },
      },
    });

    if (error) throw error;
    return data;
  }

  async function logout() {
    _user.value = null;
    await supabase.auth.signOut();
  }

  async function loginDiscord() {
    const { data, error } = await supabase.auth.signInWithOAuth({
      provider: "discord",
      options: {
        // Optional: where to redirect after login
        redirectTo: "http://localhost:5173/dashboard",
      },
    });

    if (error) console.error("Error logging in:", error.message);
  }

  return {
    user,
    isAuthenticated,
    loading,
    login,
    loginDiscord,
    register,
    logout,
  };
}

export { initAuth };
