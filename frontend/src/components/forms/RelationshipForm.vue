<template>
  <div class="bg-white px-4 py-5 sm:p-6">
    <h3 class="text-lg leading-6 font-medium text-gray-900 mb-4">
      {{ isEdit ? 'Edit Relationship' : 'Create New Relationship' }}
    </h3>
    
    <form @submit.prevent="submitForm">
      <div class="grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-6">
        <!-- Source CI -->
        <div class="sm:col-span-3">
          <label for="sourceCI" class="block text-sm font-medium text-gray-700">
            Source Configuration Item
          </label>
          <div class="mt-1">
            <select
              id="sourceCI"
              v-model="formData.source_id"
              required
              class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 rounded-md"
              :disabled="isEdit || !!preselectedSource"
            >
              <option value="">Select a CI</option>
              <option
                v-for="ci in availableCIs"
                :key="ci.id"
                :value="ci.id"
              >
                {{ ci.name }} ({{ ci.type }})
              </option>
            </select>
          </div>
        </div>

        <!-- Target CI -->
        <div class="sm:col-span-3">
          <label for="targetCI" class="block text-sm font-medium text-gray-700">
            Target Configuration Item
          </label>
          <div class="mt-1">
            <select
              id="targetCI"
              v-model="formData.target_id"
              required
              class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 rounded-md"
              :disabled="isEdit"
            >
              <option value="">Select a CI</option>
              <option
                v-for="ci in availableCIs"
                :key="ci.id"
                :value="ci.id"
                :disabled="ci.id === formData.source_id"
              >
                {{ ci.name }} ({{ ci.type }})
              </option>
            </select>
          </div>
        </div>

        <!-- Relationship Type -->
        <div class="sm:col-span-6">
          <label for="relationshipType" class="block text-sm font-medium text-gray-700">
            Relationship Type
          </label>
          <div class="mt-1">
            <select
              id="relationshipType"
              v-model="formData.type"
              required
              class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 rounded-md"
            >
              <option value="">Select a relationship type</option>
              <option
                v-for="type in relationshipTypes"
                :key="type"
                :value="type"
              >
                {{ type }}
              </option>
            </select>
          </div>
        </div>

        <!-- Description -->
        <div class="sm:col-span-6">
          <label for="description" class="block text-sm font-medium text-gray-700">
            Description (Optional)
          </label>
          <div class="mt-1">
            <textarea
              id="description"
              v-model="formData.description"
              rows="3"
              class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border border-gray-300 rounded-md"
            />
          </div>
        </div>
      </div>

      <div class="mt-5 sm:mt-6 sm:grid sm:grid-cols-2 sm:gap-3 sm:grid-flow-row-dense">
        <button
          type="submit"
          class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-primary-600 text-base font-medium text-white hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 sm:col-start-2 sm:text-sm"
          :disabled="isSubmitting"
        >
          <span v-if="isSubmitting">
            <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white inline-block" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Saving...
          </span>
          <span v-else>
            {{ isEdit ? 'Update Relationship' : 'Create Relationship' }}
          </span>
        </button>
        
        <button
          type="button"
          class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 sm:mt-0 sm:col-start-1 sm:text-sm"
          @click="$emit('cancel')"
        >
          Cancel
        </button>
      </div>
    </form>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useCIStore } from '../../stores/ci'

export default {
  name: 'RelationshipForm',
  props: {
    isEdit: {
      type: Boolean,
      default: false
    },
    relationship: {
      type: Object,
      default: null
    },
    preselectedSource: {
      type: String,
      default: null
    }
  },
  emits: ['submit', 'cancel'],
  setup(props, { emit }) {
    const ciStore = useCIStore()
    
    const isSubmitting = ref(false)
    
    const relationshipTypes = ref([
      'Depends On',
      'Connected To',
      'Part Of',
      'Hosts',
      'Uses',
      'Managed By',
      'Related To',
      'Implements',
      'Required By',
      'Includes'
    ])
    
    const formData = reactive({
      source_id: '',
      target_id: '',
      type: '',
      description: ''
    })
    
    const availableCIs = computed(() => {
      // Filter out the current CI if we're in edit mode to avoid self-references
      return ciStore.cis
    })
    
    const submitForm = async () => {
      // Validate form
      if (!formData.source_id || !formData.target_id || !formData.type) {
        return
      }
      
      // Prevent self-referencing
      if (formData.source_id === formData.target_id) {
        return
      }
      
      try {
        isSubmitting.value = true
        
        // Prepare data for API
        const submitData = {
          source_id: formData.source_id,
          target_id: formData.target_id,
          type: formData.type
        }
        
        let result
        if (props.isEdit && props.relationship) {
          // Update existing relationship
          result = await ciStore.updateRelationship(props.relationship.id, submitData)
        } else {
          // Create new relationship
          result = await ciStore.createRelationship(submitData)
        }
        
        emit('submit', result)
      } catch (error) {
        console.error('Error submitting relationship form:', error)
        // Error handling will be done by the parent component
      } finally {
        isSubmitting.value = false
      }
    }
    
    // Initialize form data
    const initializeForm = () => {
      if (props.isEdit && props.relationship) {
        formData.source_id = props.relationship.source_id
        formData.target_id = props.relationship.target_id
        formData.type = props.relationship.type
        formData.description = props.relationship.description || ''
      } else if (props.preselectedSource) {
        formData.source_id = props.preselectedSource
      }
    }
    
    // Load CIs if not already loaded
    const loadCIs = async () => {
      if (ciStore.cis.length === 0) {
        try {
          await ciStore.fetchCIs()
        } catch (error) {
          console.error('Error loading CIs:', error)
        }
      }
    }
    
    onMounted(() => {
      loadCIs()
      initializeForm()
    })
    
    // Watch for changes in props
    watch(() => props.relationship, initializeForm, { deep: true })
    watch(() => props.preselectedSource, (newValue) => {
      if (newValue) {
        formData.source_id = newValue
      }
    })
    
    return {
      isSubmitting,
      formData,
      relationshipTypes,
      availableCIs,
      submitForm
    }
  }
}
</script>