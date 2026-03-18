package models

import "time"

// ── USER ──────────────────────────────────────────────────────
type User struct {
	UserID    string    `json:"user_id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email"     binding:"required,email"`
	Phone    string `json:"phone"`
	Password string `json:"password"  binding:"required,min=6"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// ── DEPARTMENT ────────────────────────────────────────────────
type Department struct {
	DepartmentID string    `json:"department_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
}

// ── DOCTOR ────────────────────────────────────────────────────
type Doctor struct {
	DoctorID     string    `json:"doctor_id"`
	UserID       string    `json:"user_id"`
	DepartmentID string    `json:"department_id"`
	FullName     string    `json:"full_name"`
	Specialty    string    `json:"specialty"`
	Bio          string    `json:"bio"`
	Available    bool      `json:"available"`
	Department   string    `json:"department,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateDoctorRequest struct {
	DepartmentID string `json:"department_id" binding:"required"`
	FullName     string `json:"full_name"     binding:"required"`
	Specialty    string `json:"specialty"     binding:"required"`
	Bio          string `json:"bio"`
}

// ── AVAILABILITY SLOT ─────────────────────────────────────────
type AvailabilitySlot struct {
	SlotID    string    `json:"slot_id"`
	DoctorID  string    `json:"doctor_id"`
	SlotDate  string    `json:"slot_date"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	IsBooked  bool      `json:"is_booked"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateSlotRequest struct {
	DoctorID  string `json:"doctor_id"  binding:"required"`
	SlotDate  string `json:"slot_date"  binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time"   binding:"required"`
}

// ── APPOINTMENT ───────────────────────────────────────────────
type Appointment struct {
	AppointmentID   string    `json:"appointment_id"`
	UserID          string    `json:"user_id"`
	DoctorID        string    `json:"doctor_id"`
	DepartmentID    string    `json:"department_id"`
	SlotID          string    `json:"slot_id"`
	AppointmentDate string    `json:"appointment_date"`
	AppointmentTime string    `json:"appointment_time"`
	Reason          string    `json:"reason"`
	Status          string    `json:"status"`
	Notes           string    `json:"notes"`
	PatientName     string    `json:"patient_name,omitempty"`
	DoctorName      string    `json:"doctor_name,omitempty"`
	DepartmentName  string    `json:"department_name,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateAppointmentRequest struct {
	DoctorID     string `json:"doctor_id"     binding:"required"`
	DepartmentID string `json:"department_id" binding:"required"`
	SlotID       string `json:"slot_id"       binding:"required"`
	Reason       string `json:"reason"`
}

type UpdateAppointmentRequest struct {
	Status string `json:"status" binding:"required"`
	Notes  string `json:"notes"`
}

// ── QUEUE ─────────────────────────────────────────────────────
type QueueEntry struct {
	QueueID       string     `json:"queue_id"`
	AppointmentID string     `json:"appointment_id"`
	UserID        string     `json:"user_id"`
	DoctorID      string     `json:"doctor_id"`
	QueueNumber   int        `json:"queue_number"`
	Position      int        `json:"position"`
	Status        string     `json:"status"`
	PatientName   string     `json:"patient_name,omitempty"`
	DoctorName    string     `json:"doctor_name,omitempty"`
	CheckedInAt   *time.Time `json:"checked_in_at"`
	CalledAt      *time.Time `json:"called_at"`
	ServedAt      *time.Time `json:"served_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

type CheckInRequest struct {
	AppointmentID string `json:"appointment_id" binding:"required"`
}

type UpdateQueueRequest struct {
	Status string `json:"status" binding:"required"`
}

// ── GENERIC RESPONSES ─────────────────────────────────────────
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

