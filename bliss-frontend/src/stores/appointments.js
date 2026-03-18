import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/services/api'

export const useAppointmentStore = defineStore('appointments', () => {
  const appointments = ref([])
  const departments  = ref([])
  const doctors      = ref([])
  const slots        = ref([])
  const loading      = ref(false)
  const error        = ref(null)

  async function fetchDepartments() {
    const { data } = await api.get('/departments')
    departments.value = data
  }

  async function fetchDoctors(departmentId = null) {
    const params = departmentId ? { department_id: departmentId } : {}
    const { data } = await api.get('/doctors', { params })
    doctors.value = data
  }

  async function fetchSlots(doctorId, date) {
    const { data } = await api.get('/slots', { params: { doctor_id: doctorId, date } })
    slots.value = data
  }

  async function fetchMyAppointments() {
    loading.value = true
    error.value   = null
    try {
      const { data } = await api.get('/appointments')
      appointments.value = data || []
    } catch (e) {
      error.value = e.response?.data?.error || 'Failed to fetch appointments'
    } finally {
      loading.value = false
    }
  }

  async function bookAppointment(payload) {
    const { data } = await api.post('/appointments', payload)
    return data
  }

  async function cancelAppointment(id) {
    await api.patch(`/appointments/${id}`, { status: 'cancelled' })
    await fetchMyAppointments()
  }

  return {
    appointments, departments, doctors, slots, loading, error,
    fetchDepartments, fetchDoctors, fetchSlots,
    fetchMyAppointments, bookAppointment, cancelAppointment
  }
})

