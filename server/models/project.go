package models

import (
	"strings"
	"time"

	"github.com/corecollectives/mist/utils"
	"gorm.io/gorm"
)

type Project struct {
	ID          int64   `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Name        string  `gorm:"not null" json:"name"`
	Description *string `json:"description"`

	TagsString string   `gorm:"column:tags" json:"-"`
	Tags       []string `gorm:"-" json:"tags"`

	OwnerID int64 `gorm:"not null" json:"ownerId"`
	Owner   *User `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`

	ProjectMembers []User `gorm:"many2many:project_members;" json:"projectMembers"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (p *Project) BeforeSave(tx *gorm.DB) (err error) {
	if len(p.Tags) > 0 {
		p.TagsString = strings.Join(p.Tags, ",")
	} else {
		p.TagsString = ""
	}
	return
}

func (p *Project) AfterFind(tx *gorm.DB) (err error) {
	if p.TagsString != "" {
		p.Tags = strings.Split(p.TagsString, ",")
	} else {
		p.Tags = []string{}
	}
	return
}

func (p *Project) ToJSON() map[string]interface{} {
	tags := p.Tags
	if tags == nil {
		tags = []string{}
	}

	desc := ""
	if p.Description != nil {
		desc = *p.Description
	}

	return map[string]interface{}{
		"id":             p.ID,
		"name":           p.Name,
		"description":    desc,
		"tags":           tags,
		"ownerId":        p.OwnerID,
		"owner":          p.Owner,
		"projectMembers": p.ProjectMembers,
		"createdAt":      p.CreatedAt,
		"updatedAt":      p.UpdatedAt,
	}
}

func (p *Project) InsertInDB() error {
	p.ID = utils.GenerateRandomId()

	if err := db.Create(p).Error; err != nil {
		return err
	}

	owner := User{ID: p.OwnerID}
	if err := db.Model(p).Association("ProjectMembers").Append(&owner); err != nil {
		return err
	}

	p.Owner = &owner
	p.ProjectMembers = []User{owner}

	return nil
}

func GetProjectByID(projectID int64) (*Project, error) {
	var p Project
	err := db.Preload("Owner").
		Preload("ProjectMembers").
		First(&p, projectID).Error

	if err != nil {
		return nil, err
	}
	return &p, nil
}

func DeleteProjectByID(projectID int64) error {
	return db.Delete(&Project{}, projectID).Error
}

func UpdateProject(p *Project) error {
	return db.Model(p).Updates(map[string]interface{}{
		"name":        p.Name,
		"description": p.Description,
		"tags":        strings.Join(p.Tags, ","),
		"updated_at":  time.Now(),
	}).Error
}

func GetProjectsUserIsPartOf(userID int64) ([]Project, error) {
	var projects []Project

	err := db.Preload("Owner").
		Joins("JOIN project_members pm ON pm.project_id = projects.id").
		Where("pm.user_id = ?", userID).
		Find(&projects).Error

	return projects, err
}

func HasUserAccessToProject(userID, projectID int64) (bool, error) {
	var count int64
	err := db.Table("project_members").
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Count(&count).Error

	return count > 0, err
}

func IsUserProjectOwner(userID, projectID int64) (bool, error) {
	var count int64
	err := db.Model(&Project{}).
		Where("id = ? AND owner_id = ?", projectID, userID).
		Count(&count).Error
	return count > 0, err
}

func UpdateProjectMembers(projectID int64, userIDs []int64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var project Project
		if err := tx.Select("owner_id").First(&project, projectID).Error; err != nil {
			return err
		}

		ownerIncluded := false
		for _, uid := range userIDs {
			if uid == project.OwnerID {
				ownerIncluded = true
				break
			}
		}
		if !ownerIncluded {
			userIDs = append(userIDs, project.OwnerID)
		}

		var users []User
		for _, uid := range userIDs {
			users = append(users, User{ID: uid})
		}

		if err := tx.Model(&Project{ID: projectID}).Association("ProjectMembers").Replace(users); err != nil {
			return err
		}

		return nil
	})
}

//############################################################################################################
//ARCHIVED CODE BELOW

// package models

// import (
// 	"database/sql"
// 	"strings"
// 	"time"

// 	"github.com/corecollectives/mist/utils"
// )

