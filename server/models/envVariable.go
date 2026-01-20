package models

import (
	"time"

	"github.com/corecollectives/mist/utils"
)

type EnvVariable struct {
	ID int64 `gorm:"primaryKey;autoIncrement:false" json:"id"`

	AppID int64 `gorm:"uniqueIndex:idx_app_key;index;not null;constraint:OnDelete:CASCADE" json:"appId"`

	Key string `gorm:"uniqueIndex:idx_app_key;not null" json:"key"`

	Value string `gorm:"not null" json:"value"`

	IsSecret bool `gorm:"default:false" json:"isSecret,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (EnvVariable) TableName() string {
	return "envs"
}

func CreateEnvVariable(appID int64, key, value string) (*EnvVariable, error) {
	env := &EnvVariable{
		ID:    utils.GenerateRandomId(),
		AppID: appID,
		Key:   key,
		Value: value,
	}

	result := db.Create(env)
	if result.Error != nil {
		return nil, result.Error
	}
	return env, nil
}

func GetEnvVariablesByAppID(appID int64) ([]EnvVariable, error) {
	var envs []EnvVariable
	result := db.Where("app_id = ?", appID).Order("key ASC").Find(&envs)
	return envs, result.Error
}

func UpdateEnvVariable(id int64, key, value string) error {
	updates := map[string]interface{}{
		"key":   key,
		"value": value,
	}
	return db.Model(&EnvVariable{ID: id}).Updates(updates).Error
}

func DeleteEnvVariable(id int64) error {
	return db.Delete(&EnvVariable{}, id).Error
}

func GetEnvVariableByID(id int64) (*EnvVariable, error) {
	var env EnvVariable
	result := db.First(&env, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &env, nil
}

//##########################################################################################################
//ARCHIVED CODE BELOW

// package models

// import (
// 	"time"

// 	"github.com/corecollectives/mist/utils"
// )

// type EnvVariable struct {
// 	ID        int64     `json:"id"`
// 	AppID     int64     `json:"appId"`
// 	Key       string    `json:"key"`
// 	Value     string    `json:"value"`
// 	CreatedAt time.Time `json:"createdAt"`
// 	UpdatedAt time.Time `json:"updatedAt"`
// }

// func CreateEnvVariable(appID int64, key, value string) (*EnvVariable, error) {
// 	id := utils.GenerateRandomId()
// 	query := `
// 		INSERT INTO envs (id, app_id, key, value, created_at, updated_at)
// 		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
// 		RETURNING id, app_id, key, value, created_at, updated_at
// 	`
// 	var env EnvVariable
// 	err := db.QueryRow(query, id, appID, key, value).Scan(
// 		&env.ID, &env.AppID, &env.Key, &env.Value, &env.CreatedAt, &env.UpdatedAt,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &env, nil
// }

// func GetEnvVariablesByAppID(appID int64) ([]EnvVariable, error) {
// 	query := `
// 		SELECT id, app_id, key, value, created_at, updated_at
// 		FROM envs
// 		WHERE app_id = ?
// 		ORDER BY key ASC
// 	`
// 	rows, err := db.Query(query, appID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var envs []EnvVariable
// 	for rows.Next() {
// 		var env EnvVariable
// 		err := rows.Scan(&env.ID, &env.AppID, &env.Key, &env.Value, &env.CreatedAt, &env.UpdatedAt)
// 		if err != nil {
// 			return nil, err
// 		}
// 		envs = append(envs, env)
// 	}
// 	return envs, nil
// }

// func UpdateEnvVariable(id int64, key, value string) error {
// 	query := `
// 		UPDATE envs
// 		SET key = ?, value = ?, updated_at = CURRENT_TIMESTAMP
// 		WHERE id = ?
// 	`
// 	_, err := db.Exec(query, key, value, id)
// 	return err
// }

// func DeleteEnvVariable(id int64) error {
// 	query := `DELETE FROM envs WHERE id = ?`
// 	_, err := db.Exec(query, id)
// 	return err
// }

// func GetEnvVariableByID(id int64) (*EnvVariable, error) {
// 	query := `
// 		SELECT id, app_id, key, value, created_at, updated_at
// 		FROM envs
// 		WHERE id = ?
// 	`
// 	var env EnvVariable
// 	err := db.QueryRow(query, id).Scan(&env.ID, &env.AppID, &env.Key, &env.Value, &env.CreatedAt, &env.UpdatedAt)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &env, nil
// }
