<template>
  <div class="page">
    <main class="container section">
      <h1>Admin</h1>

      <div v-if="loading" data-testid="loading" class="loading-indicator">
        Loading…
      </div>

      <table v-else class="admin-table">
        <thead>
          <tr>
            <th>Email</th>
            <th>Name</th>
            <th>Role</th>
            <th>Status</th>
            <th>Created</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id">
            <td>{{ user.email }}</td>
            <td>{{ user.firstName }}</td>
            <td>{{ user.role }}</td>
            <td>{{ user.isActive ? 'Active' : 'Inactive' }}</td>
            <td>{{ formatDate(user.createdAt) }}</td>
            <td>
              <button
                v-if="user.isActive"
                :data-testid="`deactivate-${user.id}`"
                class="btn btn-secondary"
                @click="handleDeactivate(user)"
              >
                Deactivate
              </button>
              <button
                v-else
                :data-testid="`activate-${user.id}`"
                class="btn btn-primary"
                @click="handleActivate(user)"
              >
                Activate
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { getAdminUsers, activateUser, deactivateUser } from '../api/admin.js';

const users = ref([]);
const loading = ref(true);

onMounted(async () => {
  users.value = await getAdminUsers();
  loading.value = false;
});

async function handleDeactivate(user) {
  await deactivateUser(user.id);
  user.isActive = false;
}

async function handleActivate(user) {
  await activateUser(user.id);
  user.isActive = true;
}

function formatDate(iso) {
  return new Date(iso).toLocaleDateString();
}
</script>

<style scoped>
.admin-table {
  width: 100%;
  border-collapse: collapse;
}

.admin-table th,
.admin-table td {
  padding: var(--space-3) var(--space-4);
  text-align: left;
  border-bottom: 1px solid var(--color-border);
}

.loading-indicator {
  padding: var(--space-5);
  text-align: center;
  color: var(--color-text-muted);
}
</style>
