package models

import (
	"time"

	"gorm.io/gorm"
)

type Volume struct {
	ID            int64     `gorm:"primaryKey;autoIncrement:true" json:"id"`
	AppID         int64     `gorm:"uniqueIndex:idx_app_vol_name;not null;constraint:OnDelete:CASCADE" json:"appId"`
	Name          string    `gorm:"uniqueIndex:idx_app_vol_name;not null" json:"name"`
	HostPath      string    `gorm:"not null" json:"hostPath"`
	ContainerPath string    `gorm:"not null" json:"containerPath"`
	ReadOnly      bool      `gorm:"default:false" json:"readOnly"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (v *Volume) ToJson() map[string]interface{} {
	return map[string]interface{}{
		"id":            v.ID,
		"appId":         v.AppID,
		"name":          v.Name,
		"hostPath":      v.HostPath,
		"containerPath": v.ContainerPath,
		"readOnly":      v.ReadOnly,
		"createdAt":     v.CreatedAt,
	}
}

func GetVolumesByAppID(appID int64) ([]Volume, error) {
	var volumes []Volume
	err := db.Where("app_id = ?", appID).Order("created_at DESC").Find(&volumes).Error
	return volumes, err
}

func GetVolumeByID(id int64) (*Volume, error) {
	var vol Volume
	err := db.First(&vol, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &vol, nil
}

func CreateVolume(appID int64, name, hostPath, containerPath string, readOnly bool) (*Volume, error) {
	vol := &Volume{
		AppID:         appID,
		Name:          name,
		HostPath:      hostPath,
		ContainerPath: containerPath,
		ReadOnly:      readOnly,
	}
	err := db.Create(vol).Error
	if err != nil {
		return nil, err
	}
	return vol, nil
}

func UpdateVolume(id int64, name, hostPath, containerPath string, readOnly bool) error {
	return db.Model(&Volume{ID: id}).Updates(map[string]interface{}{
		"name":           name,
		"host_path":      hostPath,
		"container_path": containerPath,
		"read_only":      readOnly,
	}).Error
}

func DeleteVolume(id int64) error {
	return db.Delete(&Volume{}, id).Error
}

//##############################################################################################################
//ARCHIVED CODE BELOW

// package models

// import (
// 	"database/sql"
// 	"time"
// )

// type Volume struct {
// 	ID            int64     `db:"id" json:"id"`
// 	AppID         int64     `db:"app_id" json:"appId"`
// 	Name          string    `db:"name" json:"name"`
// 	HostPath      string    `db:"host_path" json:"hostPath"`
// 	ContainerPath string    `db:"container_path" json:"containerPath"`
// 	ReadOnly      bool      `db:"read_only" json:"readOnly"`
// 	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
// }

// func (v *Volume) ToJson() map[string]interface{} {
// 	return map[string]interface{}{
// 		"id":            v.ID,
// 		"appId":         v.AppID,
// 		"name":          v.Name,
// 		"hostPath":      v.HostPath,
// 		"containerPath": v.ContainerPath,
// 		"readOnly":      v.ReadOnly,
// 		"createdAt":     v.CreatedAt,
// 	}
// }

// func GetVolumesByAppID(appID int64) ([]Volume, error) {
// 	query := `
// 		SELECT id, app_id, name, host_path, container_path, read_only, created_at
// 		FROM volumes
// 		WHERE app_id = ?
// 		ORDER BY created_at DESC
// 	`
// 	rows, err := db.Query(query, appID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var volumes []Volume
// 	for rows.Next() {
// 		var vol Volume
// 		err := rows.Scan(&vol.ID, &vol.AppID, &vol.Name, &vol.HostPath, &vol.ContainerPath, &vol.ReadOnly, &vol.CreatedAt)
// 		if err != nil {
// 			return nil, err
// 		}
// 		volumes = append(volumes, vol)
// 	}

// 	return volumes, rows.Err()
// }

// func GetVolumeByID(id int64) (*Volume, error) {
// 	query := `
// 		SELECT id, app_id, name, host_path, container_path, read_only, created_at
// 		FROM volumes
// 		WHERE id = ?
// 	`
// 	var vol Volume
// 	err := db.QueryRow(query, id).Scan(&vol.ID, &vol.AppID, &vol.Name, &vol.HostPath, &vol.ContainerPath, &vol.ReadOnly, &vol.CreatedAt)
// 	if err == sql.ErrNoRows {
// 		return nil, nil
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &vol, nil
// }

// func CreateVolume(appID int64, name, hostPath, containerPath string, readOnly bool) (*Volume, error) {
// 	query := `
// 		INSERT INTO volumes (app_id, name, host_path, container_path, read_only, created_at)
// 		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
// 	`
// 	result, err := db.Exec(query, appID, name, hostPath, containerPath, readOnly)
// 	if err != nil {
// 		return nil, err
// 	}

// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return GetVolumeByID(id)
// }

// func UpdateVolume(id int64, name, hostPath, containerPath string, readOnly bool) error {
// 	query := `
// 		UPDATE volumes
// 		SET name = ?, host_path = ?, container_path = ?, read_only = ?
// 		WHERE id = ?
// 	`
// 	_, err := db.Exec(query, name, hostPath, containerPath, readOnly, id)
// 	return err
// }

// func DeleteVolume(id int64) error {
// 	query := `DELETE FROM volumes WHERE id = ?`
// 	_, err := db.Exec(query, id)
// 	return err
// }
