<template>
  <div class="page">
    <main class="container section">
      <header class="insights-header">
        <h1>Daily insights</h1>
        <p class="text-muted">{{ formattedDate }}</p>
      </header>

      <section v-if="insights" class="card insights-card" aria-label="Daily insight scores">
        <article class="insight-row">
          <div>
            <p class="insight-label">Intimacy</p>
            <p class="text-muted">Connection and shared understanding</p>
          </div>
          <p
            data-testid="insight-intimacy"
            :class="['insight-value', insightValueClass(insights.intimacy)]"
          >
            {{ insights.intimacy }}
          </p>
        </article>

        <article class="insight-row">
          <div>
            <p class="insight-label">Commitment</p>
            <p class="text-muted">Reliability and effort for us</p>
          </div>
          <p
            data-testid="insight-commitment"
            :class="['insight-value', insightValueClass(insights.commitment)]"
          >
            {{ insights.commitment }}
          </p>
        </article>

        <article class="insight-row">
          <div>
            <p class="insight-label">Passion</p>
            <p class="text-muted">Desire and spark</p>
          </div>
          <p
            data-testid="insight-passion"
            :class="['insight-value', insightValueClass(insights.passion)]"
          >
            {{ insights.passion }}
          </p>
        </article>
      </section>

      <section v-else-if="loadError" class="card insights-card insight-empty" role="alert" aria-label="Insights status">
        <p>{{ stateMessage }}</p>
      </section>
      <p v-else class="text-muted">Loading insights...</p>
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { getInsights } from '../api/insights.js';

const route = useRoute();
const router = useRouter();

const insights = ref(null);
const loadError = ref('');

const stateMessage = computed(() => {
  if (loadError.value === 'not_found') {
    return 'No insights available for this date.';
  }
  if (loadError.value === 'invalid_date') {
    return 'Invalid date. Use YYYYMMDD.';
  }
  if (loadError.value) {
    return 'Unable to load insights.';
  }
  return '';
});

const formattedDate = computed(() => {
  const date = route.params.date;
  if (typeof date !== 'string' || date.length !== 8) {
    return '';
  }

  return `${date.slice(0, 4)}-${date.slice(4, 6)}-${date.slice(6, 8)}`;
});

function insightValueClass(score) {
  if (score < 0) {
    return 'insight-value--unavailable';
  }
  if (score <= 24) {
    return 'insight-value--very-low';
  }
  if (score <= 49) {
    return 'insight-value--low';
  }
  if (score <= 74) {
    return 'insight-value--moderate';
  }

  return 'insight-value--high';
}

onMounted(async () => {
  try {
    insights.value = await getInsights(route.params.date);
  } catch (error) {
    if (error instanceof Error && error.message === 'unauthorized') {
      localStorage.removeItem('token');
      router.push('/login');
      return;
    }

    loadError.value = error instanceof Error ? error.message : 'failed';
  }
});
</script>

<style scoped>
.insights-header {
  margin-bottom: var(--space-6);
}

.insights-card {
  display: grid;
  gap: var(--space-4);
}

.insight-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-4);
  padding-bottom: var(--space-4);
  border-bottom: 1px solid var(--color-border);
}

.insight-row:last-child {
  padding-bottom: 0;
  border-bottom: 0;
}

.insight-label {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
}

.insight-value {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  line-height: var(--line-height-tight);
}

.insight-value--very-low {
  color: var(--color-insight-very-low);
}

.insight-value--low {
  color: var(--color-insight-low);
}

.insight-value--moderate {
  color: var(--color-insight-moderate);
}

.insight-value--high {
  color: var(--color-insight-high);
}

.insight-value--unavailable {
  color: var(--color-insight-unavailable);
}

.insight-empty {
  padding: var(--space-6);
  text-align: center;
}
</style>