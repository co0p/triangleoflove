<template>
  <main class="login-page">
    <div class="container">
      <div class="card login-card">
        <h1 class="login-heading">Create account</h1>
        <form @submit.prevent="handleSubmit" class="form-fields" novalidate>
          <div class="input-group">
            <label class="input-label" for="firstName">First name</label>
            <input
              id="firstName"
              name="firstName"
              v-model="firstName"
              class="input"
              type="text"
              autocomplete="given-name"
              placeholder="Alice"
              required
            />
          </div>
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
              autocomplete="new-password"
              placeholder="••••••••"
              required
            />
          </div>
          <div v-if="error" class="alert-error" role="alert">{{ error }}</div>
          <button type="submit" class="btn btn-primary">Create account</button>
        </form>
        <p class="login-footer">
          Already have an account?
          <router-link to="/login">Sign in</router-link>
        </p>
      </div>
    </div>
  </main>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { register } from '../api/auth.js';

const email = ref('');
const password = ref('');
const firstName = ref('');
const error = ref('');
const router = useRouter();

async function handleSubmit() {
  error.value = '';

  if (password.value.length < 8) {
    error.value = 'Password must be at least 8 characters.';
    return;
  }

  try {
    await register(email.value, password.value, firstName.value);
    router.push('/login');
  } catch (err) {
    if (err.message === 'duplicate email') {
      error.value = 'An account with that email already exists.';
    } else {
      error.value = 'Registration failed. Please try again.';
    }
  }
}
</script>

<style scoped>
.login-footer {
  text-align: center;
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
}
</style>
