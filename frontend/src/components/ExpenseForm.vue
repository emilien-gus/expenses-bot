<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content">
      <div class="modal-header">
        <h2 class="modal-title">
          {{ isEdit ? 'Редактировать расход' : 'Новый расход' }}
        </h2>
        <button class="btn-icon" title="Закрыть" @click="$emit('close')">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18" />
            <line x1="6" y1="6" x2="18" y2="18" />
          </svg>
        </button>
      </div>

      <form id="expense-form" @submit.prevent="handleSubmit">
        <div class="modal-body">
          <!-- User ID -->
          <div class="form-group">
            <label for="form-user-id" class="form-label">Telegram ID члена семьи *</label>
            <input
              id="form-user-id"
              v-model.number="form.user_id"
              type="number"
              class="form-control"
              placeholder="Например: 123456789"
              required
            />
          </div>

          <!-- Amount -->
          <div class="form-group">
            <label for="form-amount" class="form-label">Сумма (₽) *</label>
            <input
              id="form-amount"
              v-model.number="form.amount"
              type="number"
              step="0.01"
              min="0.01"
              class="form-control"
              placeholder="0.00"
              required
            />
          </div>

          <!-- Category -->
          <div class="form-group">
            <label for="form-category" class="form-label">Категория *</label>
            <select
              id="form-category"
              v-model="form.category"
              class="form-control"
              required
            >
              <option value="" disabled>Выберите категорию</option>
              <option v-for="cat in categories" :key="cat" :value="cat">
                {{ cat }}
              </option>
            </select>
          </div>

          <!-- Date -->
          <div class="form-group">
            <label for="form-spent-at" class="form-label">Дата расхода *</label>
            <input
              id="form-spent-at"
              v-model="form.spent_at"
              type="date"
              class="form-control"
              required
            />
          </div>

          <!-- Comment -->
          <div class="form-group">
            <label for="form-comment" class="form-label">Комментарий / Описание</label>
            <input
              id="form-comment"
              v-model="form.comment"
              type="text"
              class="form-control"
              placeholder="Например: кофе, такси до работы"
            />
          </div>

          <div v-if="error" class="form-error">
            {{ error }}
          </div>
        </div>

        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="$emit('close')">
            Отмена
          </button>
          <button
            id="btn-submit-expense"
            type="submit"
            class="btn btn-primary"
            :disabled="submitting"
          >
            <span v-if="submitting">Сохранение...</span>
            <span v-else>{{ isEdit ? 'Сохранить изменения' : 'Создать расход' }}</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { CATEGORIES } from '../api/expenses'

const props = defineProps({
  expense: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close', 'save'])
const categories = CATEGORIES
const isEdit = !!props.expense
const submitting = ref(false)
const error = ref('')

const form = reactive({
  user_id: props.expense ? props.expense.user_id : '',
  amount: props.expense ? props.expense.amount : '',
  category: props.expense ? props.expense.category : 'еда',
  spent_at: props.expense ? formatDateForInput(props.expense.spent_at) : getTodayDate(),
  comment: props.expense ? props.expense.comment : ''
})

function getTodayDate() {
  const d = new Date()
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

function formatDateForInput(dateStr) {
  if (!dateStr) return getTodayDate()
  return dateStr.slice(0, 10)
}

function handleSubmit() {
  error.value = ''
  if (!form.user_id || form.user_id <= 0) {
    error.value = 'Пожалуйста, укажите корректный User ID'
    return
  }
  if (!form.amount || form.amount <= 0) {
    error.value = 'Сумма должна быть больше 0'
    return
  }
  if (!form.category) {
    error.value = 'Выберите категорию расхода'
    return
  }
  if (!form.spent_at) {
    error.value = 'Укажите дату расхода'
    return
  }

  submitting.value = true
  emit('save', {
    user_id: Number(form.user_id),
    amount: Number(form.amount),
    category: form.category,
    spent_at: form.spent_at,
    comment: form.comment || ''
  })
}
</script>

<style scoped>
.form-error {
  padding: 0.75rem 1rem;
  border-radius: var(--radius-sm);
  background: rgba(244, 63, 94, 0.15);
  border: 1px solid rgba(244, 63, 94, 0.3);
  color: #fb7185;
  font-size: 0.8125rem;
  margin-top: 1rem;
}
</style>
