import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/services/api'

export const useQueueStore = defineStore('queue', () => {
  const myQueue   = ref(null)
  const fullQueue = ref([])
  const stats     = ref(null)
  const loading   = ref(false)

  async function checkIn(appointmentId) {
    const { data } = await api.post('/queue/checkin', { appointment_id: appointmentId })
    myQueue.value = data.data
    return data
  }

  async function fetchMyQueue() {
    loading.value = true
    try {
      const { data } = await api.get('/queue/my')
      myQueue.value = data
    } catch {
      myQueue.value = null
    } finally {
      loading.value = false
    }
  }

  async function fetchFullQueue(doctorId = null) {
    const params = doctorId ? { doctor_id: doctorId } : {}
    const { data } = await api.get('/queue', { params })
    fullQueue.value = data || []
  }

  async function fetchStats() {
    const { data } = await api.get('/queue/stats')
    stats.value = data
  }

  async function updateQueueEntry(id, status) {
    await api.patch(`/queue/${id}`, { status })
    await fetchFullQueue()
  }

  return { myQueue, fullQueue, stats, loading, checkIn, fetchMyQueue, fetchFullQueue, fetchStats, updateQueueEntry }
})

