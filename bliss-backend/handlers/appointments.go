package handlers

import (
	"context"
	"net/http"

	"bliss-backend/db"
	"bliss-backend/models"

	"github.com/gin-gonic/gin"
)

/*
	=========================================================
	  CREATE APPOINTMENT

=========================================================
*/
func CreateAppointment(c *gin.Context) {
	userID := c.GetString("user_id")

	var req models.CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var isBooked bool
	err := db.DB.QueryRow(context.Background(),
		`SELECT is_booked FROM availability_slots WHERE slot_id = $1`,
		req.SlotID,
	).Scan(&isBooked)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Slot not found"})
		return
	}

	if isBooked {
		c.JSON(http.StatusConflict, gin.H{"error": "This slot is already booked"})
		return
	}

	var appt models.Appointment

	query := `
		INSERT INTO appointments (
			user_id, doctor_id, department_id, branch_id, slot_id,
			appointment_date, appointment_time, reason, status
		)
		SELECT $1, $2, $3, $4, $5, slot_date, start_time, $6, 'pending'
		FROM availability_slots
		WHERE slot_id = $5
		RETURNING appointment_id, user_id, doctor_id, department_id, branch_id, slot_id,
				  appointment_date, appointment_time, reason, status, created_at, updated_at
	`

	err = db.DB.QueryRow(context.Background(), query,
		userID,
		req.DoctorID,
		req.DepartmentID,
		req.BranchID,
		req.SlotID,
		req.Reason,
	).Scan(
		&appt.AppointmentID,
		&appt.UserID,
		&appt.DoctorID,
		&appt.DepartmentID,
		&appt.BranchID,
		&appt.SlotID,
		&appt.AppointmentDate,
		&appt.AppointmentTime,
		&appt.Reason,
		&appt.Status,
		&appt.CreatedAt,
		&appt.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create appointment"})
		return
	}

	_, _ = db.DB.Exec(context.Background(),
		`UPDATE availability_slots SET is_booked = TRUE WHERE slot_id = $1`,
		req.SlotID,
	)

	c.JSON(http.StatusCreated, appt)
}

/*
	=========================================================
	  GET ALL APPOINTMENTS

=========================================================
*/
func GetAppointments(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")

	baseQuery := `
		SELECT a.appointment_id, a.user_id, a.doctor_id, a.department_id, a.branch_id, a.slot_id,
		       a.appointment_date, a.appointment_time, a.reason, a.status,
		       d.full_name, dep.name, b.branch_name
		FROM appointments a
		LEFT JOIN doctors d ON a.doctor_id = d.doctor_id
		LEFT JOIN departments dep ON a.department_id = dep.department_id
		LEFT JOIN branches b ON a.branch_id = b.branch_id
	`

	var rows interface {
		Next() bool
		Scan(...any) error
		Close()
	}
	var err error

	if role == "admin" || role == "receptionist" {
		rows, err = db.DB.Query(context.Background(),
			baseQuery+` ORDER BY a.appointment_date DESC`,
		)
	} else {
		rows, err = db.DB.Query(context.Background(),
			baseQuery+` WHERE a.user_id = $1 ORDER BY a.appointment_date DESC`,
			userID,
		)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointments"})
		return
	}
	defer rows.Close()

	var appts []models.Appointment

	for rows.Next() {
		var a models.Appointment

		if err := rows.Scan(
			&a.AppointmentID,
			&a.UserID,
			&a.DoctorID,
			&a.DepartmentID,
			&a.BranchID,
			&a.SlotID,
			&a.AppointmentDate,
			&a.AppointmentTime,
			&a.Reason,
			&a.Status,
			&a.DoctorName,
			&a.DepartmentName,
			&a.BranchName,
		); err != nil {
			continue
		}

		appts = append(appts, a)
	}

	c.JSON(http.StatusOK, appts)
}

