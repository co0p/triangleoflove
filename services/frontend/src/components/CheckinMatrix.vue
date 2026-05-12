<template>
  <section data-testid="checkin-history" aria-label="Check-in history">
    <table class="checkin-matrix-table">
      <tbody>
        <tr v-for="dim in dimensions" :key="dim.key" data-testid="matrix-row">
          <td class="checkin-matrix-label">
            <span data-testid="dimension-icon" aria-hidden="true">
              <component :is="dim.icon" :size="14" :stroke-width="2" />
            </span>
          </td>
          <td
            v-for="day in paddedData"
            :key="day.date"
            data-testid="matrix-cell"
            :class="['checkin-matrix-cell', cellClass(day[dim.key], dim.key)]"
          >
            <router-link
              :to="`/insights/${day.date}`"
              data-testid="cell-link"
              class="checkin-matrix-link"
              aria-label="View insights"
            />
          </td>
        </tr>
      </tbody>
    </table>
  </section>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { Heart, Anchor, Flame } from 'lucide-vue-next';
import { getWeeklyInsights } from '../api/insights.js';

const props = defineProps({
  past: { type: Number, default: 7 },
});

const router = useRouter();
const rawData = ref([]);
const dimensions = [
  { key: 'intimacy',   icon: Heart  },
  { key: 'commitment', icon: Anchor },
  { key: 'passion',    icon: Flame  },
];

function dateWindowUTC(past) {
  const today = new Date();
  return Array.from({ length: past }, (_, i) => {
    const d = new Date(Date.UTC(
      today.getUTCFullYear(),
      today.getUTCMonth(),
      today.getUTCDate() - (past - 1 - i),
    ));
    return d.toISOString().slice(0, 10).replace(/-/g, '');
  });
}

const paddedData = computed(() => {
  const window = dateWindowUTC(props.past);
  return window.map(date => {
    const found = rawData.value.find(d => d.date === date);
    return found ?? { date, intimacy: -1, commitment: -1, passion: -1 };
  });
});

function cellClass(score, dimKey) {
  if (score < 0) return 'cell-unavailable';
  if (score <= 24) return `${dimKey}-very-low`;
  if (score <= 49) return `${dimKey}-low`;
  if (score <= 74) return `${dimKey}-moderate`;
  return `${dimKey}-high`;
}

onMounted(async () => {
  try {
    rawData.value = await getWeeklyInsights(props.past);
  } catch (error) {
    if (error instanceof Error && error.message === 'unauthorized') {
      localStorage.removeItem('token');
      router.push('/login');
    }
  }
});
</script>

<style scoped>
.checkin-matrix-table {
  width: 100%;
  border-collapse: separate;
  border-spacing: var(--space-1);
  table-layout: fixed;
}

.checkin-matrix-label {
  width: var(--space-8);
  vertical-align: middle;
  text-align: center;
  color: var(--color-text-muted);
  padding: 0;
}

.checkin-matrix-cell {
  height: var(--space-5);
  border-radius: var(--radius-sm);
  padding: 0;
  position: relative;
}

.checkin-matrix-link {
  display: block;
  position: absolute;
  inset: 0;
  border-radius: inherit;
}

.checkin-matrix-cell.intimacy-very-low    { background-color: var(--color-intimacy-very-low); }
.checkin-matrix-cell.intimacy-low         { background-color: var(--color-intimacy-low); }
.checkin-matrix-cell.intimacy-moderate    { background-color: var(--color-intimacy-moderate); }
.checkin-matrix-cell.intimacy-high        { background-color: var(--color-intimacy-high); }
.checkin-matrix-cell.commitment-very-low  { background-color: var(--color-commitment-very-low); }
.checkin-matrix-cell.commitment-low       { background-color: var(--color-commitment-low); }
.checkin-matrix-cell.commitment-moderate  { background-color: var(--color-commitment-moderate); }
.checkin-matrix-cell.commitment-high      { background-color: var(--color-commitment-high); }
.checkin-matrix-cell.passion-very-low     { background-color: var(--color-passion-very-low); }
.checkin-matrix-cell.passion-low          { background-color: var(--color-passion-low); }
.checkin-matrix-cell.passion-moderate     { background-color: var(--color-passion-moderate); }
.checkin-matrix-cell.passion-high         { background-color: var(--color-passion-high); }
.checkin-matrix-cell.cell-unavailable     { background-color: var(--color-neutral-200); }
</style>