// type Project struct {
// 	ID             int64            `json:"id"`
// 	Name           string           `json:"name"`
// 	Description    sql.NullString   `json:"description"`
// 	Tags           []sql.NullString `json:"tags"`
// 	OwnerID        int64            `json:"ownerId"`
// 	Owner          *User            `json:"owner,omitempty"`
// 	ProjectMembers []User           `json:"projectMembers"`
// 	CreatedAt      time.Time        `json:"createdAt"`
// 	UpdatedAt      time.Time        `json:"updatedAt"`
// }

// func (p *Project) ToJSON() map[string]interface{} {
// 	tags := []string{}
// 	for _, tag := range p.Tags {
// 		if tag.Valid {
// 			tags = append(tags, tag.String)
// 		}
// 	}
// 	if len(tags) == 0 {
// 		tags = []string{}
// 	}

// 	return map[string]interface{}{
// 		"id":             p.ID,
// 		"name":           p.Name,
// 		"description":    p.Description.String,
// 		"tags":           tags,
// 		"ownerId":        p.OwnerID,
// 		"owner":          p.Owner,
// 		"projectMembers": p.ProjectMembers,
// 		"createdAt":      p.CreatedAt,
// 		"updatedAt":      p.UpdatedAt,
// 	}
// }

// func (p *Project) InsertInDB() error {
// 	p.ID = utils.GenerateRandomId()

// 	tagsStr := ""
// 	if len(p.Tags) > 0 {
// 		for i, tag := range p.Tags {
// 			if i > 0 {
// 				tagsStr += ","
// 			}
// 			tagsStr += tag.String
// 		}
// 	}

// 	var desc interface{}
// 	if p.Description.Valid {
// 		desc = p.Description.String
// 	} else {
// 		desc = nil
// 	}

// 	query := `
// 		INSERT INTO projects (id, name, description, tags, owner_id)
// 		VALUES ($1, $2, $3, $4, $5)
// 		RETURNING id, created_at, updated_at
// 	`
// 	err := db.QueryRow(query, p.ID, p.Name, desc, tagsStr, p.OwnerID).Scan(
// 		&p.ID, &p.CreatedAt, &p.UpdatedAt,
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = db.Exec(`
// 		INSERT INTO project_members (project_id, user_id)
// 		VALUES ($1, $2)
// 		ON CONFLICT DO NOTHING
// 	`, p.ID, p.OwnerID)
// 	if err != nil {
// 		return err
// 	}

// 	user, err := GetUserByID(p.OwnerID)
// 	if err != nil {
// 		return err
// 	}
// 	p.Owner = user

// 	p.ProjectMembers = []User{*user}

// 	return nil
// }

// func GetProjectByID(projectID int64) (*Project, error) {
// 	query := `
// 	SELECT
// 		p.id, p.name, p.description, p.tags, p.owner_id, p.created_at, p.updated_at,
// 		u.id, u.username, u.email, u.role, u.avatar_url, u.created_at
// 	FROM projects p
// 	JOIN users u ON p.owner_id = u.id
// 	WHERE p.id = $1
// 	`

// 	project := &Project{}
// 	owner := &User{}
// 	var tagsStr sql.NullString

// 	err := db.QueryRow(query, projectID).Scan(
// 		&project.ID,
// 		&project.Name,
// 		&project.Description,
// 		&tagsStr,
// 		&project.OwnerID,
// 		&project.CreatedAt,
// 		&project.UpdatedAt,
// 		&owner.ID,
// 		&owner.Username,
// 		&owner.Email,
// 		&owner.Role,
// 		&owner.AvatarURL,
// 		&owner.CreatedAt,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if tagsStr.Valid {
// 		strTags := strings.Split(tagsStr.String, ",")
// 		tags := make([]sql.NullString, len(strTags))
// 		for i, s := range strTags {
// 			tags[i] = sql.NullString{
// 				String: s,
// 				Valid:  true,
// 			}
// 		}
// 		project.Tags = tags
// 	}

// 	project.Owner = owner

// 	memberQuery := `
// 	SELECT u.id, u.username, u.email, u.role, u.avatar_url, u.created_at
// 	FROM users u
// 	JOIN project_members pm ON u.id = pm.user_id
// 	WHERE pm.project_id = $1
// 	`
// 	rows, err := db.Query(memberQuery, projectID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var members []User
// 	for rows.Next() {
// 		var member User
// 		if err := rows.Scan(
// 			&member.ID,
// 			&member.Username,
// 			&member.Email,
// 			&member.Role,
// 			&member.AvatarURL,
// 			&member.CreatedAt,
// 		); err != nil {
// 			return nil, err
// 		}
// 		members = append(members, member)
// 	}
// 	project.ProjectMembers = members

