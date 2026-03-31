<template>
  <div class="dashboard-page">
    <NavBar :firstName="firstName" />
    <header class="container section">
      <h1>Welcome back, {{ firstName }}</h1>
    </header>
    <main class="container section">
      <router-link data-testid="checkin-link" to="/checkin" class="btn btn-primary checkin-entry">
        Daily check-in
      </router-link>
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
