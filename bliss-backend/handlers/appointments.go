package handlers

import (
	"context"
	"net/http"

	"bliss-backend/db"
	"bliss-backend/models"

	"github.com/gin-gonic/gin"
)

func CreateAppointment(c *gin.Context) {
	userID := c.GetString("user_id")

	var req models.CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var isBooked bool
	err := db.DB.QueryRow(context.Background(),
		`SELECT is_booked FROM availability_slots WHERE slot_id = $1`, req.SlotID,
	).Scan(&isBooked)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Slot not found"})
		return
	}
	if isBooked {
		c.JSON(http.StatusConflict, gin.H{"error": "This slot has already been booked"})
		return
	}

	var appt models.Appointment
	query := `
		INSERT INTO appointments (user_id, doctor_id, department_id, slot_id, appointment_date, appointment_time, reason, status)
		SELECT $1, $2, $3, $4, slot_date, start_time, $5, 'pending'
		FROM availability_slots WHERE slot_id = $4
		RETURNING appointment_id, user_id, doctor_id, department_id, slot_id,
		          appointment_date, appointment_time, reason, status, created_at, updated_at
	`
	err = db.DB.QueryRow(context.Background(), query,
		userID, req.DoctorID, req.DepartmentID, req.SlotID, req.Reason,
	).Scan(
		&appt.AppointmentID, &appt.UserID, &appt.DoctorID, &appt.DepartmentID, &appt.SlotID,
		&appt.AppointmentDate, &appt.AppointmentTime, &appt.Reason, &appt.Status,
		&appt.CreatedAt, &appt.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create appointment"})
		return
	}

	db.DB.Exec(context.Background(),
		`UPDATE availability_slots SET is_booked = TRUE WHERE slot_id = $1`, req.SlotID,
	)

	go queueNotification(appt.AppointmentID, userID, "confirmation")

	c.JSON(http.StatusCreated, models.SuccessResponse{Message: "Appointment booked successfully", Data: appt})
}

func GetAppointments(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")

	if role == "admin" || role == "receptionist" {
		query := `
			SELECT a.appointment_id, a.user_id, a.doctor_id, a.department_id, a.slot_id,
			       a.appointment_date, a.appointment_time, a.reason, a.status, a.notes,
			       u.full_name, d.full_name, dep.name, a.created_at, a.updated_at
			FROM appointments a
			LEFT JOIN users u       ON a.user_id       = u.user_id
			LEFT JOIN doctors d     ON a.doctor_id     = d.doctor_id
			LEFT JOIN departments dep ON a.department_id = dep.department_id
			ORDER BY a.appointment_date DESC, a.appointment_time DESC
		`
		rows, err := db.DB.Query(context.Background(), query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointments"})
			return
		}
		defer rows.Close()

		var appts []models.Appointment
		for rows.Next() {
			var a models.Appointment
			if err := rows.Scan(
				&a.AppointmentID, &a.UserID, &a.DoctorID, &a.DepartmentID, &a.SlotID,
				&a.AppointmentDate, &a.AppointmentTime, &a.Reason, &a.Status, &a.Notes,
				&a.PatientName, &a.DoctorName, &a.DepartmentName,
				&a.CreatedAt, &a.UpdatedAt,
			); err == nil {
				appts = append(appts, a)
			}
		}
		if appts == nil {
			appts = []models.Appointment{}
		}
		c.JSON(http.StatusOK, appts)
		return
	}

	query := `
		SELECT a.appointment_id, a.user_id, a.doctor_id, a.department_id, a.slot_id,
		       a.appointment_date, a.appointment_time, a.reason, a.status, a.notes,
		       d.full_name, dep.name, a.created_at, a.updated_at
		FROM appointments a
		LEFT JOIN doctors d       ON a.doctor_id     = d.doctor_id
		LEFT JOIN departments dep ON a.department_id = dep.department_id
		WHERE a.user_id = $1
		ORDER BY a.appointment_date DESC
	`
	rows, err := db.DB.Query(context.Background(), query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointments"})
		return
	}
	defer rows.Close()

	var appts []models.Appointment
	for rows.Next() {
		var a models.Appointment
		if err := rows.Scan(
			&a.AppointmentID, &a.UserID, &a.DoctorID, &a.DepartmentID, &a.SlotID,
			&a.AppointmentDate, &a.AppointmentTime, &a.Reason, &a.Status, &a.Notes,
			&a.DoctorName, &a.DepartmentName,
			&a.CreatedAt, &a.UpdatedAt,
		); err == nil {
			appts = append(appts, a)
		}
	}
	if appts == nil {
		appts = []models.Appointment{}
	}
	c.JSON(http.StatusOK, appts)
}