// 	return project, nil
// }

// func DeleteProjectByID(projectID int64) error {
// 	query := `DELETE FROM projects WHERE id = $1`
// 	_, err := db.Exec(query, projectID)
// 	return err
// }

// func UpdateProject(p *Project) error {
// 	query := `
// 		UPDATE projects
// 		SET name = $1, description = $2, tags = $3, updated_at = CURRENT_TIMESTAMP
// 		WHERE id = $4
// 		RETURNING updated_at
// 	`
// 	tagsStr := ""
// 	if len(p.Tags) > 0 {
// 		for i, tag := range p.Tags {
// 			if i > 0 {
// 				tagsStr += ","
// 			}
// 			tagsStr += tag.String
// 		}
// 	}

// 	return db.QueryRow(query, p.Name, p.Description, tagsStr, p.ID).Scan(&p.UpdatedAt)
// }

// func GetProjectsUserIsPartOf(userID int64) ([]Project, error) {
// 	query := `
// 	SELECT
// 		p.id, p.name, p.description, p.tags, p.owner_id, p.created_at, p.updated_at,
// 		u.id, u.username, u.email, u.role, u.avatar_url, u.created_at
// 	FROM projects p
// 	JOIN project_members pm ON p.id = pm.project_id
// 	JOIN users u ON p.owner_id = u.id
// 	WHERE pm.user_id = $1
// 	`

// 	rows, err := db.Query(query, userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var projects []Project
// 	for rows.Next() {
// 		var project Project
// 		var owner User
// 		var tagsStr sql.NullString

// 		err := rows.Scan(
// 			&project.ID,
// 			&project.Name,
// 			&project.Description,
// 			&tagsStr,
// 			&project.OwnerID,
// 			&project.CreatedAt,
// 			&project.UpdatedAt,
// 			&owner.ID,
// 			&owner.Username,
// 			&owner.Email,
// 			&owner.Role,
// 			&owner.AvatarURL,
// 			&owner.CreatedAt,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if tagsStr.Valid {
// 			strTags := strings.Split(tagsStr.String, ",")
// 			tags := make([]sql.NullString, len(strTags))
// 			for i, s := range strTags {
// 				tags[i] = sql.NullString{
// 					String: s,
// 					Valid:  true,
// 				}
// 			}
// 			project.Tags = tags
// 		}

// 		project.Owner = &owner
// 		projects = append(projects, project)
// 	}

// 	return projects, nil
// }

// func HasUserAccessToProject(userID, projectID int64) (bool, error) {
// 	query := `
// 		SELECT COUNT(1)
// 		FROM project_members
// 		WHERE project_id = $1 AND user_id = $2
// 	`
// 	var count int
// 	err := db.QueryRow(query, projectID, userID).Scan(&count)
// 	if err != nil {
// 		return false, err
// 	}
// 	return count > 0, nil
// }

// func IsUserProjectOwner(userID, projectID int64) (bool, error) {
// 	query := `SELECT owner_id FROM projects WHERE id = $1`
// 	var ownerID int64
// 	err := db.QueryRow(query, projectID).Scan(&ownerID)
// 	if err != nil {
// 		return false, err
// 	}
// 	return ownerID == userID, nil
// }

// func UpdateProjectMembers(projectID int64, userIDs []int64) error {
// 	query := `SELECT owner_id FROM projects WHERE id = $1`
// 	var ownerID int64
// 	err := db.QueryRow(query, projectID).Scan(&ownerID)
// 	if err != nil {
// 		return err
// 	}

// 	tx, err := db.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()

// 	_, err = tx.Exec(`DELETE FROM project_members WHERE project_id = $1 AND user_id != $2`, projectID, ownerID)
// 	if err != nil {
// 		return err
// 	}

// 	ownerIncluded := false
// 	for _, userID := range userIDs {
// 		if userID == ownerID {
// 			ownerIncluded = true
// 			break
// 		}
// 	}
// 	if !ownerIncluded {
// 		userIDs = append(userIDs, ownerID)
// 	}

// 	for _, userID := range userIDs {
// 		_, err = tx.Exec(`
// 			INSERT INTO project_members (project_id, user_id)
// 			VALUES ($1, $2)
// 			ON CONFLICT (project_id, user_id) DO NOTHING
// 		`, projectID, userID)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return tx.Commit()
// }
