<template>
  <section data-testid="monthly-matrix" aria-label="Monthly insights">
    <table class="matrix-table">
      <tbody>
        <tr v-for="dimension in dimensions" :key="dimension" data-testid="matrix-row">
          <td
            v-for="day in paddedData"
            :key="day.date"
            data-testid="matrix-cell"
            :class="['matrix-cell', cellClass(day[dimension])]"
          ></td>
        </tr>
      </tbody>
    </table>
  </section>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { getWeeklyInsights } from '../api/insights.js';

const props = defineProps({
  past: { type: Number, default: 7 },
});

const router = useRouter();
const rawData = ref([]);
const dimensions = ['intimacy', 'commitment', 'passion'];

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

function cellClass(score) {
  if (score < 0) return 'insight-value--unavailable';
  if (score <= 24) return 'insight-value--very-low';
  if (score <= 49) return 'insight-value--low';
  if (score <= 74) return 'insight-value--moderate';
  return 'insight-value--high';
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
.matrix-table {
  width: 100%;
  border-collapse: collapse;
  table-layout: fixed;
}

.matrix-cell {
  height: var(--space-5);
  padding: 0;
}

.matrix-cell.insight-value--very-low  { background-color: var(--color-insight-very-low); }
.matrix-cell.insight-value--low       { background-color: var(--color-insight-low); }
.matrix-cell.insight-value--moderate  { background-color: var(--color-insight-moderate); }
.matrix-cell.insight-value--high      { background-color: var(--color-insight-high); }
.matrix-cell.insight-value--unavailable { background-color: var(--color-neutral-200); }
</style>
