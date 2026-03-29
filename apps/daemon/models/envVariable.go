package models

import (
	"time"

	"github.com/corecollectives/mist/utils"
)

type EnvVariable struct {
	ID        int64     `gorm:"primaryKey;autoIncrement:false" json:"id"`
	AppID     int64     `gorm:"uniqueIndex:idx_app_key;index;not null;constraint:OnDelete:CASCADE" json:"appId"`
	Key       string    `gorm:"uniqueIndex:idx_app_key;not null" json:"key"`
	Value     string    `gorm:"not null" json:"value"`
	IsSecret  bool      `gorm:"default:false" json:"isSecret,omitempty"`
	Runtime   *bool     `gorm:"default:false" json:"runtime,omitempty"`
	Buildtime *bool     `gorm:"default:false" json:"buildtime,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (EnvVariable) TableName() string {
	return "envs"
}

func (E *EnvVariable) IsRuntime() bool {
	if E.Runtime == nil {
		if E.Buildtime == nil {
			return true
		} else {
			return false
		}
	}
	return true
}

func (E *EnvVariable) IsBuildtime() bool {
	if E.Buildtime == nil {
		return false
	} else {
		if *E.Buildtime == true {
			return true
		}
	}
	return false
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

func CreateEnvVariableWithType(appID int64, key, value string, runtime, buildtime *bool) (*EnvVariable, error) {
	env := &EnvVariable{
		ID:    utils.GenerateRandomId(),
		AppID: appID,
		Key:   key,
		Value: value,
	}

	if runtime == nil && buildtime == nil {
		defaultRuntime := true
		runtime = &defaultRuntime
	}

	if runtime != nil {
		env.Runtime = runtime
	}
	if buildtime != nil {
		env.Buildtime = buildtime
	}

	result := db.Create(env)
	if result.Error != nil {
		return nil, result.Error
	}
	return env, nil
}

func UpdateEnvVariableWithType(id int64, key, value string, runtime, buildtime *bool) error {
	updates := map[string]interface{}{
		"key":   key,
		"value": value,
	}

	if runtime != nil {
		updates["runtime"] = *runtime
	}
	if buildtime != nil {
		updates["buildtime"] = *buildtime
	}

	return db.Model(&EnvVariable{ID: id}).Updates(updates).Error
}
