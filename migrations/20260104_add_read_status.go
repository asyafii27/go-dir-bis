package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

// AddReadStatusToMessages menambahkan status 'read' ke enum status di chat_messages
func AddReadStatusToMessages(db *gorm.DB) error {
	sql := "ALTER TABLE chat_messages MODIFY COLUMN status ENUM('sent','delivered','read','edited') DEFAULT 'sent'"

	if err := db.Exec(sql).Error; err != nil {
		return fmt.Errorf("gagal menambahkan status 'read': %v", err)
	}

	fmt.Println("✓ Berhasil menambahkan status 'read' ke chat_messages")
	return nil
}
