import { defineStore } from 'pinia'

export const useUIStore = defineStore('ui', {
  state: () => ({
    isLoading: false,
    sidebarOpen: false,
    notification: {
      show: false,
      type: 'info',
      message: '',
      title: ''
    },
    modal: {
      isOpen: false,
      title: '',
      content: null,
      size: 'md',
      actions: []
    }
  }),
  
  actions: {
    setLoading(isLoading) {
      this.isLoading = isLoading
    },
    
    toggleSidebar() {
      this.sidebarOpen = !this.sidebarOpen
    },
    
    setSidebarOpen(isOpen) {
      this.sidebarOpen = isOpen
    },
    
    showNotification(type, message, title = '') {
      this.notification = {
        show: true,
        type,
        message,
        title
      }
      
      // Auto-hide notification after 5 seconds
      setTimeout(() => {
        this.hideNotification()
      }, 5000)
    },
    
    showSuccess(message, title = '') {
      this.showNotification('success', message, title)
    },
    
    showError(message, title = '') {
      this.showNotification('error', message, title)
    },
    
    showWarning(message, title = '') {
      this.showNotification('warning', message, title)
    },
    
    showInfo(message, title = '') {
      this.showNotification('info', message, title)
    },
    
    hideNotification() {
      this.notification.show = false
    },
    
    openModal(options = {}) {
      this.modal = {
        isOpen: true,
        title: options.title || '',
        content: options.content || null,
        size: options.size || 'md',
        actions: options.actions || []
      }
    },
    
    closeModal() {
      this.modal.isOpen = false
    },
    
    showConfirmDialog(title, message, onConfirm, onCancel = null) {
      this.openModal({
        title,
        content: message,
        size: 'sm',
        actions: [
          {
            text: 'Cancel',
            handler: () => {
              this.closeModal()
              if (onCancel) onCancel()
            }
          },
          {
            text: 'Confirm',
            primary: true,
            handler: () => {
              this.closeModal()
              onConfirm()
            }
          }
        ]
      })
    }
  }
})