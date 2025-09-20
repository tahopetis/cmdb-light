<template>
  <transition
    enter-active-class="ease-out duration-300"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
    leave-active-class="ease-in duration-200"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
  >
    <div
      v-if="isOpen"
      class="fixed inset-0 z-50 overflow-y-auto"
      aria-labelledby="modal-title"
      role="dialog"
      aria-modal="true"
    >
      <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <transition
          enter-active-class="ease-out duration-300"
          enter-from-class="opacity-0"
          enter-to-class="opacity-100"
          leave-active-class="ease-in duration-200"
          leave-from-class="opacity-100"
          leave-to-class="opacity-0"
        >
          <div
            v-if="isOpen"
            class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"
            aria-hidden="true"
            @click="close"
          ></div>
        </transition>

        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>

        <transition
          enter-active-class="ease-out duration-300"
          enter-from-class="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
          enter-to-class="opacity-100 translate-y-0 sm:scale-100"
          leave-active-class="ease-in duration-200"
          leave-from-class="opacity-100 translate-y-0 sm:scale-100"
          leave-to-class="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
        >
          <div
            v-if="isOpen"
            class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle"
            :class="[
              size === 'sm' ? 'sm:max-w-sm' : '',
              size === 'md' ? 'sm:max-w-lg' : '',
              size === 'lg' ? 'sm:max-w-2xl' : '',
              size === 'xl' ? 'sm:max-w-4xl' : '',
              size === 'full' ? 'sm:max-w-full sm:m-4' : 'sm:max-w-lg'
            ]"
          >
            <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
              <div class="sm:flex sm:items-start">
                <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left w-full">
                  <h3
                    v-if="title"
                    class="text-lg leading-6 font-medium text-gray-900"
                    id="modal-title"
                  >
                    {{ title }}
                  </h3>
                  <div class="mt-2">
                    <component :is="content" v-if="content && typeof content === 'object'" />
                    <div v-else-if="content && typeof content === 'string'" v-html="content"></div>
                    <div v-else>
                      <slot></slot>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            
            <div
              v-if="actions.length > 0"
              class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse"
            >
              <button
                v-for="(action, index) in actions"
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
              
              <button
                v-if="!actions.some(a => a.text === 'Cancel')"
                @click="close"
                type="button"
                class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
              >
                Cancel
              </button>
            </div>
          </div>
        </transition>
      </div>
    </div>
  </transition>
</template>

<script>
import { computed } from 'vue'
import { useUIStore } from '../../stores/ui'

export default {
  name: 'Modal',
  setup() {
    const uiStore = useUIStore()
    
    const isOpen = computed(() => uiStore.modal.isOpen)
    const title = computed(() => uiStore.modal.title)
    const content = computed(() => uiStore.modal.content)
    const size = computed(() => uiStore.modal.size)
    const actions = computed(() => uiStore.modal.actions)
    
    const close = () => {
      uiStore.closeModal()
    }
    
    return {
      isOpen,
      title,
      content,
      size,
      actions,
      close
    }
  }
}
</script>