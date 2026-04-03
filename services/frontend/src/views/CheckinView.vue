<template>
  <div class="checkin-page">
    <NavBar :firstName="firstName" />
    <main class="container section">
      <h1>Daily check-in</h1>
      <p class="text-muted">30–60 seconds · private note optional</p>
      <div class="checkin-dimensions">
        <div v-for="dim in relationshipDimensions" :key="dim.key" class="checkin-row">
          <label :for="dim.key" class="checkin-label">{{ dim.label }}</label>
          <p class="checkin-description">{{ dim.description }}</p>
          <div class="checkin-slider-row">
            <input
              :id="dim.key"
              type="range"
              min="-5"
              max="5"
              :value="ratings[dim.key] === null ? 0 : ratings[dim.key]"
              :data-unset="ratings[dim.key] === null ? 'true' : undefined"
              :class="['checkin-slider', { 'checkin-slider--unset': ratings[dim.key] === null }]"
              @input="onSliderInput(dim.key, $event)"
            />
            <span class="checkin-value-badge" :class="{ 'checkin-value-badge--unset': ratings[dim.key] === null }">
              {{ ratings[dim.key] === null ? '—' : (ratings[dim.key] > 0 ? '+' : '') + ratings[dim.key] }}
            </span>
          </div>
        </div>
      </div>
      <hr class="checkin-divider" />
      <div class="checkin-dimensions">
        <div v-for="dim in personalDimensions" :key="dim.key" class="checkin-row">
          <label :for="dim.key" class="checkin-label">{{ dim.label }}</label>
          <p class="checkin-description">{{ dim.description }}</p>
          <div class="checkin-slider-row">
            <input
              :id="dim.key"
              type="range"
              min="-5"
              max="5"
              :value="ratings[dim.key] === null ? 0 : ratings[dim.key]"
              :data-unset="ratings[dim.key] === null ? 'true' : undefined"
              :class="['checkin-slider', { 'checkin-slider--unset': ratings[dim.key] === null }]"
              @input="onSliderInput(dim.key, $event)"
            />
            <span class="checkin-value-badge" :class="{ 'checkin-value-badge--unset': ratings[dim.key] === null }">
              {{ ratings[dim.key] === null ? '—' : (ratings[dim.key] > 0 ? '+' : '') + ratings[dim.key] }}
            </span>
          </div>
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
  { key: 'felt_close',            label: 'Felt close today',          description: 'Did you feel emotionally connected?' },
  { key: 'positive_energy',       label: 'Positive energy / fun',     description: 'Was there lightness or joy between you?' },
  { key: 'supported',             label: 'Supported / team',          description: 'Did you feel like you had each other\'s backs?' },
  { key: 'communication_healthy', label: 'Communication healthy',     description: 'Were you able to talk openly and be heard?' },
  { key: 'stress_level',          label: 'My stress level',           description: 'How much did personal stress weigh on you?' },
];

const relationshipDimensions = dimensions.slice(0, 4);
const personalDimensions = dimensions.slice(4);

const ratings = reactive({
  felt_close: null,
  positive_energy: null,
  supported: null,
  communication_healthy: null,
  stress_level: null,
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
  const val = Number(event.target.value);
  if (ratings[key] === null && val === 0) return;
  ratings[key] = val;
}

async function save() {
  error.value = '';
  confirmed.value = false;
  const payload = {
    ...ratings,
    note: note.value,
  };
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

.checkin-slider-row {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.checkin-slider {
  flex: 1;
}

.checkin-value-badge {
  min-width: 2.5rem;
  text-align: center;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-text);
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  padding: 2px var(--space-2);
}

.checkin-value-badge--unset {
  color: var(--color-text-muted);
  border-color: transparent;
  background: transparent;
}

.checkin-slider--unset {
  opacity: 0.4;
}

.checkin-description {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
  margin: 0;
}

.checkin-divider {
  border: none;
  border-top: 1px solid var(--color-border);
  margin: var(--space-6) 0;
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
