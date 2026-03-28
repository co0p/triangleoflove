<template>
  <header>
    <span>Hello, {{ firstName }}</span>
  </header>
  <main>
    <h1>Dashboard</h1>
  </main>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { getMe } from '../api/users.js';

const firstName = ref('');
const router = useRouter();

onMounted(async () => {
  try {
    const profile = await getMe();
    firstName.value = profile.firstName;
  } catch {
    localStorage.removeItem('token');
    router.push('/login');
  }
});
</script>
