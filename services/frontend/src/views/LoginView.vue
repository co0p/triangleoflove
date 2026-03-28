<template>
  <main>
    <h1>Sign in</h1>
    <form @submit.prevent="handleSubmit">
      <label>
        Email
        <input v-model="email" type="email" autocomplete="email" required />
      </label>
      <label>
        Password
        <input v-model="password" type="password" autocomplete="current-password" required />
      </label>
      <p v-if="error" role="alert">{{ error }}</p>
      <button type="submit">Continue</button>
    </form>
  </main>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { login } from '../api/auth.js';

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
