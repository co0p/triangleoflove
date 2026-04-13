<template>
  <div class="page">
    <main class="container section">
      <h1>Profile</h1>

      <!-- Security card -->
      <div class="card profile-card">
        <h2>Security</h2>
        <form class="form-fields" @submit.prevent="submitChangePassword">
          <div class="input-group">
            <label class="input-label" for="current_password">Current password</label>
            <input
              id="current_password"
              name="current_password"
              type="password"
              class="input"
              v-model="currentPassword"
              autocomplete="current-password"
            />
          </div>
          <div class="input-group">
            <label class="input-label" for="new_password">New password</label>
            <input
              id="new_password"
              name="new_password"
              type="password"
              class="input"
              v-model="newPassword"
              autocomplete="new-password"
            />
          </div>
          <div class="input-group">
            <label class="input-label" for="confirm_password">Confirm new password</label>
            <input
              id="confirm_password"
              name="confirm_password"
              type="password"
              class="input"
              v-model="confirmPassword"
              autocomplete="new-password"
            />
          </div>
          <button type="submit" class="btn btn-primary" :disabled="loading">
            {{ loading ? 'Saving…' : 'Change password' }}
          </button>
          <p v-if="passwordError" class="alert-error" role="alert">{{ passwordError }}</p>
          <p v-if="passwordSuccess" class="password-success">{{ passwordSuccess }}</p>
        </form>
      </div>

      <!-- Account card -->
      <div class="card profile-card">
        <h2>Account</h2>
        <div class="form-fields">
          <div class="input-group">
            <label class="input-label" for="email">Email</label>
            <input
              id="email"
              type="email"
              class="input"
              :value="email"
              readonly
            />
          </div>
          <button class="btn btn-secondary" @click="logout">Log out</button>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { getMe } from '../api/users.js';
import { changePassword } from '../api/auth.js';

const router = useRouter();

const email = ref('');
const currentPassword = ref('');
const newPassword = ref('');
const confirmPassword = ref('');
const loading = ref(false);
const passwordError = ref('');
const passwordSuccess = ref('');

onMounted(async () => {
  try {
    const profile = await getMe();
    email.value = profile.email;
  } catch {
    router.push('/login');
  }
});

async function submitChangePassword() {
  passwordError.value = '';
  passwordSuccess.value = '';

  if (newPassword.value !== confirmPassword.value) {
    passwordError.value = 'New passwords do not match.';
    return;
  }

  loading.value = true;
  try {
    await changePassword(currentPassword.value, newPassword.value);
    passwordSuccess.value = 'Password changed successfully.';
    currentPassword.value = '';
    newPassword.value = '';
    confirmPassword.value = '';
  } catch (err) {
    passwordError.value = err.message || 'Something went wrong. Please try again.';
  } finally {
    loading.value = false;
  }
}

function logout() {
  localStorage.removeItem('token');
  router.push('/login');
}
</script>

<style scoped>
.profile-card {
  margin-top: var(--space-6);
}

.password-success {
  color: #15803d;
  font-size: var(--font-size-sm);
  margin-top: var(--space-2);
}
</style>

