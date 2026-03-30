<template>
  <nav class="navbar">
    <div class="navbar-brand">
      <img :src="logoSrc" alt="Triangle of Love" width="32" height="32" />
    </div>
    <span class="navbar-greeting navbar-greeting--truncate" v-if="firstName">Hello, {{ firstName }}</span>
    <div class="avatar" aria-hidden="true">{{ initials }}</div>
  </nav>
</template>

<script setup>
import { computed } from 'vue';
import logoSrc from '../assets/logo.svg';

const props = defineProps({
  firstName: {
    type: String,
    default: ''
  }
});

const initials = computed(() =>
  props.firstName ? props.firstName.charAt(0).toUpperCase() : '?'
);
</script>
