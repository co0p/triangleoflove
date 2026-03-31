<template>
  <div class="checkin-page">
    <NavBar :firstName="firstName" />
    <main class="container section">
      <h1>Daily check-in</h1>
      <p class="text-muted">30–60 seconds · private note optional</p>
      <div class="checkin-dimensions">
        <div v-for="dim in dimensions" :key="dim.key" class="checkin-row">
          <label :for="dim.key" class="checkin-label">{{ dim.label }}</label>
          <input
            :id="dim.key"
            type="range"
            min="1"
            max="10"
            :value="ratings[dim.key] === -1 ? 5 : ratings[dim.key]"
            :data-unset="ratings[dim.key] === -1 ? 'true' : undefined"
            :class="['checkin-slider', { 'checkin-slider--unset': ratings[dim.key] === -1 }]"
            @input="onSliderInput(dim.key, $event)"
          />
        </div>
      </div>
      <textarea
        v-model="note"
        class="checkin-note"
        placeholder="Optional private note"
      />
      <button data-testid="save-checkin" class="btn btn-primary checkin-save" @click="save">
        Save Check-in
      </button>
      <p v-if="confirmed" data-testid="checkin-confirmation" class="checkin-confirmed">
        Check-in saved.
      </p>
      <p v-if="error" role="alert" class="checkin-error">
        {{ error }}
      </p>
    </main>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import NavBar from '../components/NavBar.vue';
import { getTodayCheckin, saveTodayCheckin } from '../api/checkin.js';

const firstName = '';

const dimensions = [
  { key: 'felt_close',            label: 'Felt close today' },
  { key: 'positive_energy',       label: 'Positive energy / fun' },
  { key: 'supported',             label: 'Supported / team' },
  { key: 'communication_healthy', label: 'Communication healthy' },
  { key: 'stress_level',          label: 'My stress level' },
];

const ratings = reactive({
  felt_close: -1,
  positive_energy: -1,
  supported: -1,
  communication_healthy: -1,
  stress_level: -1,
});

const note = ref('');
const confirmed = ref(false);
const error = ref('');

onMounted(async () => {
  const existing = await getTodayCheckin();
  if (existing) {
    Object.keys(ratings).forEach(k => { ratings[k] = existing[k]; });
    note.value = existing.note;
  }
});

function onSliderInput(key, event) {
  ratings[key] = Number(event.target.value);
}

async function save() {
  error.value = '';
  confirmed.value = false;
  const payload = {
    ...ratings,
    note: note.value,
  };
  // Send -1 for any unset slider as-is; backend stores it
  try {
    await saveTodayCheckin(payload);
    confirmed.value = true;
  } catch {
    error.value = 'Could not save check-in. Please try again.';
  }
}
</script>

<style scoped>
.checkin-page {
  min-height: 100vh;
  background-color: var(--color-bg);
}

.checkin-dimensions {
  margin-top: var(--space-6);
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
}

.checkin-row {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.checkin-label {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-text);
}

.checkin-slider--unset {
  opacity: 0.4;
}

.checkin-note {
  margin-top: var(--space-6);
  width: 100%;
  min-height: 6rem;
  border: 2px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: var(--space-3);
  font-family: var(--font-sans);
  font-size: var(--font-size-sm);
  color: var(--color-text);
  background: var(--color-surface);
  resize: none;
  box-sizing: border-box;
}

.checkin-save {
  margin-top: var(--space-6);
  width: 100%;
}

.checkin-confirmed {
  margin-top: var(--space-4);
  color: var(--color-primary-dark);
  font-size: var(--font-size-sm);
  text-align: center;
}

.checkin-error {
  margin-top: var(--space-4);
  color: var(--color-error);
  font-size: var(--font-size-sm);
  text-align: center;
}
</style>
