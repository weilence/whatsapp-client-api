package whatsapp

import (
	"database/sql"
)

func createMessagesTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS messages (
		id TEXT PRIMARY KEY,
		our_jid TEXT,
		their_jid TEXT,
		message TEXT,
		revoked INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	return err
}

func insertMessage(db *sql.DB, id, outJid, theirJid, message string) error {
	_, err := db.Exec(`INSERT INTO messages (id, our_jid, their_jid, message) VALUES (?, ?, ?, ?)`, id, outJid, theirJid, message)
	return err
}

func revokeMessage(db *sql.DB, id string) error {
	_, err := db.Exec(`UPDATE messages SET revoked = ? WHERE id = ?`, 1, id)
	return err
}
