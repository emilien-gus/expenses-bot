<template>
  <div class="table-card glass-card">
    <div class="table-header">
      <div class="table-title">
        <span>История трат</span>
        <span class="count-tag">{{ items.length }}</span>
      </div>

      <div class="table-actions">
        <div class="sort-control">
          <label for="table-sort" class="sort-label">Сортировка:</label>
          <select id="table-sort" v-model="sortBy" class="sort-select" @change="sortItems">
            <option value="date-desc">Сначала новые</option>
            <option value="date-asc">Сначала старые</option>
            <option value="amount-desc">По убыванию суммы</option>
            <option value="amount-asc">По возрастанию суммы</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Loading Skeleton -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Загрузка трат...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="items.length === 0" class="empty-state">
      <div class="empty-icon">
        <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <rect x="2" y="5" width="20" height="14" rx="2" />
          <line x1="2" y1="10" x2="22" y2="10" />
        </svg>
      </div>
      <h3 class="empty-title">Траты не найдены</h3>
      <p class="empty-subtitle">Нет расходов, удовлетворяющих заданным фильтрам</p>
      <button class="btn btn-secondary btn-sm" @click="$emit('reset-filters')">
        Сбросить фильтры
      </button>
    </div>

    <!-- Expenses Table -->
    <div v-else class="table-wrapper">
      <table class="expenses-table">
        <thead>
          <tr>
            <th class="th-date">Дата</th>
            <th class="th-user">Член семьи</th>
            <th class="th-category">Категория</th>
            <th class="th-comment">Описание</th>
            <th class="th-amount text-right">Сумма</th>
            <th class="th-actions text-right">Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="expense in sortedItems" :key="expense.id" class="table-row">
            <td class="td-date">
              <span class="date-main">{{ formatDate(expense.spent_at) }}</span>
            </td>
            <td class="td-user">
              <span class="user-pill">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
                  <circle cx="12" cy="7" r="4" />
                </svg>
                ID: {{ expense.user_id }}
              </span>
            </td>
            <td class="td-category">
              <span :class="['category-badge', `badge-${expense.category}`]">
                {{ expense.category }}
              </span>
            </td>
            <td class="td-comment">
              <span class="comment-text">{{ expense.comment || '—' }}</span>
            </td>
            <td class="td-amount text-right">
              <span class="amount-value">{{ formatAmount(expense.amount) }} ₽</span>
            </td>
            <td class="td-actions text-right">
              <div class="row-actions">
                <button
                  :id="`btn-edit-${expense.id}`"
                  class="btn-icon"
                  title="Редактировать"
                  @click="$emit('edit', expense)"
                >
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                    <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
                  </svg>
                </button>
                <button
                  :id="`btn-delete-${expense.id}`"
                  class="btn-icon btn-icon-danger"
                  title="Удалить"
                  @click="$emit('delete', expense)"
                >
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="3 6 5 6 21 6" />
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
                  </svg>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  items: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  }
})

defineEmits(['edit', 'delete', 'reset-filters'])

const sortBy = ref('date-desc')

const sortedItems = computed(() => {
  const list = [...props.items]

  switch (sortBy.value) {
    case 'date-desc':
      return list.sort((a, b) => new Date(b.spent_at) - new Date(a.spent_at) || b.id - a.id)
    case 'date-asc':
      return list.sort((a, b) => new Date(a.spent_at) - new Date(b.spent_at) || a.id - b.id)
    case 'amount-desc':
      return list.sort((a, b) => b.amount - a.amount)
    case 'amount-asc':
      return list.sort((a, b) => a.amount - b.amount)
    default:
      return list
  }
})

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return new Intl.DateTimeFormat('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric'
  }).format(d)
}

function formatAmount(num) {
  return new Intl.NumberFormat('ru-RU', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(num)
}
</script>

<style scoped>
.table-card {
  padding: 1.5rem;
  overflow: hidden;
}

.table-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 1rem;
  margin-bottom: 1.25rem;
}

.table-title {
  display: flex;
  align-items: center;
  gap: 0.625rem;
  font-size: 1.125rem;
  font-weight: 700;
  color: var(--text-primary);
}

.count-tag {
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  color: var(--text-muted);
  font-size: 0.75rem;
  font-weight: 600;
  padding: 0.125rem 0.5rem;
  border-radius: var(--radius-full);
}

.sort-control {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.sort-label {
  font-size: 0.8125rem;
  color: var(--text-muted);
}

.sort-select {
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  color: var(--text-secondary);
  font-size: 0.8125rem;
  padding: 0.375rem 0.75rem;
  border-radius: var(--radius-sm);
  outline: none;
  cursor: pointer;
  transition: var(--transition);
}

.sort-select:focus {
  border-color: var(--border-focus);
}

.table-wrapper {
  overflow-x: auto;
}

.expenses-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.expenses-table th {
  padding: 0.75rem 1rem;
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  border-bottom: 1px solid var(--border-subtle);
}

.expenses-table td {
  padding: 1rem;
  border-bottom: 1px solid var(--border-subtle);
  font-size: 0.875rem;
  vertical-align: middle;
}

.table-row {
  transition: var(--transition);
}

.table-row:hover {
  background: rgba(255, 255, 255, 0.02);
}

.text-right {
  text-align: right;
}

.date-main {
  font-weight: 500;
  color: var(--text-primary);
}

.user-pill {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  padding: 0.25rem 0.5rem;
  border-radius: var(--radius-sm);
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--text-secondary);
}

.comment-text {
  color: var(--text-secondary);
}

.amount-value {
  font-weight: 700;
  color: var(--text-primary);
}

.row-actions {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
}

.btn-icon-danger:hover {
  background: rgba(244, 63, 94, 0.15);
  color: #fb7185;
  border-color: rgba(244, 63, 94, 0.3);
}

/* Loading & Empty states */
.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 1.5rem;
  text-align: center;
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid rgba(99, 102, 241, 0.2);
  border-top-color: var(--accent-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-icon {
  width: 64px;
  height: 64px;
  border-radius: var(--radius-lg);
  background: var(--bg-card);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
  margin-bottom: 1rem;
}

.empty-title {
  font-size: 1.125rem;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 0.375rem;
}

.empty-subtitle {
  font-size: 0.875rem;
  color: var(--text-muted);
  margin-bottom: 1.25rem;
}
</style>
