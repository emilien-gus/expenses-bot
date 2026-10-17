const API_BASE = '/api'

export const CATEGORIES = [
  'еда',
  'транспорт',
  'жильё',
  'развлечения',
  'здоровье',
  'покупки',
  'другое'
]

export async function fetchExpenses(filters = {}) {
  const params = new URLSearchParams()
  if (filters.userId) params.append('user_id', filters.userId)
  if (filters.category) params.append('category', filters.category)
  if (filters.dateFrom) params.append('date_from', filters.dateFrom)
  if (filters.dateTo) params.append('date_to', filters.dateTo)

  const url = `${API_BASE}/expenses?${params.toString()}`
  const res = await fetch(url)
  if (!res.ok) {
    const errorData = await res.json().catch(() => ({}))
    throw new Error(errorData.error || `Ошибка загрузки (${res.status})`)
  }
  return res.json()
}

export async function fetchUsers() {
  const res = await fetch(`${API_BASE}/users`)
  if (!res.ok) {
    throw new Error(`Ошибка загрузки списка пользователей (${res.status})`)
  }
  return res.json()
}

export async function createExpense(expense) {
  const res = await fetch(`${API_BASE}/expenses`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(expense)
  })
  if (!res.ok) {
    const errorData = await res.json().catch(() => ({}))
    throw new Error(errorData.error || `Ошибка создания траты (${res.status})`)
  }
  return res.json()
}

export async function updateExpense(id, expense) {
  const res = await fetch(`${API_BASE}/expenses/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(expense)
  })
  if (!res.ok) {
    const errorData = await res.json().catch(() => ({}))
    throw new Error(errorData.error || `Ошибка обновления траты (${res.status})`)
  }
  return res.json()
}

export async function deleteExpense(id) {
  const res = await fetch(`${API_BASE}/expenses/${id}`, {
    method: 'DELETE'
  })
  if (!res.ok) {
    const errorData = await res.json().catch(() => ({}))
    throw new Error(errorData.error || `Ошибка удаления траты (${res.status})`)
  }
  return res.json()
}
