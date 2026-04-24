import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/services/api'

export const useAppointmentStore = defineStore('appointments', () => {
  const appointments = ref([])

  // CORE DATA
  const branches     = ref([])
  const departments  = ref([])
  const doctors      = ref([])
  const slots        = ref([])

  // UI STATE
  const loading      = ref(false)
  const error        = ref(null)

  /* ---------------------------------------------
     FETCH BRANCHES
  ----------------------------------------------*/
  async function fetchBranches() {
    try {
      const { data } = await api.get('/branches')
      branches.value = data || []
    } catch (err) {
      console.error("Failed to load branches:", err)
      branches.value = []
    }
  }

  /* ---------------------------------------------
     FETCH DEPARTMENTS (OPTIONAL branch_id)
  ----------------------------------------------*/
  async function fetchDepartments(branchId = null) {
    try {
      const params = branchId ? { branch_id: branchId } : {}
      const { data } = await api.get('/departments', { params })
      departments.value = data || []
    } catch (err) {
      console.error("Failed to load departments:", err)
      departments.value = []
    }
  }

  /* ---------------------------------------------
     FETCH DOCTORS (BY department_id)
  ----------------------------------------------*/
  async function fetchDoctors(departmentId = null) {
    try {
      const params = departmentId ? { department_id: departmentId } : {}
      const { data } = await api.get('/doctors', { params })
      doctors.value = data || []
    } catch (err) {
      console.error("Failed to load doctors:", err)
      doctors.value = []
    }
  }

  /* ---------------------------------------------
     FETCH SLOTS (BY doctor + date)
  ----------------------------------------------*/
  async function fetchSlots(doctorId, date) {
    try {
      const { data } = await api.get('/slots', {
        params: { doctor_id: doctorId, date }
      })
      slots.value = data || []
    } catch (err) {
      console.error("Failed to load slots:", err)
      slots.value = []
    }
  }

  /* ---------------------------------------------
     FETCH MY APPOINTMENTS
  ----------------------------------------------*/
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

  /* ---------------------------------------------
     BOOK AN APPOINTMENT
  ----------------------------------------------*/
  async function bookAppointment(payload) {
    const { data } = await api.post('/appointments', payload)
    return data
  }

  /* ---------------------------------------------
     CANCEL APPOINTMENT
  ----------------------------------------------*/
  async function cancelAppointment(id) {
    await api.patch(`/appointments/${id}`, { status: 'cancelled' })
    await fetchMyAppointments()
  }

  /* ---------------------------------------------
     RETURN STATE & METHODS
  ----------------------------------------------*/
  return {
    // state
    appointments,
    branches,
    departments,
    doctors,
    slots,
    loading,
    error,

    // actions
    fetchBranches,
    fetchDepartments,
    fetchDoctors,
    fetchSlots,
    fetchMyAppointments,
    bookAppointment,
    cancelAppointment
  }
})