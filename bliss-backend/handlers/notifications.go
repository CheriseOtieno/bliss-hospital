package handlers

import (
	"context"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"bliss-backend/db"

	"github.com/gin-gonic/gin"
)

func queueNotification(appointmentID, userID, notifType string) {
	var message string
	switch notifType {
	case "confirmation":
		message = "Your appointment at Bliss Hospital has been confirmed. Please arrive 10 minutes early."
	case "reminder":
		message = "Reminder: You have an appointment at Bliss Hospital tomorrow."
	case "cancellation":
		message = "Your Bliss Hospital appointment has been cancelled. Visit our portal to rebook."
	case "queue_call":
		message = "It is now your turn! Please proceed to the consultation room."
	case "delay":
		message = "We apologize — your appointment is running slightly behind schedule. Please remain seated."
	default:
		message = "You have a notification from Bliss Hospital."
	}

	var email, phone, fullName string
	err := db.DB.QueryRow(context.Background(),
		`SELECT email, COALESCE(phone, ''), full_name FROM users WHERE user_id = $1`, userID,
	).Scan(&email, &phone, &fullName)
	if err != nil {
		log.Printf("Notification: failed to fetch user %s: %v", userID, err)
		return
	}

	message = "Dear " + fullName + ", " + message

	var apptIDParam interface{}
	if appointmentID != "" {
		apptIDParam = appointmentID
	}

	db.DB.Exec(context.Background(),
		`INSERT INTO notifications (user_id, appointment_id, type, channel, message, status)
		 VALUES ($1, $2, $3, 'email', $4, 'pending')`,
		userID, apptIDParam, notifType, message,
	)

	if email != "" {
		if err := sendEmail(email, notifType, message); err != nil {
			log.Printf("Notification: email to %s failed: %v", email, err)
			db.DB.Exec(context.Background(),
				`UPDATE notifications SET status = 'failed'
				 WHERE user_id = $1 AND type = $2 AND status = 'pending'`,
				userID, notifType,
			)
		} else {
			db.DB.Exec(context.Background(),
				`UPDATE notifications SET status = 'sent', sent_at = NOW()
				 WHERE user_id = $1 AND type = $2 AND status = 'pending'`,
				userID, notifType,
			)
		}
	}
}

func sendEmail(to, subject, body string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")

	if host == "" || user == "" || pass == "" {
		log.Println("Notification: SMTP not configured, skipping email send")
		return nil
	}

	auth := smtp.PlainAuth("", user, pass, host)
	msg := []byte(
		"To: " + to + "\r\n" +
			"From: Bliss Hospital <" + user + ">\r\n" +
			"Subject: Bliss Hospital - " + subject + "\r\n" +
			"MIME-version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n" +
			body + "\r\n",
	)
	return smtp.SendMail(host+":"+port, auth, user, []string{to}, msg)
}

func GetNotifications(c *gin.Context) {
	userID := c.GetString("user_id")

	rows, err := db.DB.Query(context.Background(),
		`SELECT notification_id, type, channel, message, status, created_at
		 FROM notifications WHERE user_id = $1
		 ORDER BY created_at DESC LIMIT 20`,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}
	defer rows.Close()

	var notifs []map[string]interface{}
	for rows.Next() {
		var id, notifType, channel, message, status, createdAt string
		if err := rows.Scan(&id, &notifType, &channel, &message, &status, &createdAt); err == nil {
			notifs = append(notifs, map[string]interface{}{
				"notification_id": id,
				"type":            notifType,
				"channel":         channel,
				"message":         message,
				"status":          status,
				"created_at":      createdAt,
			})
		}
	}
	if notifs == nil {
		notifs = []map[string]interface{}{}
	}
	c.JSON(http.StatusOK, notifs)
}

