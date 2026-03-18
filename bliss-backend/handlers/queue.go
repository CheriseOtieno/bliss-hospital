package handlers

import (
	"context"
	"net/http"
	"time"

	"bliss-backend/db"
	"bliss-backend/models"

	"github.com/gin-gonic/gin"
)

func CheckIn(c *gin.Context) {
	userID := c.GetString("user_id")

	var req models.CheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var apptUserID, doctorID, apptStatus string
	err := db.DB.QueryRow(context.Background(),
		`SELECT user_id, doctor_id, status FROM appointments WHERE appointment_id = $1`,
		req.AppointmentID,
	).Scan(&apptUserID, &doctorID, &apptStatus)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}
	if apptUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "This appointment does not belong to you"})
		return
	}
	if apptStatus == "cancelled" || apptStatus == "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot check in for a " + apptStatus + " appointment"})
		return
	}

	var existingID string
	db.DB.QueryRow(context.Background(),
		`SELECT queue_id FROM queue WHERE appointment_id = $1`, req.AppointmentID,
	).Scan(&existingID)
	if existingID != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "Already checked in for this appointment"})
		return
	}

	var maxQueueNum int
	db.DB.QueryRow(context.Background(),
		`SELECT COALESCE(MAX(queue_number), 0) FROM queue WHERE doctor_id = $1 AND DATE(created_at) = CURRENT_DATE`,
		doctorID,
	).Scan(&maxQueueNum)

	var position int
	db.DB.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM queue WHERE doctor_id = $1 AND status IN ('waiting','called') AND DATE(created_at) = CURRENT_DATE`,
		doctorID,
	).Scan(&position)

	now := time.Now()
	var entry models.QueueEntry
	err = db.DB.QueryRow(context.Background(),
		`INSERT INTO queue (appointment_id, user_id, doctor_id, queue_number, position, status, checked_in_at)
		 VALUES ($1, $2, $3, $4, $5, 'waiting', $6)
		 RETURNING queue_id, appointment_id, user_id, doctor_id, queue_number, position, status, checked_in_at, created_at`,
		req.AppointmentID, userID, doctorID, maxQueueNum+1, position+1, now,
	).Scan(
		&entry.QueueID, &entry.AppointmentID, &entry.UserID, &entry.DoctorID,
		&entry.QueueNumber, &entry.Position, &entry.Status, &entry.CheckedInAt, &entry.CreatedAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check in"})
		return
	}

	db.DB.Exec(context.Background(),
		`UPDATE appointments SET status = 'confirmed' WHERE appointment_id = $1`, req.AppointmentID,
	)

	c.JSON(http.StatusCreated, models.SuccessResponse{Message: "Check-in successful", Data: entry})
}

func GetMyQueueStatus(c *gin.Context) {
	userID := c.GetString("user_id")

	var entry models.QueueEntry
	err := db.DB.QueryRow(context.Background(),
		`SELECT q.queue_id, q.appointment_id, q.user_id, q.doctor_id,
		        q.queue_number, q.position, q.status,
		        COALESCE(d.full_name, '') AS doctor_name,
		        q.checked_in_at, q.called_at, q.served_at, q.created_at
		 FROM queue q
		 LEFT JOIN doctors d ON q.doctor_id = d.doctor_id
		 WHERE q.user_id = $1
		   AND q.status IN ('waiting', 'called', 'serving')
		   AND DATE(q.created_at) = CURRENT_DATE
		 ORDER BY q.created_at DESC LIMIT 1`,
		userID,
	).Scan(
		&entry.QueueID, &entry.AppointmentID, &entry.UserID, &entry.DoctorID,
		&entry.QueueNumber, &entry.Position, &entry.Status, &entry.DoctorName,
		&entry.CheckedInAt, &entry.CalledAt, &entry.ServedAt, &entry.CreatedAt,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active queue entry found for today"})
		return
	}

	var livePosition int
	db.DB.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM queue
		 WHERE doctor_id = $1 AND status IN ('waiting','called')
		   AND queue_number < $2 AND DATE(created_at) = CURRENT_DATE`,
		entry.DoctorID, entry.QueueNumber,
	).Scan(&livePosition)
	entry.Position = livePosition + 1

	c.JSON(http.StatusOK, entry)
}

