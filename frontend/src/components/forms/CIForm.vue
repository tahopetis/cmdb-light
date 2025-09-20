<template>
  <form @submit.prevent="handleSubmit">
    <div class="space-y-6">
      <!-- Name -->
      <div>
        <label for="name" class="block text-sm font-medium text-gray-700">
          Name <span class="text-red-500">*</span>
        </label>
        <div class="mt-1">
          <input
            type="text"
            id="name"
            v-model="formData.name"
            required
            class="form-input"
            :class="{ 'border-red-500': errors.name }"
          />
          <p v-if="errors.name" class="mt-1 text-sm text-red-600">
            {{ errors.name }}
          </p>
        </div>
      </div>
      
      <!-- Type -->
      <div>
        <label for="type" class="block text-sm font-medium text-gray-700">
          Type <span class="text-red-500">*</span>
        </label>
        <div class="mt-1">
          <select
            id="type"
            v-model="formData.type"
            required
            class="form-input"
            :class="{ 'border-red-500': errors.type }"
          >
            <option value="">Select a type</option>
            <option
              v-for="type in ciStore.ciTypes"
              :key="type"
              :value="type"
            >
              {{ type }}
            </option>
          </select>
          <p v-if="errors.type" class="mt-1 text-sm text-red-600">
            {{ errors.type }}
          </p>
        </div>
      </div>
      
      <!-- Status -->
      <div>
        <label for="status" class="block text-sm font-medium text-gray-700">
          Status <span class="text-red-500">*</span>
        </label>
        <div class="mt-1">
          <select
            id="status"
            v-model="formData.status"
            required
            class="form-input"
            :class="{ 'border-red-500': errors.status }"
          >
            <option value="">Select a status</option>
            <option value="Active">Active</option>
            <option value="Inactive">Inactive</option>
            <option value="Pending">Pending</option>
            <option value="Maintenance">Maintenance</option>
            <option value="Retired">Retired</option>
          </select>
          <p v-if="errors.status" class="mt-1 text-sm text-red-600">
            {{ errors.status }}
          </p>
        </div>
      </div>
      
      <!-- Description -->
      <div>
        <label for="description" class="block text-sm font-medium text-gray-700">
          Description
        </label>
        <div class="mt-1">
          <textarea
            id="description"
            v-model="formData.description"
            rows="3"
            class="form-input"
          ></textarea>
        </div>
      </div>
      
      <!-- Tags -->
      <div>
        <label for="tags" class="block text-sm font-medium text-gray-700">
          Tags
        </label>
        <div class="mt-1">
          <input
            type="text"
            id="tags"
            v-model="tagsInput"
            @keydown.enter.prevent="addTag"
            @keydown.comma.prevent="addTag"
            class="form-input"
            placeholder="Type a tag and press Enter or Comma"
          />
          <p class="mt-1 text-sm text-gray-500">
            Press Enter or Comma to add a tag
          </p>
        </div>
        
        <div v-if="formData.tags.length > 0" class="mt-2">
          <div class="flex flex-wrap gap-2">
            <span
              v-for="(tag, index) in formData.tags"
              :key="index"
              class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-primary-100 text-primary-800"
            >
              {{ tag }}
              <button
                type="button"
                @click="removeTag(index)"
                class="flex-shrink-0 ml-1.5 h-4 w-4 rounded-full inline-flex items-center justify-center text-primary-400 hover:bg-primary-200 hover:text-primary-500 focus:outline-none"
              >
                <span class="sr-only">Remove tag</span>
                <svg class="h-2 w-2" stroke="currentColor" fill="none" viewBox="0 0 8 8">
                  <path stroke-linecap="round" stroke-width="1.5" d="M1 1l6 6m0-6L1 7" />
                </svg>
              </button>
            </span>
          </div>
        </div>
      </div>
      
      <!-- Custom Attributes -->
      <div>
        <div class="flex items-center justify-between">
          <label class="block text-sm font-medium text-gray-700">
            Custom Attributes
          </label>
          <button
            type="button"
            @click="addCustomAttribute"
            class="inline-flex items-center px-2.5 py-1.5 border border-transparent text-xs font-medium rounded text-primary-700 bg-primary-100 hover:bg-primary-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            Add Attribute
          </button>
        </div>
        
        <div v-if="formData.custom_attributes.length > 0" class="mt-2 space-y-3">
          <div
            v-for="(attr, index) in formData.custom_attributes"
            :key="index"
            class="grid grid-cols-1 gap-3 sm:grid-cols-12"
          >
            <div class="sm:col-span-5">
              <input
                type="text"
                v-model="attr.key"
                placeholder="Attribute name"
                class="form-input"
              />
            </div>
            <div class="sm:col-span-6">
              <input
                type="text"
                v-model="attr.value"
                placeholder="Attribute value"
                class="form-input"
              />
            </div>
            <div class="sm:col-span-1 flex items-center justify-center">
              <button
                type="button"
                @click="removeCustomAttribute(index)"
                class="text-red-600 hover:text-red-900"
              >
                <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Form actions -->
      <div class="flex justify-end space-x-3">
        <button
          type="button"
          @click="$emit('cancel')"
          class="btn btn-secondary"
        >
          Cancel
        </button>
        <button
          type="submit"
          :disabled="isSubmitting"
          class="btn btn-primary"
          :class="{ 'opacity-75 cursor-not-allowed': isSubmitting }"
        >
          <span v-if="isSubmitting" class="flex items-center">
            <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Saving...
          </span>
          <span v-else>
            {{ isEdit ? 'Update' : 'Create' }} Configuration Item
          </span>
        </button>
      </div>
    </div>
  </form>
