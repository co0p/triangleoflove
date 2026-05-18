<template>
  <main class="login-page">
    <div class="container">
      <div class="card login-card">
        <div class="login-logo">
          <img :src="logoSrc" alt="Triangle of Love" width="64" height="64" />
        </div>
        <h1 class="login-heading">Sign in</h1>
        <form @submit.prevent="handleSubmit" class="form-fields" novalidate>
          <div class="input-group">
            <label class="input-label" for="email">Email</label>
            <input
              id="email"
              v-model="email"
              class="input"
              type="email"
              autocomplete="email"
              placeholder="you@example.com"
              required
            />
          </div>
          <div class="input-group">
            <label class="input-label" for="password">Password</label>
            <input
              id="password"
              v-model="password"
              class="input"
              type="password"
              autocomplete="current-password"
              placeholder="••••••••"
              required
            />
          </div>
          <p v-if="successMessage" class="login-success" role="status">{{ successMessage }}</p>
          <div v-if="error" class="alert-error" role="alert">{{ error }}</div>
          <button type="submit" class="btn btn-primary">Sign in</button>
        </form>
        <p class="login-footer">
          Need an account?
          <router-link to="/register">Sign up</router-link>
        </p>
      </div>
    </div>
  </main>
</template>

<script setup>
import { computed, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { login } from '../api/auth.js';
import { useCurrentUser } from '../composables/useCurrentUser.js';
import logoSrc from '../assets/logo.svg';

const email = ref('');
const password = ref('');
const error = ref('');
const route = useRoute();
const router = useRouter();
const { load } = useCurrentUser();
const successMessage = computed(() =>
  route.query.registered === '1' ? 'Account created. You can now sign in.' : ''
);

async function handleSubmit() {
  error.value = '';
  try {
    const { token } = await login(email.value, password.value);
    localStorage.setItem('token', token);
    await load();
    router.push('/dashboard');
  } catch {
    error.value = 'Invalid email or password.';
  }
}
</script>

<style scoped>
.login-logo {
  display: flex;
  justify-content: center;
}

.login-footer {
  text-align: center;
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
}

.login-success {
  padding: var(--space-3) var(--space-4);
  border: 1px solid var(--color-border);
  border-left: 3px solid var(--color-accent);
  border-radius: var(--radius-md);
  background: var(--color-surface);
  color: var(--color-text);
  font-size: var(--font-size-sm);
  line-height: var(--line-height-normal);
}
</style>