func GetQueue(c *gin.Context) {
	doctorID := c.Query("doctor_id")

	query := `
		SELECT q.queue_id, q.appointment_id, q.user_id, q.doctor_id,
		       q.queue_number, q.position, q.status,
		       COALESCE(u.full_name, '') AS patient_name,
		       COALESCE(d.full_name, '') AS doctor_name,
		       q.checked_in_at, q.called_at, q.served_at, q.created_at
		FROM queue q
		LEFT JOIN users u   ON q.user_id   = u.user_id
		LEFT JOIN doctors d ON q.doctor_id = d.doctor_id
		WHERE DATE(q.created_at) = CURRENT_DATE
	`
	args := []interface{}{}
	if doctorID != "" {
		query += ` AND q.doctor_id = $1`
		args = append(args, doctorID)
	}
	query += ` ORDER BY q.queue_number ASC`

	rows, err := db.DB.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch queue"})
		return
	}
	defer rows.Close()

	var entries []models.QueueEntry
	for rows.Next() {
		var e models.QueueEntry
		if err := rows.Scan(
			&e.QueueID, &e.AppointmentID, &e.UserID, &e.DoctorID,
			&e.QueueNumber, &e.Position, &e.Status,
			&e.PatientName, &e.DoctorName,
			&e.CheckedInAt, &e.CalledAt, &e.ServedAt, &e.CreatedAt,
		); err == nil {
			entries = append(entries, e)
		}
	}
	if entries == nil {
		entries = []models.QueueEntry{}
	}
	c.JSON(http.StatusOK, entries)
}

func UpdateQueueEntry(c *gin.Context) {
	queueID := c.Param("id")

	var req models.UpdateQueueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()

	switch req.Status {
	case "called":
		db.DB.Exec(context.Background(),
			`UPDATE queue SET status = 'called', called_at = $1 WHERE queue_id = $2`, now, queueID,
		)
		var userID, apptID string
		db.DB.QueryRow(context.Background(),
			`SELECT user_id, appointment_id FROM queue WHERE queue_id = $1`, queueID,
		).Scan(&userID, &apptID)
		go queueNotification(apptID, userID, "queue_call")

	case "serving":
		db.DB.Exec(context.Background(),
			`UPDATE queue SET status = 'serving' WHERE queue_id = $1`, queueID,
		)

	case "done":
		db.DB.Exec(context.Background(),
			`UPDATE queue SET status = 'done', served_at = $1 WHERE queue_id = $2`, now, queueID,
		)
		db.DB.Exec(context.Background(),
			`UPDATE appointments SET status = 'completed'
			 WHERE appointment_id = (SELECT appointment_id FROM queue WHERE queue_id = $1)`, queueID,
		)

	case "skipped":
		db.DB.Exec(context.Background(),
			`UPDATE queue SET status = 'skipped' WHERE queue_id = $1`, queueID,
		)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Use: called, serving, done, skipped"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Queue entry updated successfully"})
}

func GetQueueStats(c *gin.Context) {
	var waiting, serving, done, skipped int
	db.DB.QueryRow(context.Background(),
		`SELECT
			COUNT(*) FILTER (WHERE status = 'waiting'),
			COUNT(*) FILTER (WHERE status IN ('called','serving')),
			COUNT(*) FILTER (WHERE status = 'done'),
			COUNT(*) FILTER (WHERE status = 'skipped')
		 FROM queue WHERE DATE(created_at) = CURRENT_DATE`,
	).Scan(&waiting, &serving, &done, &skipped)

	c.JSON(http.StatusOK, gin.H{
		"waiting": waiting,
		"serving": serving,
		"done":    done,
		"skipped": skipped,
		"total":   waiting + serving + done + skipped,
	})
}

