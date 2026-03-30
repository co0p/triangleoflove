<template>
  <div class="dashboard-page">
    <NavBar :firstName="firstName" />
    <main class="container section">
      <h1>Dashboard</h1>
      <p class="text-muted">Welcome back. More here soon.</p>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { getMe } from '../api/users.js';
import NavBar from '../components/NavBar.vue';

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

<style scoped>
.dashboard-page {
  min-height: 100vh;
  background-color: var(--color-bg);
}
</style>
