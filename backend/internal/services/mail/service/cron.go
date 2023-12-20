package service

import (
	"backend/internal/proto"
	"context"
	"log"
	"time"
)

func (s *Service) NotifyUsers() {
	now := time.Now()
	dayOfMonth := now.Day()
	hour := now.Hour()
	response, err := s.DB.GetEmailsForNotification(context.Background(), &proto.GetEmailsForNotificationRequest{
		Day:  int32(dayOfMonth),
		Hour: int32(hour),
	})
	if err != nil {
		log.Printf("cron notify users: %v", err)
		return
	}

	if err := s.Mail.SendPlainMessage("Пора вносить показания", "notification.gohtml", response.Emails...); err != nil {
		log.Printf("err sending mail for emails=%v: %v", response.Emails, err)
	}
}
