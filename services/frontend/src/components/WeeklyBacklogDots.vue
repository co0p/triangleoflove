<template>
  <div
    class="weekly-backlog"
    role="img"
    :aria-label="`Weekly check-in backlog: ${filledDots} of 3 dots`"
    data-testid="weekly-backlog-dots"
  >
    <span
      v-for="index in 3"
      :key="index"
      class="weekly-backlog-dot"
      :class="{ 'weekly-backlog-dot--filled': index <= filledDots }"
      data-testid="weekly-backlog-dot"
      aria-hidden="true"
    />
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  weeklyInsights: { type: Array, required: true },
  referenceDate: { type: Date, default: () => new Date() },
});

const DAY_MS = 86400000;

function parseCompactUTCDate(date) {
  const year = Number(date.slice(0, 4));
  const month = Number(date.slice(4, 6));
  const day = Number(date.slice(6, 8));
  return new Date(Date.UTC(year, month - 1, day));
}

function toUTCStart(date) {
  return new Date(Date.UTC(date.getUTCFullYear(), date.getUTCMonth(), date.getUTCDate()));
}

function mapDaysBehindToDots(daysBehind) {
  if (daysBehind <= 1) return 1;
  if (daysBehind <= 3) return 2;
  return 3;
}

const filledDots = computed(() => {
  if (!props.weeklyInsights.length) {
    return 3;
  }

  const latestDate = props.weeklyInsights
    .map(entry => entry.date)
    .sort((a, b) => b.localeCompare(a))[0];
  const latestUTC = parseCompactUTCDate(latestDate);
  const todayUTC = toUTCStart(props.referenceDate);
  const daysBehind = Math.max(0, Math.floor((todayUTC - latestUTC) / DAY_MS));

  return mapDaysBehindToDots(daysBehind);
});
</script>

<style scoped>
.weekly-backlog {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
}

.weekly-backlog-dot {
  width: 0.5rem;
  height: 0.5rem;
  border-radius: var(--radius-full);
  background-color: var(--color-neutral-300);
}

.weekly-backlog-dot--filled {
  background-color: var(--color-neutral-500);
}
</style>
