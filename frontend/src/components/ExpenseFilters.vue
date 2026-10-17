<template>
  <div class="filters-card glass-card">
    <div class="filters-header">
      <div class="filters-title">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3" />
        </svg>
        <span>Фильтры</span>
      </div>

      <button
        v-if="hasActiveFilters"
        id="btn-reset-filters"
        class="btn btn-secondary btn-sm"
        @click="resetFilters"
      >
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="18" y1="6" x2="6" y2="18" />
          <line x1="6" y1="6" x2="18" y2="18" />
        </svg>
        <span>Сбросить</span>
      </button>
    </div>

    <div class="filters-grid">
      <!-- Family member filter -->
      <div class="filter-item">
        <label for="filter-user" class="form-label">Член семьи (User ID)</label>
        <select
          id="filter-user"
          v-model="localFilters.userId"
          class="form-control"
          @change="emitChange"
        >
          <option value="">Все члены семьи</option>
          <option v-for="user in users" :key="user" :value="user">
            Пользователь #{{ user }}
          </option>
        </select>
      </div>

      <!-- Category filter -->
      <div class="filter-item">
        <label for="filter-category" class="form-label">Категория</label>
        <select
          id="filter-category"
          v-model="localFilters.category"
          class="form-control"
          @change="emitChange"
        >
          <option value="">Все категории</option>
          <option v-for="cat in categories" :key="cat" :value="cat">
            {{ cat }}
          </option>
        </select>
      </div>

      <!-- Date from -->
      <div class="filter-item">
        <label for="filter-date-from" class="form-label">Период: с</label>
        <input
          id="filter-date-from"
          v-model="localFilters.dateFrom"
          type="date"
          class="form-control"
          @change="emitChange"
        />
      </div>

      <!-- Date to -->
      <div class="filter-item">
        <label for="filter-date-to" class="form-label">по</label>
        <input
          id="filter-date-to"
          v-model="localFilters.dateTo"
          type="date"
          class="form-control"
          @change="emitChange"
        />
      </div>
    </div>

    <!-- Quick period presets -->
    <div class="presets-row">
      <span class="presets-label">Быстрый выбор:</span>
      <button
        class="preset-chip"
        :class="{ active: isPresetActive('today') }"
        @click="applyPreset('today')"
      >
        Сегодня
      </button>
      <button
        class="preset-chip"
        :class="{ active: isPresetActive('week') }"
        @click="applyPreset('week')"
      >
        Последние 7 дней
      </button>
      <button
        class="preset-chip"
        :class="{ active: isPresetActive('month') }"
        @click="applyPreset('month')"
      >
        Этот месяц
      </button>
      <button
        class="preset-chip"
        :class="{ active: isPresetActive('all') }"
        @click="applyPreset('all')"
      >
        За всё время
      </button>
    </div>
  </div>
</template>

<script setup>
import { reactive, computed, watch } from 'vue'
import { CATEGORIES } from '../api/expenses'

const props = defineProps({
  users: {
    type: Array,
    default: () => []
  },
  modelValue: {
    type: Object,
    default: () => ({
      userId: '',
      category: '',
      dateFrom: '',
      dateTo: ''
    })
  }
})

const emit = defineEmits(['update:modelValue', 'change'])
const categories = CATEGORIES

const localFilters = reactive({
  userId: props.modelValue.userId || '',
  category: props.modelValue.category || '',
  dateFrom: props.modelValue.dateFrom || '',
  dateTo: props.modelValue.dateTo || ''
})

watch(() => props.modelValue, (newVal) => {
  localFilters.userId = newVal.userId || ''
  localFilters.category = newVal.category || ''
  localFilters.dateFrom = newVal.dateFrom || ''
  localFilters.dateTo = newVal.dateTo || ''
}, { deep: true })

const hasActiveFilters = computed(() => {
  return localFilters.userId || localFilters.category || localFilters.dateFrom || localFilters.dateTo
})

function emitChange() {
  const updated = { ...localFilters }
  emit('update:modelValue', updated)
  emit('change', updated)
}

function resetFilters() {
  localFilters.userId = ''
  localFilters.category = ''
  localFilters.dateFrom = ''
  localFilters.dateTo = ''
  emitChange()
}

function formatDate(d) {
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

function applyPreset(preset) {
  const now = new Date()

  if (preset === 'today') {
    const today = formatDate(now)
    localFilters.dateFrom = today
    localFilters.dateTo = today
  } else if (preset === 'week') {
    const weekAgo = new Date(now)
    weekAgo.setDate(weekAgo.getDate() - 7)
    localFilters.dateFrom = formatDate(weekAgo)
    localFilters.dateTo = formatDate(now)
  } else if (preset === 'month') {
    const firstDay = new Date(now.getFullYear(), now.getMonth(), 1)
    localFilters.dateFrom = formatDate(firstDay)
    localFilters.dateTo = formatDate(now)
  } else if (preset === 'all') {
    localFilters.dateFrom = ''
    localFilters.dateTo = ''
  }

  emitChange()
}

function isPresetActive(preset) {
  const now = new Date()
  const today = formatDate(now)

  if (preset === 'today') {
    return localFilters.dateFrom === today && localFilters.dateTo === today
  }
  if (preset === 'month') {
    const firstDay = formatDate(new Date(now.getFullYear(), now.getMonth(), 1))
    return localFilters.dateFrom === firstDay && localFilters.dateTo === today
  }
  if (preset === 'all') {
    return !localFilters.dateFrom && !localFilters.dateTo
  }
  return false
}
</script>

<style scoped>
.filters-card {
  padding: 1.5rem;
  margin-bottom: 2rem;
}

.filters-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1.25rem;
}

.filters-title {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 1rem;
  font-weight: 700;
  color: var(--text-primary);
}

.btn-sm {
  padding: 0.375rem 0.75rem;
  font-size: 0.75rem;
}

.filters-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
}

.presets-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: 1.25rem;
  padding-top: 1rem;
  border-top: 1px solid var(--border-subtle);
}

.presets-label {
  font-size: 0.75rem;
  color: var(--text-muted);
  margin-right: 0.25rem;
}

.preset-chip {
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  color: var(--text-secondary);
  font-size: 0.75rem;
  font-weight: 500;
  padding: 0.25rem 0.625rem;
  border-radius: var(--radius-full);
  cursor: pointer;
  transition: var(--transition);
}

.preset-chip:hover {
  background: var(--bg-card-hover);
  color: var(--text-primary);
  border-color: var(--border-hover);
}

.preset-chip.active {
  background: rgba(99, 102, 241, 0.2);
  color: #a5b4fc;
  border-color: rgba(99, 102, 241, 0.4);
}
</style>
