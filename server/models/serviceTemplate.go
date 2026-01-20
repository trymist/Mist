package models

import (
	"time"

	"gorm.io/gorm"
)

type ServiceTemplateCategory string

const (
	CategoryDatabase ServiceTemplateCategory = "database"
	CategoryCache    ServiceTemplateCategory = "cache"
	CategoryQueue    ServiceTemplateCategory = "queue"
	CategoryStorage  ServiceTemplateCategory = "storage"
	CategoryOther    ServiceTemplateCategory = "other"
)

type ServiceTemplate struct {
	ID int64 `gorm:"primaryKey;autoIncrement:true" json:"id"`

	Name string `gorm:"uniqueIndex;not null" json:"name"`

	DisplayName string `gorm:"not null" json:"displayName"`

	Category ServiceTemplateCategory `gorm:"default:'database';index" json:"category"`

	Description *string `json:"description,omitempty"`
	IconURL     *string `json:"iconUrl,omitempty"`

	DockerImage        string  `gorm:"not null" json:"dockerImage"`
	DockerImageVersion *string `json:"dockerImageVersion,omitempty"`

	DefaultPort     int     `gorm:"not null" json:"defaultPort"`
	DefaultEnvVars  *string `json:"defaultEnvVars,omitempty"`
	RequiredEnvVars *string `json:"requiredEnvVars,omitempty"`

	DefaultVolumePath *string `json:"defaultVolumePath,omitempty"`

	VolumeRequired bool `gorm:"default:true" json:"volumeRequired"`

	RecommendedCPU    *float64 `json:"recommendedCpu,omitempty"`
	RecommendedMemory *int     `json:"recommendedMemory,omitempty"`
	MinMemory         *int     `json:"minMemory,omitempty"`

	HealthcheckCommand *string `json:"healthcheckCommand,omitempty"`

	HealthcheckInterval int `gorm:"default:30" json:"healthcheckInterval"`

	AdminUIImage      *string `json:"adminUiImage,omitempty"`
	AdminUIPort       *int    `json:"adminUiPort,omitempty"`
	SetupInstructions *string `json:"setupInstructions,omitempty"`

	IsActive bool `gorm:"default:true;index" json:"isActive"`

	IsFeatured bool `gorm:"default:false" json:"isFeatured"`

	SortOrder int `gorm:"default:0;index" json:"sortOrder"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (st *ServiceTemplate) ToJson() map[string]interface{} {
	return map[string]interface{}{
		"id":                  st.ID,
		"name":                st.Name,
		"displayName":         st.DisplayName,
		"category":            st.Category,
		"description":         st.Description,
		"iconUrl":             st.IconURL,
		"dockerImage":         st.DockerImage,
		"dockerImageVersion":  st.DockerImageVersion,
		"defaultPort":         st.DefaultPort,
		"defaultEnvVars":      st.DefaultEnvVars,
		"requiredEnvVars":     st.RequiredEnvVars,
		"defaultVolumePath":   st.DefaultVolumePath,
		"volumeRequired":      st.VolumeRequired,
		"recommendedCpu":      st.RecommendedCPU,
		"recommendedMemory":   st.RecommendedMemory,
		"minMemory":           st.MinMemory,
		"healthcheckCommand":  st.HealthcheckCommand,
		"healthcheckInterval": st.HealthcheckInterval,
		"adminUiImage":        st.AdminUIImage,
		"adminUiPort":         st.AdminUIPort,
		"setupInstructions":   st.SetupInstructions,
		"isActive":            st.IsActive,
		"isFeatured":          st.IsFeatured,
		"sortOrder":           st.SortOrder,
		"createdAt":           st.CreatedAt,
		"updatedAt":           st.UpdatedAt,
	}
}

func GetAllServiceTemplates() ([]ServiceTemplate, error) {
	var templates []ServiceTemplate
	err := db.Where("is_active = ?", true).
		Order("sort_order, display_name").
		Find(&templates).Error
	return templates, err
}

func GetServiceTemplateByName(name string) (*ServiceTemplate, error) {
	var template ServiceTemplate
	err := db.Where("name = ? AND is_active = ?", name, true).
		First(&template).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &template, nil
}

func GetServiceTemplatesByCategory(category ServiceTemplateCategory) ([]ServiceTemplate, error) {
	var templates []ServiceTemplate
	err := db.Where("category = ? AND is_active = ?", category, true).
		Order("sort_order, display_name").
		Find(&templates).Error
	return templates, err
}

//###############################################################################################################
//ARCHIVED CODE BELOW

// package models

// import (
// 	"database/sql"
// 	"time"
// )

// type ServiceTemplateCategory string

// const (
// 	CategoryDatabase ServiceTemplateCategory = "database"
// 	CategoryCache    ServiceTemplateCategory = "cache"
// 	CategoryQueue    ServiceTemplateCategory = "queue"
// 	CategoryStorage  ServiceTemplateCategory = "storage"
// 	CategoryOther    ServiceTemplateCategory = "other"
// )

// type ServiceTemplate struct {
// 	ID                  int64                   `db:"id" json:"id"`
// 	Name                string                  `db:"name" json:"name"`
// 	DisplayName         string                  `db:"display_name" json:"displayName"`
// 	Category            ServiceTemplateCategory `db:"category" json:"category"`
// 	Description         *string                 `db:"description" json:"description,omitempty"`
// 	IconURL             *string                 `db:"icon_url" json:"iconUrl,omitempty"`
// 	DockerImage         string                  `db:"docker_image" json:"dockerImage"`
// 	DockerImageVersion  *string                 `db:"docker_image_version" json:"dockerImageVersion,omitempty"`
// 	DefaultPort         int                     `db:"default_port" json:"defaultPort"`
// 	DefaultEnvVars      *string                 `db:"default_env_vars" json:"defaultEnvVars,omitempty"`
// 	RequiredEnvVars     *string                 `db:"required_env_vars" json:"requiredEnvVars,omitempty"`
// 	DefaultVolumePath   *string                 `db:"default_volume_path" json:"defaultVolumePath,omitempty"`
// 	VolumeRequired      bool                    `db:"volume_required" json:"volumeRequired"`
// 	RecommendedCPU      *float64                `db:"recommended_cpu" json:"recommendedCpu,omitempty"`
// 	RecommendedMemory   *int                    `db:"recommended_memory" json:"recommendedMemory,omitempty"`
// 	MinMemory           *int                    `db:"min_memory" json:"minMemory,omitempty"`
// 	HealthcheckCommand  *string                 `db:"healthcheck_command" json:"healthcheckCommand,omitempty"`
// 	HealthcheckInterval int                     `db:"healthcheck_interval" json:"healthcheckInterval"`
// 	AdminUIImage        *string                 `db:"admin_ui_image" json:"adminUiImage,omitempty"`
// 	AdminUIPort         *int                    `db:"admin_ui_port" json:"adminUiPort,omitempty"`
// 	SetupInstructions   *string                 `db:"setup_instructions" json:"setupInstructions,omitempty"`
// 	IsActive            bool                    `db:"is_active" json:"isActive"`
// 	IsFeatured          bool                    `db:"is_featured" json:"isFeatured"`
// 	SortOrder           int                     `db:"sort_order" json:"sortOrder"`
// 	CreatedAt           time.Time               `db:"created_at" json:"createdAt"`
// 	UpdatedAt           time.Time               `db:"updated_at" json:"updatedAt"`
// }

// func (st *ServiceTemplate) ToJson() map[string]interface{} {
// 	return map[string]interface{}{
// 		"id":                  st.ID,
// 		"name":                st.Name,
// 		"displayName":         st.DisplayName,
// 		"category":            st.Category,
// 		"description":         st.Description,
// 		"iconUrl":             st.IconURL,
// 		"dockerImage":         st.DockerImage,
// 		"dockerImageVersion":  st.DockerImageVersion,
// 		"defaultPort":         st.DefaultPort,
// 		"defaultEnvVars":      st.DefaultEnvVars,
// 		"requiredEnvVars":     st.RequiredEnvVars,
// 		"defaultVolumePath":   st.DefaultVolumePath,
// 		"volumeRequired":      st.VolumeRequired,
// 		"recommendedCpu":      st.RecommendedCPU,
// 		"recommendedMemory":   st.RecommendedMemory,
// 		"minMemory":           st.MinMemory,
// 		"healthcheckCommand":  st.HealthcheckCommand,
// 		"healthcheckInterval": st.HealthcheckInterval,
// 		"adminUiImage":        st.AdminUIImage,
// 		"adminUiPort":         st.AdminUIPort,
// 		"setupInstructions":   st.SetupInstructions,
// 		"isActive":            st.IsActive,
// 		"isFeatured":          st.IsFeatured,
// 		"sortOrder":           st.SortOrder,
// 		"createdAt":           st.CreatedAt,
// 		"updatedAt":           st.UpdatedAt,
// 	}
// }

// func GetAllServiceTemplates() ([]ServiceTemplate, error) {
// 	var templates []ServiceTemplate
// 	query := `
// 	SELECT id, name, display_name, category, description, icon_url,
// 	       docker_image, docker_image_version, default_port, default_env_vars,
// 	       required_env_vars, default_volume_path, volume_required,
// 	       recommended_cpu, recommended_memory, min_memory,
// 	       healthcheck_command, healthcheck_interval,
// 	       admin_ui_image, admin_ui_port, setup_instructions,
// 	       is_active, is_featured, sort_order, created_at, updated_at
// 	FROM service_templates
// 	WHERE is_active = 1
// 	ORDER BY sort_order, display_name
// 	`
// 	rows, err := db.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var template ServiceTemplate
// 		err := rows.Scan(
// 			&template.ID, &template.Name, &template.DisplayName, &template.Category,
// 			&template.Description, &template.IconURL, &template.DockerImage,
// 			&template.DockerImageVersion, &template.DefaultPort, &template.DefaultEnvVars,
// 			&template.RequiredEnvVars, &template.DefaultVolumePath, &template.VolumeRequired,
// 			&template.RecommendedCPU, &template.RecommendedMemory, &template.MinMemory,
// 			&template.HealthcheckCommand, &template.HealthcheckInterval,
// 			&template.AdminUIImage, &template.AdminUIPort, &template.SetupInstructions,
// 			&template.IsActive, &template.IsFeatured, &template.SortOrder,
// 			&template.CreatedAt, &template.UpdatedAt,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		templates = append(templates, template)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return templates, nil
// }

// func GetServiceTemplateByName(name string) (*ServiceTemplate, error) {
// 	var template ServiceTemplate
// 	query := `
// 	SELECT id, name, display_name, category, description, icon_url,
// 	       docker_image, docker_image_version, default_port, default_env_vars,
// 	       required_env_vars, default_volume_path, volume_required,
// 	       recommended_cpu, recommended_memory, min_memory,
// 	       healthcheck_command, healthcheck_interval,
// 	       admin_ui_image, admin_ui_port, setup_instructions,
// 	       is_active, is_featured, sort_order, created_at, updated_at
// 	FROM service_templates
// 	WHERE name = ? AND is_active = 1
// 	`
// 	err := db.QueryRow(query, name).Scan(
// 		&template.ID, &template.Name, &template.DisplayName, &template.Category,
// 		&template.Description, &template.IconURL, &template.DockerImage,
// 		&template.DockerImageVersion, &template.DefaultPort, &template.DefaultEnvVars,
// 		&template.RequiredEnvVars, &template.DefaultVolumePath, &template.VolumeRequired,
// 		&template.RecommendedCPU, &template.RecommendedMemory, &template.MinMemory,
// 		&template.HealthcheckCommand, &template.HealthcheckInterval,
// 		&template.AdminUIImage, &template.AdminUIPort, &template.SetupInstructions,
// 		&template.IsActive, &template.IsFeatured, &template.SortOrder,
// 		&template.CreatedAt, &template.UpdatedAt,
// 	)
// 	if err == sql.ErrNoRows {
// 		return nil, nil
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &template, nil
// }

// func GetServiceTemplatesByCategory(category ServiceTemplateCategory) ([]ServiceTemplate, error) {
// 	var templates []ServiceTemplate
// 	query := `
// 	SELECT id, name, display_name, category, description, icon_url,
// 	       docker_image, docker_image_version, default_port, default_env_vars,
// 	       required_env_vars, default_volume_path, volume_required,
// 	       recommended_cpu, recommended_memory, min_memory,
// 	       healthcheck_command, healthcheck_interval,
// 	       admin_ui_image, admin_ui_port, setup_instructions,
// 	       is_active, is_featured, sort_order, created_at, updated_at
// 	FROM service_templates
// 	WHERE category = ? AND is_active = 1
// 	ORDER BY sort_order, display_name
// 	`
// 	rows, err := db.Query(query, category)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var template ServiceTemplate
// 		err := rows.Scan(
// 			&template.ID, &template.Name, &template.DisplayName, &template.Category,
// 			&template.Description, &template.IconURL, &template.DockerImage,
// 			&template.DockerImageVersion, &template.DefaultPort, &template.DefaultEnvVars,
// 			&template.RequiredEnvVars, &template.DefaultVolumePath, &template.VolumeRequired,
// 			&template.RecommendedCPU, &template.RecommendedMemory, &template.MinMemory,
// 			&template.HealthcheckCommand, &template.HealthcheckInterval,
// 			&template.AdminUIImage, &template.AdminUIPort, &template.SetupInstructions,
// 			&template.IsActive, &template.IsFeatured, &template.SortOrder,
// 			&template.CreatedAt, &template.UpdatedAt,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		templates = append(templates, template)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return templates, nil
// }
