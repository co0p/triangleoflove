<template>
  <div class="page">
    <NavBar :firstName="firstName" />
    <header class="container section">
      <h1>Welcome back, {{ firstName }}</h1>
    </header>
    <main class="container section">
      <div class="card dashboard-pairing-card">
        <div class="dashboard-pairing-status">
          <span v-if="partnerName" class="dashboard-pairing-dot dashboard-pairing-dot--connected" aria-hidden="true"></span>
          <span v-else class="dashboard-pairing-dot dashboard-pairing-dot--unpaired" aria-hidden="true"></span>
          <p v-if="partnerName" data-testid="pairing-status" class="dashboard-pairing-label">
            Connected with <strong>{{ partnerName }}</strong>
          </p>
          <p v-else class="dashboard-pairing-label text-muted">
            Not connected yet
          </p>
        </div>
        <router-link to="/pairing" class="btn btn-secondary dashboard-pairing-link">
          {{ partnerName ? 'View pairing' : 'Connect with partner' }}
        </router-link>
      </div>
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
.dashboard-pairing-card {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  margin-bottom: var(--space-6);
}

.dashboard-pairing-status {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.dashboard-pairing-dot {
  flex-shrink: 0;
  width: 0.625rem;
  height: 0.625rem;
  border-radius: var(--radius-full);
}

.dashboard-pairing-dot--connected {
  background-color: var(--color-primary);
}

.dashboard-pairing-dot--unpaired {
  background-color: var(--color-neutral-400);
}

.dashboard-pairing-label {
  font-size: var(--font-size-sm);
  line-height: var(--line-height-normal);
}

.dashboard-pairing-link {
  align-self: stretch;
}

.checkin-entry {
  display: flex;
}
</style>
