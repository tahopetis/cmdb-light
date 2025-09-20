<template>
  <div id="app">
    <router-view />
    <Modal />
  </div>
</template>

<script>
import Modal from './components/ui/Modal.vue'

export default {
  name: 'App',
  components: {
    Modal
  }
}
</script>

<style>
/* Global styles */
#app {
  font-family: 'Inter', system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  height: 100vh;
}

/* Custom utility classes */
.card {
  @apply bg-white shadow overflow-hidden sm:rounded-lg;
}

.form-input {
  @apply shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 rounded-md;
}

.btn {
  @apply inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500;
}

.btn-primary {
  @apply text-white bg-primary-600 hover:bg-primary-700;
}

.btn-secondary {
  @apply text-gray-700 bg-white border-gray-300 hover:bg-gray-50;
}

.line-clamp-2 {
  overflow: hidden;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}
</style>