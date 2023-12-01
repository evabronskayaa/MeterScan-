package service

func (s *Service) NotifyUsers() {
	/*var users []schema.User
	now := time.Now()
	dayOfMonth := now.Day()
	hour := now.Hour()

	if err := s.DB.Model(&schema.User{}).
		Preload("Settings").
		Joins("inner join user_settings on user_settings.user_id = users.id").
		Where("users.verified_at IS NOT NULL").
		Where("user_settings.notification_day_of_month = ?", dayOfMonth).
		Where("user_settings.notification_hour = ?", hour).
		Find(&users).Error; err != nil {
		log.Printf("cron notify users: %v", err)
		return
	}

	for _, user := range users {
		if err := user.SendMessageToMail(s.MailClient, "Пора вносить показания", "notification.gohtml", nil); err != nil {
			log.Printf("err sending mail for userId=%v: %v", user.ID, err)
		}
	}*/
}
