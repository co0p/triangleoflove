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
          <div v-if="error" class="alert-error" role="alert">{{ error }}</div>
          <button type="submit" class="btn btn-primary">Sign in</button>
        </form>
      </div>
    </div>
  </main>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { login } from '../api/auth.js';
import logoSrc from '../assets/logo.svg';

const email = ref('');
const password = ref('');
const error = ref('');
const router = useRouter();

async function handleSubmit() {
  error.value = '';
  try {
    const { token } = await login(email.value, password.value);
    localStorage.setItem('token', token);
    router.push('/dashboard');
  } catch {
    error.value = 'Invalid email or password.';
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  /* align-items: flex-start keeps the card near the top when the soft keyboard
     opens and shrinks the viewport — prevents the form being cut off mid-field. */
  align-items: flex-start;
  padding-block: var(--space-10);
  background-color: var(--color-bg);
}

.login-card {
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
}

.login-logo {
  display: flex;
  justify-content: center;
}

.login-heading {
  text-align: center;
}
</style>
