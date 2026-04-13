<template>
  <NavBar v-if="showNavBar" />
  <RouterView />
</template>

<script setup>
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import NavBar from './components/NavBar.vue';

const route = useRoute();
const showNavBar = computed(() => route.path !== '/login');
</script>

