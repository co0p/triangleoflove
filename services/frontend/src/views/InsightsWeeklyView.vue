<template>
  <div class="page">
    <main class="container section">
      <table>
        <tbody>
          <tr
            v-for="dimension in dimensions"
            :key="dimension"
            data-testid="weekly-row"
          >
            <td
              v-for="day in weeklyData"
              :key="day.date"
              data-testid="weekly-cell"
              :class="['weekly-cell', cellClass(day[dimension])]"
            >
              {{ day[dimension] >= 0 ? day[dimension] : '' }}
            </td>
          </tr>
        </tbody>
      </table>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { getWeeklyInsights } from '../api/insights.js';

const router = useRouter();
const weeklyData = ref([]);
const dimensions = ['intimacy', 'commitment', 'passion'];

function cellClass(score) {
  if (score < 0) return 'insight-value--unavailable';
  if (score <= 24) return 'insight-value--very-low';
  if (score <= 49) return 'insight-value--low';
  if (score <= 74) return 'insight-value--moderate';
  return 'insight-value--high';
}

onMounted(async () => {
  try {
    const data = await getWeeklyInsights(7);
    weeklyData.value = [...data].sort((a, b) => a.date.localeCompare(b.date));
  } catch (error) {
    if (error instanceof Error && error.message === 'unauthorized') {
      localStorage.removeItem('token');
      router.push('/login');
    }
  }
});
</script>
