<template>
  <header class="header-container">
    <div class="header-top">
      <div class="brand">
        <div class="brand-icon">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6" />
          </svg>
        </div>
        <div>
          <h1 class="brand-title">Семейный бюджет</h1>
          <p class="brand-subtitle">Панель управления расходами главы семьи</p>
        </div>
      </div>

      <button id="btn-add-expense" class="btn btn-primary" @click="$emit('open-create')">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <line x1="12" y1="5" x2="12" y2="19" />
          <line x1="5" y1="12" x2="19" y2="12" />
        </svg>
        <span>Добавить трату</span>
      </button>
    </div>

    <!-- Stats summary cards -->
    <div class="stats-grid">
      <div class="stat-card glass-card">
        <div class="stat-label">Всего расходов</div>
        <div class="stat-value highlight-amount">{{ formattedTotalSum }} ₽</div>
        <div class="stat-meta">за выбранный период</div>
      </div>

      <div class="stat-card glass-card">
        <div class="stat-label">Количество трат</div>
        <div class="stat-value">{{ totalCount }}</div>
        <div class="stat-meta">записей в списке</div>
      </div>

      <div class="stat-card glass-card">
        <div class="stat-label">Членов семьи</div>
        <div class="stat-value">{{ userCount }}</div>
        <div class="stat-meta">активных участников</div>
      </div>
    </div>
  </header>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  totalSum: {
    type: Number,
    default: 0
  },
  totalCount: {
    type: Number,
    default: 0
  },
  userCount: {
    type: Number,
    default: 0
  }
})

defineEmits(['open-create'])

const formattedTotalSum = computed(() => {
  return new Intl.NumberFormat('ru-RU', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(props.totalSum)
})
</script>

<style scoped>
.header-container {
  margin-bottom: 2rem;
}

.header-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 1.25rem;
  margin-bottom: 2rem;
}

.brand {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.brand-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-md);
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.2) 0%, rgba(129, 140, 248, 0.2) 100%);
  border: 1px solid rgba(99, 102, 241, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--accent-primary);
  box-shadow: 0 0 20px var(--accent-glow);
}

.brand-title {
  font-size: 1.625rem;
  font-weight: 800;
  letter-spacing: -0.025em;
  color: var(--text-primary);
}

.brand-subtitle {
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 1.25rem;
}

.stat-card {
  padding: 1.25rem 1.5rem;
}

.stat-label {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 0.375rem;
}

.stat-value {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: -0.02em;
}

.highlight-amount {
  background: linear-gradient(135deg, #a5b4fc 0%, #6366f1 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.stat-meta {
  font-size: 0.75rem;
  color: var(--text-muted);
  margin-top: 0.25rem;
}
</style>