/*
	=========================================================
	  GET SINGLE APPOINTMENT

=========================================================
*/
func GetAppointment(c *gin.Context) {
	appointmentID := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("role")

	var a models.Appointment

	query := `
		SELECT a.appointment_id, a.user_id, a.doctor_id, a.department_id, a.branch_id, a.slot_id,
		       a.appointment_date, a.appointment_time, a.reason, a.status,
		       u.full_name, d.full_name, dep.name, b.branch_name,
		       a.created_at, a.updated_at
		FROM appointments a
		LEFT JOIN users u ON a.user_id = u.user_id
		LEFT JOIN doctors d ON a.doctor_id = d.doctor_id
		LEFT JOIN departments dep ON a.department_id = dep.department_id
		LEFT JOIN branches b ON a.branch_id = b.branch_id
		WHERE a.appointment_id = $1
	`

	err := db.DB.QueryRow(context.Background(), query, appointmentID).Scan(
		&a.AppointmentID,
		&a.UserID,
		&a.DoctorID,
		&a.DepartmentID,
		&a.BranchID,
		&a.SlotID,
		&a.AppointmentDate,
		&a.AppointmentTime,
		&a.Reason,
		&a.Status,
		&a.PatientName,
		&a.DoctorName,
		&a.DepartmentName,
		&a.BranchName,
		&a.CreatedAt,
		&a.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	if role == "patient" && a.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	c.JSON(http.StatusOK, a)
}

/*
	=========================================================
	  UPDATE APPOINTMENT

=========================================================
*/
func UpdateAppointment(c *gin.Context) {
	appointmentID := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("role")

	var req models.UpdateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ownerID string
	err := db.DB.QueryRow(context.Background(),
		`SELECT user_id FROM appointments WHERE appointment_id = $1`,
		appointmentID,
	).Scan(&ownerID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	if role == "patient" && ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	_, err = db.DB.Exec(context.Background(),
		`UPDATE appointments SET status = $1, notes = $2 WHERE appointment_id = $3`,
		req.Status, req.Notes, appointmentID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update appointment"})
		return
	}

	if req.Status == "cancelled" {
		db.DB.Exec(context.Background(),
			`UPDATE availability_slots 
			 SET is_booked = FALSE
			 WHERE slot_id = (
				SELECT slot_id FROM appointments WHERE appointment_id = $1
			 )`,
			appointmentID,
		)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment updated"})
}

/*
	=========================================================
	  GET BRANCHES

=========================================================
*/
func GetBranches(c *gin.Context) {
	rows, err := db.DB.Query(context.Background(),
		`SELECT branch_id, branch_name FROM branches ORDER BY branch_name`,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch branches"})
		return
	}
	defer rows.Close()

	var branches []models.Branch

	for rows.Next() {
		var b models.Branch
		_ = rows.Scan(&b.BranchID, &b.Name) // IGNORE: name field in Branch struct is actually branch_name in DB
		branches = append(branches, b)
	}

	c.JSON(http.StatusOK, branches)
}

/*
	=========================================================
	  GET DEPARTMENTS

=========================================================
*/
func GetDepartments(c *gin.Context) {
	rows, err := db.DB.Query(context.Background(),
		`SELECT department_id, name, description FROM departments ORDER BY name`,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch departments"})
		return
	}
	defer rows.Close()

	var depts []models.Department

	for rows.Next() {
		var d models.Department
		_ = rows.Scan(&d.DepartmentID, &d.Name, &d.Description)
		depts = append(depts, d)
	}

	c.JSON(http.StatusOK, depts)
}

/*
	=========================================================
	  GET DOCTORS

=========================================================
*/
func GetDoctors(c *gin.Context) {
	deptID := c.Query("department_id")

	baseQuery := `
		SELECT doctor_id, department_id, full_name, specialty, bio
		FROM doctors
		WHERE available = TRUE
	`

	var rows interface {
		Next() bool
		Scan(...any) error
		Close()
	}
	var err error

	if deptID != "" {
		rows, err = db.DB.Query(context.Background(), baseQuery+` AND department_id = $1`, deptID)
	} else {
		rows, err = db.DB.Query(context.Background(), baseQuery)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch doctors"})
		return
	}
	defer rows.Close()

	var doctors []models.Doctor

	for rows.Next() {
		var d models.Doctor
		_ = rows.Scan(&d.DoctorID, &d.DepartmentID, &d.FullName, &d.Specialty, &d.Bio)
		doctors = append(doctors, d)
	}

	c.JSON(http.StatusOK, doctors)
}

/*
	=========================================================
	  GET SLOTS

=========================================================
*/
func GetSlots(c *gin.Context) {
	doctorID := c.Query("doctor_id")

	if doctorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "doctor_id is required"})
		return
	}

	rows, err := db.DB.Query(context.Background(),
		`SELECT slot_id, slot_date, start_time, end_time, is_booked, created_at
		 FROM availability_slots
		 WHERE doctor_id = $1 AND is_booked = FALSE
		 ORDER BY slot_date, start_time`,
		doctorID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch slots"})
		return
	}
	defer rows.Close()

	var slots []models.AvailabilitySlot

	for rows.Next() {
		var s models.AvailabilitySlot
		_ = rows.Scan(
			&s.SlotID,
			&s.SlotDate,
			&s.StartTime,
			&s.EndTime,
			&s.IsBooked,
			&s.CreatedAt,
		)
		slots = append(slots, s)
	}

	c.JSON(http.StatusOK, slots)
}