func GetAppointment(c *gin.Context) {
	appointmentID := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("role")

	var a models.Appointment
	query := `
		SELECT a.appointment_id, a.user_id, a.doctor_id, a.department_id, a.slot_id,
		       a.appointment_date, a.appointment_time, a.reason, a.status, a.notes,
		       u.full_name, d.full_name, dep.name, a.created_at, a.updated_at
		FROM appointments a
		LEFT JOIN users u         ON a.user_id       = u.user_id
		LEFT JOIN doctors d       ON a.doctor_id     = d.doctor_id
		LEFT JOIN departments dep ON a.department_id = dep.department_id
		WHERE a.appointment_id = $1
	`
	err := db.DB.QueryRow(context.Background(), query, appointmentID).Scan(
		&a.AppointmentID, &a.UserID, &a.DoctorID, &a.DepartmentID, &a.SlotID,
		&a.AppointmentDate, &a.AppointmentTime, &a.Reason, &a.Status, &a.Notes,
		&a.PatientName, &a.DoctorName, &a.DepartmentName,
		&a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	if role == "patient" && a.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, a)
}

func UpdateAppointment(c *gin.Context) {
	appointmentID := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("role")

	var req models.UpdateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if role == "patient" {
		if req.Status != "cancelled" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Patients can only cancel appointments"})
			return
		}
		var ownerID string
		db.DB.QueryRow(context.Background(),
			`SELECT user_id FROM appointments WHERE appointment_id = $1`, appointmentID,
		).Scan(&ownerID)
		if ownerID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	_, err := db.DB.Exec(context.Background(),
		`UPDATE appointments SET status = $1, notes = $2 WHERE appointment_id = $3`,
		req.Status, req.Notes, appointmentID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update appointment"})
		return
	}

	if req.Status == "cancelled" {
		db.DB.Exec(context.Background(),
			`UPDATE availability_slots SET is_booked = FALSE
			 WHERE slot_id = (SELECT slot_id FROM appointments WHERE appointment_id = $1)`,
			appointmentID,
		)
		go queueNotification(appointmentID, userID, "cancellation")
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Appointment updated successfully"})
}

func GetDepartments(c *gin.Context) {
	rows, err := db.DB.Query(context.Background(),
		`SELECT department_id, name, description, created_at FROM departments ORDER BY name`,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch departments"})
		return
	}
	defer rows.Close()

	var depts []models.Department
	for rows.Next() {
		var d models.Department
		if err := rows.Scan(&d.DepartmentID, &d.Name, &d.Description, &d.CreatedAt); err == nil {
			depts = append(depts, d)
		}
	}
	if depts == nil {
		depts = []models.Department{}
	}
	c.JSON(http.StatusOK, depts)
}

func GetDoctors(c *gin.Context) {
	deptID := c.Query("department_id")

	query := `
		SELECT d.doctor_id, COALESCE(d.user_id::text, ''), d.department_id, d.full_name,
		       d.specialty, COALESCE(d.bio, ''), d.available, COALESCE(dep.name, ''), d.created_at
		FROM doctors d
		LEFT JOIN departments dep ON d.department_id = dep.department_id
		WHERE d.available = TRUE
	`
	args := []interface{}{}
	if deptID != "" {
		query += ` AND d.department_id = $1`
		args = append(args, deptID)
	}
	query += ` ORDER BY d.full_name`

	rows, err := db.DB.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch doctors"})
		return
	}
	defer rows.Close()

	var doctors []models.Doctor
	for rows.Next() {
		var d models.Doctor
		if err := rows.Scan(
			&d.DoctorID, &d.UserID, &d.DepartmentID, &d.FullName,
			&d.Specialty, &d.Bio, &d.Available, &d.Department, &d.CreatedAt,
		); err == nil {
			doctors = append(doctors, d)
		}
	}
	if doctors == nil {
		doctors = []models.Doctor{}
	}
	c.JSON(http.StatusOK, doctors)
}

func GetSlots(c *gin.Context) {
	doctorID := c.Query("doctor_id")
	date := c.Query("date")

	if doctorID == "" || date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "doctor_id and date are required"})
		return
	}

	rows, err := db.DB.Query(context.Background(),
		`SELECT slot_id, doctor_id, slot_date, start_time, end_time, is_booked, created_at
		 FROM availability_slots
		 WHERE doctor_id = $1 AND slot_date = $2 AND is_booked = FALSE
		 ORDER BY start_time`,
		doctorID, date,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch slots"})
		return
	}
	defer rows.Close()

	var slots []models.AvailabilitySlot
	for rows.Next() {
		var s models.AvailabilitySlot
		if err := rows.Scan(
			&s.SlotID, &s.DoctorID, &s.SlotDate, &s.StartTime, &s.EndTime, &s.IsBooked, &s.CreatedAt,
		); err == nil {
			slots = append(slots, s)
		}
	}
	if slots == nil {
		slots = []models.AvailabilitySlot{}
	}
	c.JSON(http.StatusOK, slots)
}

func CreateSlot(c *gin.Context) {
	var req models.CreateSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var slot models.AvailabilitySlot
	err := db.DB.QueryRow(context.Background(),
		`INSERT INTO availability_slots (doctor_id, slot_date, start_time, end_time)
		 VALUES ($1, $2, $3, $4)
		 RETURNING slot_id, doctor_id, slot_date, start_time, end_time, is_booked, created_at`,
		req.DoctorID, req.SlotDate, req.StartTime, req.EndTime,
	).Scan(&slot.SlotID, &slot.DoctorID, &slot.SlotDate, &slot.StartTime, &slot.EndTime, &slot.IsBooked, &slot.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create slot"})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{Message: "Slot created successfully", Data: slot})
}

