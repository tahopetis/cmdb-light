<template>
  <div class="py-6">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center">
        <h1 class="text-2xl font-semibold text-gray-900">
          {{ isEdit ? 'Edit Configuration Item' : 'Create Configuration Item' }}
        </h1>
        <div class="flex space-x-3">
          <router-link
            :to="isEdit ? `/cis/${ciId}` : '/cis'"
            class="btn btn-secondary"
          >
            Cancel
          </router-link>
        </div>
      </div>
      
      <div class="mt-6">
        <div class="card">
          <div class="px-4 py-5 sm:p-6">
            <CIForm
              :ci="currentCI"
              :is-edit="isEdit"
              @submit="handleSubmit"
              @cancel="handleCancel"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useCIStore } from '../stores/ci'
import { useUIStore } from '../stores/ui'
import CIForm from '../components/forms/CIForm.vue'

export default {
  name: 'CreateEditCIView',
  components: {
    CIForm
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const ciStore = useCIStore()
    const uiStore = useUIStore()
    
    const isEdit = computed(() => route.name === 'EditCI')
    const ciId = computed(() => route.params.id)
    const currentCI = computed(() => ciStore.currentCI)
    
    const handleSubmit = async (formData) => {
      uiStore.setLoading(true)
      
      try {
        let result
        
        if (isEdit.value) {
          // Update existing CI
          result = await ciStore.updateCI(ciId.value, formData)
          
          if (result.success) {
            uiStore.showSuccess('Configuration item updated successfully')
            router.push(`/cis/${ciId.value}`)
          } else {
            uiStore.showError(result.error || 'Failed to update configuration item')
          }
        } else {
          // Create new CI
          result = await ciStore.createCI(formData)
          
          if (result.success) {
            uiStore.showSuccess('Configuration item created successfully')
            router.push('/cis')
          } else {
            uiStore.showError(result.error || 'Failed to create configuration item')
          }
        }
      } catch (error) {
        console.error('Error submitting CI form:', error)
        uiStore.showError('An unexpected error occurred')
      } finally {
        uiStore.setLoading(false)
      }
    }
    
    const handleCancel = () => {
      if (isEdit.value) {
        router.push(`/cis/${ciId.value}`)
      } else {
        router.push('/cis')
      }
    }
    
    onMounted(async () => {
      if (isEdit.value && ciId.value) {
        try {
          await ciStore.fetchCIById(ciId.value)
        } catch (error) {
          console.error('Error fetching CI:', error)
          uiStore.showError('Failed to fetch configuration item')
          router.push('/cis')
        }
      }
    })
    
    return {
      isEdit,
      ciId,
      currentCI,
      handleSubmit,
      handleCancel
    }
  }
}
</script>