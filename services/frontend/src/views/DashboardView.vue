<template>
  <div class="page">
    <NavBar :firstName="firstName" />
    <header class="container section">
      <h1>Welcome back, {{ firstName }}</h1>
    </header>
    <main class="container section">
      <p v-if="partnerName" data-testid="pairing-status">You are paired with {{ partnerName }}</p>
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
import { getCoupleStatus } from '../api/pairing.js';
import NavBar from '../components/NavBar.vue';

const firstName = ref('');
const partnerName = ref('');
const router = useRouter();

onMounted(async () => {
  try {
    const [profile, status] = await Promise.all([getMe(), getCoupleStatus()]);
    firstName.value = profile.firstName;
    if (status.paired) partnerName.value = status.partner_first_name;
  } catch {
    localStorage.removeItem('token');
    router.push('/login');
  }
});
</script>

<style scoped>
</style>
