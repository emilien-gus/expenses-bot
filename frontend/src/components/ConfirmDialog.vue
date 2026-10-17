<template>
  <div class="modal-overlay" @click.self="$emit('cancel')">
    <div class="modal-content modal-confirm">
      <div class="modal-header">
        <h2 class="modal-title">Подтверждение удаления</h2>
        <button class="btn-icon" title="Закрыть" @click="$emit('cancel')">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18" />
            <line x1="6" y1="6" x2="18" y2="18" />
          </svg>
        </button>
      </div>

      <div class="modal-body">
        <p class="confirm-message">
          Вы действительно хотите удалить эту запись о расходе?
        </p>

        <div v-if="expense" class="item-preview">
          <div class="preview-row">
            <span class="preview-label">Описание:</span>
            <span class="preview-value">{{ expense.comment || 'Без описания' }}</span>
          </div>
          <div class="preview-row">
            <span class="preview-label">Сумма:</span>
            <span class="preview-value font-bold text-rose">{{ formatAmount(expense.amount) }} ₽</span>
          </div>
          <div class="preview-row">
            <span class="preview-label">Категория:</span>
            <span class="preview-value">{{ expense.category }}</span>
          </div>
          <div class="preview-row">
            <span class="preview-label">Член семьи:</span>
            <span class="preview-value">ID: {{ expense.user_id }}</span>
          </div>
        </div>

        <p class="warning-text">Это действие нельзя будет отменить.</p>
      </div>

      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" @click="$emit('cancel')">
          Отмена
        </button>
        <button
          id="btn-confirm-delete"
          type="button"
          class="btn btn-danger"
          @click="$emit('confirm')"
        >
          Удалить
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  expense: {
    type: Object,
    default: null
  }
})

defineEmits(['confirm', 'cancel'])

function formatAmount(num) {
  if (num == null) return '0.00'
  return new Intl.NumberFormat('ru-RU', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(num)
}
</script>

<style scoped>
.modal-confirm {
  max-width: 440px;
}

.confirm-message {
  font-size: 0.9375rem;
  color: var(--text-primary);
  margin-bottom: 1rem;
}

.item-preview {
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  padding: 0.875rem 1rem;
  margin-bottom: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.preview-row {
  display: flex;
  justify-content: space-between;
  font-size: 0.8125rem;
}

.preview-label {
  color: var(--text-muted);
}

.preview-value {
  color: var(--text-primary);
  font-weight: 500;
}

.font-bold {
  font-weight: 700;
}

.text-rose {
  color: #fb7185;
}

.warning-text {
  font-size: 0.75rem;
  color: var(--text-muted);
}
</style>
