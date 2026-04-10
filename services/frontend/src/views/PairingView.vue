<template>
  <div class="page">
    <NavBar />
    <main class="container section">
      <h1>Pairing</h1>

      <div v-if="pairedWith" data-testid="paired-status">
        <p>You are paired with <span data-testid="partner-name">{{ pairedWith }}</span></p>
        <p data-testid="paired-since" class="paired-since">Since {{ pairedSinceFormatted }}</p>
      </div>

      <template v-else>
        <div class="card pairing-card">
          <div class="input-group">
            <p class="input-label">Your invite code</p>
            <p data-testid="invite-code" class="invite-code">{{ inviteCode }}</p>
          </div>
          <button class="btn btn-secondary" @click="regenerate">Regenerate</button>
        </div>

        <div class="card pairing-card">
          <div class="input-group">
            <label class="input-label" for="partner-code">Partner's code</label>
            <input
              id="partner-code"
              data-testid="partner-code-input"
              type="text"
              class="input"
              maxlength="6"
              placeholder="Enter 6-character code"
              v-model="partnerCode"
            />
          </div>
          <button
            data-testid="connect-button"
            class="btn btn-primary"
            @click="connect"
          >Connect</button>
          <p v-if="connectError" role="alert" data-testid="connect-error" class="error-text">{{ connectError }}</p>
          <p v-if="connectSuccess" data-testid="connect-success" class="success-text">Connected!</p>
        </div>
      </template>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { getPairing, getCoupleStatus, regeneratePairing, connectPairing } from '../api/pairing.js';
import NavBar from '../components/NavBar.vue';

const inviteCode = ref('');
const partnerCode = ref('');
const connectError = ref('');
const connectSuccess = ref(false);
const pairedWith = ref('');
const pairedSince = ref('');
const router = useRouter();

const pairedSinceFormatted = computed(() => {
  if (!pairedSince.value) return '';
  return new Date(pairedSince.value).toLocaleDateString(undefined, { year: 'numeric', month: 'long', day: 'numeric' });
});

onMounted(async () => {
  try {
    const status = await getCoupleStatus();
    if (status.paired) {
      pairedWith.value = status.partner_first_name;
      pairedSince.value = status.paired_since;
      return;
    }
    const data = await getPairing();
    inviteCode.value = data.invite_code;
  } catch {
    localStorage.removeItem('token');
    router.push('/login');
  }
});

async function regenerate() {
  const data = await regeneratePairing();
  inviteCode.value = data.invite_code;
}

async function connect() {
  connectError.value = '';
  connectSuccess.value = false;
  try {
    await connectPairing(partnerCode.value);
    connectSuccess.value = true;
  } catch (err) {
    connectError.value = err.message;
  }
}
</script>

<style scoped>
.pairing-card {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  margin-bottom: var(--space-6);
}

.invite-code {
  font-size: 2rem;
  font-weight: bold;
  letter-spacing: 0.25em;
  padding: var(--space-2) 0;
}

.paired-since {
  color: var(--color-text-muted);
}
</style>
