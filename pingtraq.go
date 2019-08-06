package pingtraq

import (
	"database/sql"
	"net/http"
)

func Init(db string) error {
	if err := initDB(db); err != nil {
		return err
	}

	return nil
}

func IsPing(name string) string {
	var id string
	err := queryRow("SELECT id FROM ping WHERE name = ?", name).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return ""
	}
	if id == "" {
		return ""
	}

	return id
}

func AddPing(name string) error {
	stmt, err := prepare("INSERT INTO ping VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(newUUID(), name, now())
	if err != nil {
		return err
	}

	return nil
}

func AddPingRecord(id string, r *http.Request) error {
	ua := r.Header.Get("user-agent")
	bl := r.Header.Get("button-battery-level")

	var cl string
	if bl == "" {
		bl = "100"
		cl = "unknown"
	} else {
		cl = "flic-button"
	}

	stmt, err := prepare("INSERT INTO ping_record VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(newUUID(), id, cl, r.RemoteAddr, ua, bl, now())
	if err != nil {
		return err
	}

	return nil
}

func ListPing() ([]string, error) {
	rows, err := query("SELECT name FROM ping ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		names = append(names, name)
	}

	return names, nil
}

type PingRecord struct {
	Client       string
	BatteryLevel string
	Address      string
	CreatedAt    string
}

func ListPingRecords(name string) ([]PingRecord, error) {
	rows, err := query("SELECT pr.client, pr.battery_level, pr.address, pr.created_at FROM ping_record pr INNER JOIN ping p ON p.id = pr.ping_id WHERE p.name = ? ORDER BY pr.created_at", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []PingRecord
	for rows.Next() {
		var client, batteryLevel, address, createdAt string
		rows.Scan(&client, &batteryLevel, &address, &createdAt)
		records = append(records, PingRecord{client, batteryLevel, address, createdAt})
	}

	return records, nil
}