</template>

<script>
import { ref, reactive, computed, watch } from 'vue'
import { useCIStore } from '../../stores/ci'

export default {
  name: 'CIForm',
  props: {
    ci: {
      type: Object,
      default: null
    },
    isEdit: {
      type: Boolean,
      default: false
    }
  },
  emits: ['submit', 'cancel'],
  setup(props, { emit }) {
    const ciStore = useCIStore()
    
    const isSubmitting = ref(false)
    const tagsInput = ref('')
    const errors = reactive({})
    
    const formData = reactive({
      name: '',
      type: '',
      status: 'Active',
      description: '',
      tags: [],
      attributes: {} // For JSONB attributes
    })
    
    // Initialize form data if editing
    if (props.isEdit && props.ci) {
      formData.name = props.ci.name || ''
      formData.type = props.ci.type || ''
      formData.status = props.ci.status || 'Active'
      formData.description = props.ci.description || ''
      formData.tags = props.ci.tags || []
      formData.attributes = props.ci.attributes || {}
      
      // Convert attributes to custom_attributes array for form display
      if (props.ci.attributes) {
        formData.custom_attributes = Object.keys(props.ci.attributes).map(key => ({
          key,
          value: props.ci.attributes[key]
        }))
      } else {
        formData.custom_attributes = []
      }
    } else {
      formData.custom_attributes = []
    }
    
    const validateForm = () => {
      // Reset errors
      Object.keys(errors).forEach(key => {
        errors[key] = ''
      })
      
      let isValid = true
      
      // Validate name
      if (!formData.name.trim()) {
        errors.name = 'Name is required'
        isValid = false
      }
      
      // Validate type
      if (!formData.type) {
        errors.type = 'Type is required'
        isValid = false
      }
      
      // Validate status
      if (!formData.status) {
        errors.status = 'Status is required'
        isValid = false
      }
      
      return isValid
    }
    
    const handleSubmit = async () => {
      if (!validateForm()) return
      
      isSubmitting.value = true
      
      try {
        // Prepare data for submission
        const submitData = {
          name: formData.name.trim(),
          type: formData.type,
          status: formData.status,
          description: formData.description.trim(),
          tags: formData.tags,
          attributes: {}
        }
        
        // Convert custom_attributes array to attributes object
        formData.custom_attributes.forEach(attr => {
          if (attr.key.trim() && attr.value.trim()) {
            submitData.attributes[attr.key.trim()] = attr.value.trim()
          }
        })
        
        emit('submit', submitData)
      } catch (error) {
        console.error('Error submitting form:', error)
      } finally {
        isSubmitting.value = false
      }
    }
    
    const addTag = () => {
      const tag = tagsInput.value.trim()
      if (tag && !formData.tags.includes(tag)) {
        formData.tags.push(tag)
        tagsInput.value = ''
      }
    }
    
    const removeTag = (index) => {
      formData.tags.splice(index, 1)
    }
    
    const addCustomAttribute = () => {
      formData.custom_attributes.push({ key: '', value: '' })
    }
    
    const removeCustomAttribute = (index) => {
      formData.custom_attributes.splice(index, 1)
    }
    
    return {
      ciStore,
      isSubmitting,
      tagsInput,
      errors,
      formData,
      isEdit: props.isEdit,
      handleSubmit,
      addTag,
      removeTag,
      addCustomAttribute,
      removeCustomAttribute
    }
  }
}
</script>