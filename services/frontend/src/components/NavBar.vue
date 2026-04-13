<template>
  <nav class="navbar">
    <div class="navbar-brand">
      <router-link to="/dashboard" aria-label="Go to dashboard">
        <img :src="logoSrc" alt="Triangle of Love" width="32" height="32" />
      </router-link>
    </div>
    <router-link to="/profile" class="navbar-profile-link" aria-label="Go to profile">
      <span class="navbar-greeting navbar-greeting--truncate" v-if="firstName">Hello, {{ firstName }}</span>

      <div class="avatar">{{ initials }}</div>
    </router-link>
  </nav>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import logoSrc from '../assets/logo.svg';
import { getMe } from '../api/users.js';

const firstName = ref('');

onMounted(async () => {
  try {
    const profile = await getMe();
    firstName.value = profile.firstName;
  } catch {
    // token invalid or missing — router guard will redirect
  }
});

const initials = computed(() =>
  firstName.value ? firstName.value.charAt(0).toUpperCase() : '?'
);
</script>
