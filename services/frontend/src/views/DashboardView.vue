<template>
  <div class="page">
    <header class="container section">
      <h1>Welcome back, {{ firstName }}</h1>
    </header>
    <main class="container section">
      <section class="block">
        <CheckinMatrix :past="31" />
      </section>

      <div class="card dashboard-card">
        <div class="dashboard-status-row">
          <span class="dashboard-dot" :class="checkedIn ? 'dashboard-dot--done' : 'dashboard-dot--pending'" aria-hidden="true"></span>
          <p class="dashboard-status-label" :class="{ 'text-muted': !checkedIn }">
            {{ checkedIn ? 'Weekly check-in logged today' : 'No check-in logged today' }}
          </p>
        </div>
        <div class="dashboard-weekly-row">
          <p class="dashboard-weekly-label text-muted">Weekly rhythm</p>
          <WeeklyBacklogDots :weekly-insights="weeklyInsights" />
        </div>
        <router-link data-testid="session-link" to="/session" class="btn btn-primary">
          Weekly check-in
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
import { getCoupleStatus } from '../api/pairing.js';
import { getTodaySession } from '../api/checkin.js';
import { getWeeklyInsights } from '../api/insights.js';
import { useCurrentUser } from '../composables/useCurrentUser.js';
import CheckinMatrix from '../components/CheckinMatrix.vue';
import WeeklyBacklogDots from '../components/WeeklyBacklogDots.vue';

const { firstName, load } = useCurrentUser();
const partnerName = ref('');
const checkedIn = ref(false);
const weeklyInsights = ref([]);
const router = useRouter();

onMounted(async () => {
  try {
    const [, status, todaySession, insights] = await Promise.all([
      load(),
      getCoupleStatus(),
      getTodaySession(),
      getWeeklyInsights(7),
    ]);

    if (status.paired) partnerName.value = status.partner_first_name;
    checkedIn.value = todaySession !== null;
    weeklyInsights.value = insights;
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

.dashboard-weekly-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 44px;
}

.dashboard-weekly-label {
  font-size: var(--font-size-xs);
}
</style>
