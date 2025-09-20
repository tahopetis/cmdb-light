import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

// Layouts
import MainLayout from '../components/layout/MainLayout.vue'

// Views
import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'
import CIListView from '../views/CIListView.vue'
import CIDetailView from '../views/CIDetailView.vue'
import CreateEditCIView from '../views/CreateEditCIView.vue'
import GraphView from '../views/GraphView.vue'
import AuditLogsView from '../views/AuditLogsView.vue'
import SettingsView from '../views/SettingsView.vue'
import NotFoundView from '../views/NotFoundView.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: LoginView,
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: MainLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: DashboardView
      },
      {
        path: 'cis',
        name: 'CIList',
        component: CIListView
      },
      {
        path: 'cis/create',
        name: 'CreateCI',
        component: CreateEditCIView,
        meta: { requiresAdmin: true }
      },
      {
        path: 'cis/:id',
        name: 'CIDetail',
        component: CIDetailView
      },
      {
        path: 'cis/:id/edit',
        name: 'EditCI',
        component: CreateEditCIView,
        meta: { requiresAdmin: true }
      },
      {
        path: 'graph',
        name: 'Graph',
        component: GraphView
      },
      {
        path: 'audit-logs',
        name: 'AuditLogs',
        component: AuditLogsView,
        meta: { requiresAdmin: true }
      },
      {
        path: 'settings',
        name: 'Settings',
        component: SettingsView
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFoundView
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guard
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // Check if user is authenticated
  const isAuthenticated = authStore.checkAuth()
  
  // If route requires authentication and user is not authenticated
  if (to.meta.requiresAuth !== false && !isAuthenticated) {
    // Redirect to login with redirect parameter
    next({
      path: '/login',
      query: { redirect: to.fullPath }
    })
    return
  }
  
  // If route requires admin and user is not admin
  if (to.meta.requiresAdmin && !authStore.isAdmin) {
    // Redirect to dashboard
    next({ path: '/' })
    return
  }
  
  // If user is authenticated and trying to access login page
  if (to.path === '/login' && isAuthenticated) {
    // Redirect to dashboard
    next({ path: '/' })
    return
  }
  
  // Otherwise, proceed with navigation
  next()
})

export default router