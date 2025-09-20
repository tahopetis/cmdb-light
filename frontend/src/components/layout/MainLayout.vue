<template>
  <div class="h-screen flex overflow-hidden bg-gray-100">
    <!-- Sidebar -->
    <Sidebar :sidebar-open="uiStore.sidebarOpen" />
    
    <!-- Main content -->
    <div class="flex flex-col w-0 flex-1 overflow-hidden">
      <!-- Header -->
      <Header />
      
      <!-- Main panel -->
      <main class="flex-1 relative overflow-y-auto focus:outline-none">
        <!-- Loading overlay -->
        <div
          v-if="uiStore.isLoading"
          class="absolute inset-0 bg-gray-500 bg-opacity-25 z-10 flex items-center justify-center"
        >
          <div class="bg-white p-4 rounded-lg shadow-lg">
            <div class="flex items-center">
              <svg class="animate-spin h-5 w-5 text-primary-600 mr-3" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span class="text-gray-700">Loading...</span>
            </div>
          </div>
        </div>
        
        <!-- Page content -->
        <div class="py-6">
          <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <slot />
          </div>
        </div>
      </main>
    </div>
    
    <!-- Notification -->
    <div
      v-if="uiStore.notification.show"
      class="fixed bottom-4 right-4 z-50 max-w-md w-full bg-white shadow-lg rounded-lg pointer-events-auto ring-1 ring-black ring-opacity-5 overflow-hidden"
    >
      <div class="p-4">
        <div class="flex items-start">
          <div class="flex-shrink-0">
            <svg
              v-if="uiStore.notification.type === 'success'"
              class="h-6 w-6 text-green-400"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              aria-hidden="true"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <svg
              v-else-if="uiStore.notification.type === 'error'"
              class="h-6 w-6 text-red-400"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              aria-hidden="true"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
              />
            </svg>
            <svg
              v-else-if="uiStore.notification.type === 'warning'"
              class="h-6 w-6 text-yellow-400"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              aria-hidden="true"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
              />
            </svg>
            <svg
              v-else
              class="h-6 w-6 text-blue-400"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              aria-hidden="true"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          </div>
          <div class="ml-3 w-0 flex-1 pt-0.5">
            <p
              v-if="uiStore.notification.title"
              class="text-sm font-medium text-gray-900"
            >
              {{ uiStore.notification.title }}
            </p>
            <p class="text-sm text-gray-500">
              {{ uiStore.notification.message }}
            </p>
          </div>
          <div class="ml-4 flex-shrink-0 flex">
            <button
              @click="uiStore.hideNotification"
              class="bg-white rounded-md inline-flex text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
            >
              <span class="sr-only">Close</span>
              <svg
                class="h-5 w-5"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 20 20"
                fill="currentColor"
                aria-hidden="true"
              >
                <path
                  fill-rule="evenodd"
                  d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                  clip-rule="evenodd"
                />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Modal -->
    <div
      v-if="uiStore.modal.isOpen"
      class="fixed z-50 inset-0 overflow-y-auto"
      aria-labelledby="modal-title"
      role="dialog"
      aria-modal="true"
    >
      <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div
          class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"
          aria-hidden="true"
          @click="uiStore.closeModal"
        ></div>
        
        <span
          class="hidden sm:inline-block sm:align-middle sm:h-screen"
          aria-hidden="true"
          >&#8203;</span
        >
        
        <div
          class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle"
          :class="[
            uiStore.modal.size === 'sm' ? 'sm:max-w-lg' : '',
            uiStore.modal.size === 'md' ? 'sm:max-w-xl' : '',
            uiStore.modal.size === 'lg' ? 'sm:max-w-3xl' : '',
            uiStore.modal.size === 'xl' ? 'sm:max-w-5xl' : '',
            'sm:w-full'
          ]"
        >
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <div class="sm:flex sm:items-start">
              <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left w-full">
                <h3
                  class="text-lg leading-6 font-medium text-gray-900"
                  id="modal-title"
                >
                  {{ uiStore.modal.title }}
                </h3>
                <div class="mt-4">
                  <component
                    :is="uiStore.modal.content"
                    v-if="typeof uiStore.modal.content === 'object'"
                  />
                  <p v-else class="text-sm text-gray-500">
                    {{ uiStore.modal.content }}
                  </p>
                </div>
              </div>
            </div>
          </div>
          <div
            v-if="uiStore.modal.actions.length > 0"
            class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse"
          >
            <button
              v-for="(action, index) in uiStore.modal.actions"
              :key="index"
              @click="action.handler"
              type="button"
              class="w-full inline-flex justify-center rounded-md border shadow-sm px-4 py-2 text-base font-medium focus:outline-none focus:ring-2 focus:ring-offset-2 sm:ml-3 sm:w-auto sm:text-sm"
              :class="[
                action.primary
                  ? 'border-transparent bg-primary-600 text-white hover:bg-primary-700 focus:ring-primary-500'
                  : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-50 focus:ring-primary-500'
              ]"
            >
              {{ action.text }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'
import Header from './Header.vue'
import Sidebar from './Sidebar.vue'
import { useUIStore } from '../../stores/ui'

export default {
  name: 'MainLayout',
  components: {
    Header,
    Sidebar
  },
  setup() {
    const uiStore = useUIStore()
    
    return {
      uiStore
    }
  }
}
</script>