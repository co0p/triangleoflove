<template>
  <div class="page">
    <header class="container section">
      <h1>Welcome back, {{ firstName }}</h1>
    </header>
    <main class="container section">
      <div class="card dashboard-card">
        <div class="dashboard-status-row">
          <span class="dashboard-dot" :class="checkedIn ? 'dashboard-dot--done' : 'dashboard-dot--pending'" aria-hidden="true"></span>
          <p class="dashboard-status-label" :class="{ 'text-muted': !checkedIn }">
            {{ checkedIn ? 'Checked in today' : 'Not checked in yet' }}
          </p>
        </div>
        <router-link data-testid="checkin-link" to="/checkin" class="btn btn-primary">
          Daily check-in
        </router-link>
      </div>

      <div class="card dashboard-card">
        <div class="dashboard-status-row">
          <span class="dashboard-dot" :class="partnerName ? 'dashboard-dot--done' : 'dashboard-dot--pending'" aria-hidden="true"></span>
          <p v-if="partnerName" data-testid="pairing-status" class="dashboard-status-label">
            Connected with <strong>{{ partnerName }}</strong>
          </p>
          <p v-else class="dashboard-status-label text-muted">
            Not connected yet
          </p>
        </div>
        <router-link to="/pairing" class="btn btn-secondary">
          {{ partnerName ? 'View pairing' : 'Connect with partner' }}
        </router-link>
      </div>

      <div class="card dashboard-card">
        <router-link data-testid="insights-link" to="/insights" class="btn btn-secondary">
          Weekly insights
        </router-link>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { getMe } from '../api/users.js';
import { getCoupleStatus } from '../api/pairing.js';
import { getTodayCheckin } from '../api/checkin.js';

const firstName = ref('');
const partnerName = ref('');
const checkedIn = ref(false);
const router = useRouter();

onMounted(async () => {
  try {
    const [profile, status, todayCheckin] = await Promise.all([getMe(), getCoupleStatus(), getTodayCheckin()]);
    firstName.value = profile.firstName;
    if (status.paired) partnerName.value = status.partner_first_name;
    checkedIn.value = todayCheckin !== null;
  } catch (error) {
    if (!(error instanceof Error) || error.message !== 'unauthorized') {
      return;
    }
    localStorage.removeItem('token');
    router.push('/login');
  }
});
</script>

<style scoped>
.dashboard-card {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  margin-bottom: var(--space-6);
}

.dashboard-status-row {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.dashboard-dot {
  flex-shrink: 0;
  width: 0.625rem;
  height: 0.625rem;
  border-radius: var(--radius-full);
}

.dashboard-dot--done {
  background-color: var(--color-primary);
}

.dashboard-dot--pending {
  background-color: var(--color-neutral-400);
}

.dashboard-status-label {
  font-size: var(--font-size-sm);
  line-height: var(--line-height-normal);
}
</style>