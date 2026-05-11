<template>
  <main class="login-page">
    <div class="container">
      <div class="card login-card">
        <h1 class="login-heading">Create account</h1>
        <form @submit.prevent="handleSubmit" class="form-fields" novalidate>
          <div class="input-group">
            <label class="input-label" for="firstName">First name</label>
            <input
              id="firstName"
              name="firstName"
              v-model="firstName"
              class="input"
              type="text"
              autocomplete="given-name"
              placeholder="Alice"
              required
            />
          </div>
          <div class="input-group">
            <label class="input-label" for="email">Email</label>
            <input
              id="email"
              v-model="email"
              class="input"
              type="email"
              autocomplete="email"
              placeholder="you@example.com"
              required
            />
            <p
              class="input-hint validation-rule"
              :class="{ 'validation-rule--met': emailLooksValid }"
              data-testid="email-format-indicator"
            >
              <span class="validation-rule__badge">{{ emailLooksValid ? 'Ready' : 'Waiting' }}</span>
              Valid email format
            </p>
          </div>
          <div class="input-group">
            <label class="input-label" for="password">Password</label>
            <input
              id="password"
              name="password"
              v-model="password"
              class="input"
              type="password"
              autocomplete="new-password"
              placeholder="••••••••"
              required
            />
            <ul class="validation-list" aria-label="Password requirements">
              <li
                class="input-hint validation-rule"
                :class="{ 'validation-rule--met': passwordLongEnough }"
                data-testid="password-length-indicator"
              >
                <span class="validation-rule__badge">{{ passwordLongEnough ? 'Ready' : 'Waiting' }}</span>
                At least 8 characters
              </li>
              <li
                class="input-hint validation-rule"
                :class="{ 'validation-rule--met': passwordHasSpecialChar }"
                data-testid="password-special-indicator"
              >
                <span class="validation-rule__badge">{{ passwordHasSpecialChar ? 'Ready' : 'Waiting' }}</span>
                Includes a special character
              </li>
            </ul>
          </div>
          <div class="input-group">
            <label class="input-label" for="passwordConfirmation">Confirm password</label>
            <input
              id="passwordConfirmation"
              name="passwordConfirmation"
              v-model="passwordConfirmation"
              class="input"
              type="password"
              autocomplete="new-password"
              placeholder="••••••••"
              required
            />
            <p
              class="input-hint validation-rule"
              :class="{ 'validation-rule--met': passwordsMatch }"
              data-testid="password-match-indicator"
            >
              <span class="validation-rule__badge">{{ passwordsMatch ? 'Ready' : 'Waiting' }}</span>
              Passwords match
            </p>
          </div>
          <div v-if="error" class="alert-error" role="alert">{{ error }}</div>
          <button type="submit" class="btn btn-primary">Create account</button>
        </form>
        <p class="login-footer">
          Already have an account?
          <router-link to="/login">Sign in</router-link>
        </p>
      </div>
    </div>
  </main>
</template>

<script setup>
import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';
import { register } from '../api/auth.js';

const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

const email = ref('');
const password = ref('');
const passwordConfirmation = ref('');
const firstName = ref('');
const error = ref('');
const router = useRouter();
const emailLooksValid = computed(() => emailPattern.test(email.value));
const passwordLongEnough = computed(() => password.value.length >= 8);
const passwordHasSpecialChar = computed(() => /[^A-Za-z0-9]/.test(password.value));
const passwordsMatch = computed(
  () => passwordConfirmation.value.length > 0 && password.value === passwordConfirmation.value
);

async function handleSubmit() {
  error.value = '';

  if (!emailLooksValid.value) {
    error.value = 'Enter a valid email address.';
    return;
  }

  if (!passwordConfirmation.value) {
    error.value = 'Please confirm your password.';
    return;
  }

  if (password.value !== passwordConfirmation.value) {
    error.value = 'Passwords must match.';
    return;
  }

  if (password.value.length < 8) {
    error.value = 'Password must be at least 8 characters.';
    return;
  }

  if (!passwordHasSpecialChar.value) {
    error.value = 'Password must include a special character.';
    return;
  }

  try {
    await register(email.value, password.value, firstName.value);
    router.push({ path: '/login', query: { registered: '1' } });
  } catch (err) {
    if (err.message === 'duplicate email') {
      error.value = 'An account with that email already exists.';
    } else if (err.message) {
      error.value = err.message;
    } else {
      error.value = 'Registration failed. Please try again.';
    }
  }
}
</script>

<style scoped>
.validation-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  margin: var(--space-2) 0 0;
  padding: 0;
  list-style: none;
}

.validation-rule {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  margin-top: var(--space-2);
}

.validation-rule--met {
  color: var(--color-primary);
}

.validation-rule__badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 4.5rem;
  padding: 0 var(--space-2);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-full);
  background: var(--color-surface);
  color: var(--color-text-muted);
  font-size: var(--font-size-xs);
}

.validation-rule--met .validation-rule__badge {
  border-color: var(--color-primary);
  color: var(--color-primary);
}

.login-footer {
  text-align: center;
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
}
</style>
