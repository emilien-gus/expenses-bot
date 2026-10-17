<template>
  <div class="app-container">
    <!-- Header with statistics & Add button -->
    <AppHeader
      :total-sum="totalSum"
      :total-count="totalCount"
      :user-count="users.length"
      @open-create="openCreateModal"
    />

    <!-- Filters -->
    <ExpenseFilters
      v-model="filters"
      :users="users"
      @change="loadExpenses"
    />

    <!-- Table of expenses -->
    <ExpenseTable
      :items="expenses"
      :loading="loading"
      @edit="openEditModal"
      @delete="openDeleteModal"
      @reset-filters="handleResetFilters"
    />

    <!-- Expense Create / Edit Modal -->
    <ExpenseForm
      v-if="formModalOpen"
      :expense="activeExpense"
      @close="closeFormModal"
      @save="handleSaveExpense"
    />

    <!-- Delete Confirmation Modal -->
    <ConfirmDialog
      v-if="deleteModalOpen"
      :expense="expenseToDelete"
      @cancel="closeDeleteModal"
      @confirm="handleConfirmDelete"
    />

    <!-- Toast Notifications -->
    <div class="toast-container">
      <div
        v-for="toast in toasts"
        :key="toast.id"
        :class="['toast', `toast-${toast.type}`]"
      >
        <svg v-if="toast.type === 'success'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#10b981" stroke-width="2.5">
          <polyline points="20 6 9 17 4 12" />
        </svg>
        <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#f43f5e" stroke-width="2.5">
          <circle cx="12" cy="12" r="10" />
          <line x1="12" y1="8" x2="12" y2="12" />
          <line x1="12" y1="16" x2="12.01" y2="16" />
        </svg>
        <span>{{ toast.message }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import AppHeader from './components/AppHeader.vue'
import ExpenseFilters from './components/ExpenseFilters.vue'
import ExpenseTable from './components/ExpenseTable.vue'
import ExpenseForm from './components/ExpenseForm.vue'
import ConfirmDialog from './components/ConfirmDialog.vue'
import {
  fetchExpenses,
  fetchUsers,
  createExpense,
  updateExpense,
  deleteExpense
} from './api/expenses'

// State
const expenses = ref([])
const totalCount = ref(0)
const users = ref([])
const loading = ref(false)

const filters = reactive({
  userId: '',
  category: '',
  dateFrom: '',
  dateTo: ''
})

const formModalOpen = ref(false)
const activeExpense = ref(null)

const deleteModalOpen = ref(false)
const expenseToDelete = ref(null)

const toasts = ref([])

// Computed
const totalSum = computed(() => {
  return expenses.value.reduce((acc, curr) => acc + (Number(curr.amount) || 0), 0)
})

// Data Loading
async function loadExpenses() {
  loading.value = true
  try {
    const res = await fetchExpenses(filters)
    expenses.value = res.data || []
    totalCount.value = res.total || 0
  } catch (err) {
    showToast(err.message || 'Ошибка при загрузке трат', 'error')
  } finally {
    loading.value = false
  }
}

async function loadUsers() {
  try {
    const res = await fetchUsers()
    users.value = res.data || []
  } catch (err) {
    console.error('Failed to load users:', err)
  }
}

// Modal management
function openCreateModal() {
  activeExpense.value = null
  formModalOpen.value = true
}

function openEditModal(expense) {
  activeExpense.value = { ...expense }
  formModalOpen.value = true
}

function closeFormModal() {
  formModalOpen.value = false
  activeExpense.value = null
}

function openDeleteModal(expense) {
  expenseToDelete.value = expense
  deleteModalOpen.value = true
}

function closeDeleteModal() {
  deleteModalOpen.value = false
  expenseToDelete.value = null
}

function handleResetFilters() {
  filters.userId = ''
  filters.category = ''
  filters.dateFrom = ''
  filters.dateTo = ''
  loadExpenses()
}

// CRUD operations
async function handleSaveExpense(formData) {
  try {
    if (activeExpense.value && activeExpense.value.id) {
      await updateExpense(activeExpense.value.id, formData)
      showToast('Расход успешно обновлен', 'success')
    } else {
      await createExpense(formData)
      showToast('Расход успешно добавлен', 'success')
    }
    closeFormModal()
    await Promise.all([loadExpenses(), loadUsers()])
  } catch (err) {
    showToast(err.message || 'Ошибка при сохранении', 'error')
  }
}

async function handleConfirmDelete() {
  if (!expenseToDelete.value) return
  try {
    await deleteExpense(expenseToDelete.value.id)
    showToast('Расход успешно удален', 'success')
    closeDeleteModal()
    await Promise.all([loadExpenses(), loadUsers()])
  } catch (err) {
    showToast(err.message || 'Ошибка при удалении', 'error')
  }
}

// Toasts
function showToast(message, type = 'success') {
  const id = Date.now() + Math.random()
  toasts.value.push({ id, message, type })
  setTimeout(() => {
    toasts.value = toasts.value.filter(t => t.id !== id)
  }, 3500)
}

onMounted(() => {
  loadExpenses()
  loadUsers()
})
</script>
