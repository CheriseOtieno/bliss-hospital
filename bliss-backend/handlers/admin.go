package handlers

import (
	"context"
	"net/http"

	"bliss-backend/db"
	"bliss-backend/models"

	"github.com/gin-gonic/gin"
)

func CreateDoctor(c *gin.Context) {
	var req models.CreateDoctorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var doc models.Doctor
	err := db.DB.QueryRow(context.Background(),
		`INSERT INTO doctors (department_id, full_name, specialty, bio)
		 VALUES ($1, $2, $3, $4)
		 RETURNING doctor_id, COALESCE(user_id::text, ''), department_id, full_name, specialty, COALESCE(bio, ''), available, created_at`,
		req.DepartmentID, req.FullName, req.Specialty, req.Bio,
	).Scan(&doc.DoctorID, &doc.UserID, &doc.DepartmentID, &doc.FullName, &doc.Specialty, &doc.Bio, &doc.Available, &doc.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create doctor"})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{Message: "Doctor created successfully", Data: doc})
}

func UpdateDoctor(c *gin.Context) {
	doctorID := c.Param("id")

	var body struct {
		Available bool   `json:"available"`
		Specialty string `json:"specialty"`
		Bio       string `json:"bio"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.DB.Exec(context.Background(),
		`UPDATE doctors SET available = $1, specialty = $2, bio = $3 WHERE doctor_id = $4`,
		body.Available, body.Specialty, body.Bio, doctorID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update doctor"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Doctor updated successfully"})
}

func GetDashboardStats(c *gin.Context) {
	var totalPatients, totalDoctors, todayAppointments, pendingAppointments, waitingQueue int

	db.DB.QueryRow(context.Background(), `SELECT COUNT(*) FROM users WHERE role = 'patient'`).Scan(&totalPatients)
	db.DB.QueryRow(context.Background(), `SELECT COUNT(*) FROM doctors WHERE available = TRUE`).Scan(&totalDoctors)
	db.DB.QueryRow(context.Background(), `SELECT COUNT(*) FROM appointments WHERE appointment_date = CURRENT_DATE`).Scan(&todayAppointments)
	db.DB.QueryRow(context.Background(), `SELECT COUNT(*) FROM appointments WHERE status = 'pending'`).Scan(&pendingAppointments)
	db.DB.QueryRow(context.Background(), `SELECT COUNT(*) FROM queue WHERE status = 'waiting' AND DATE(created_at) = CURRENT_DATE`).Scan(&waitingQueue)

	c.JSON(http.StatusOK, gin.H{
		"total_patients":       totalPatients,
		"total_doctors":        totalDoctors,
		"today_appointments":   todayAppointments,
		"pending_appointments": pendingAppointments,
		"waiting_queue":        waitingQueue,
	})
}

func GetUsers(c *gin.Context) {
	rows, err := db.DB.Query(context.Background(),
		`SELECT user_id, full_name, email, COALESCE(phone, ''), role, is_active, created_at
		 FROM users ORDER BY created_at DESC`,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.UserID, &u.FullName, &u.Email, &u.Phone, &u.Role, &u.IsActive, &u.CreatedAt); err == nil {
			users = append(users, u)
		}
	}
	if users == nil {
		users = []models.User{}
	}
	c.JSON(http.StatusOK, users)
}

