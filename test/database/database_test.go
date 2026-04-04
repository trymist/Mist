package db

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	mistdb "github.com/corecollectives/mist/db"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	err = mistdb.MigrateDB(db)
	if err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	models.SetDB(db)
	return db
}

func cleanupDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	sqlDB.Close()
}

func TestInitDB_Success(t *testing.T) {
	os.Setenv("ENV", "dev")
	defer os.Unsetenv("ENV")

	db, err := mistdb.InitDB()
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer cleanupDB(db)

	if db == nil {
		t.Fatal("db should not be nil")
	}
}

func TestInitDB_FailsOnInvalidPath(t *testing.T) {
	origEnv := os.Getenv("ENV")
	os.Setenv("ENV", "dev")
	defer os.Setenv("ENV", origEnv)

	t.Skip("Skipping test that depends on filesystem permissions")
}

func TestMigrateDb_CreatesAllTables(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	expectedTables := []string{
		"users", "api_tokens", "apps", "audit_logs", "backups",
		"deployments", "envs", "projects", "project_members",
		"git_providers", "github_installations", "app_repositories", "domains",
		"volumes", "crons", "registries", "sessions", "notifications",
	}

	for _, table := range expectedTables {
		if !db.Migrator().HasTable(table) {
			t.Errorf("table %s should exist", table)
		}
	}
}

func TestMigrateDb_SkipsExistingTables(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	err := mistdb.MigrateDB(db)
	if err != nil {
		t.Errorf("re-running migrations should not fail: %v", err)
	}
}

func TestUser_Create(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	user := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hash",
	}

	err := user.Create()
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if user.ID == 0 {
		t.Error("user ID should be set after creation")
	}

	if user.Role != "user" {
		t.Error("default role should be user")
	}
}

func TestUser_Create_DuplicateUsername(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	user1 := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "uniqueuser",
		Email:        "user1@example.com",
		PasswordHash: "hash",
	}
	user2 := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "uniqueuser",
		Email:        "user2@example.com",
		PasswordHash: "hash",
	}

	if err := user1.Create(); err != nil {
		t.Fatalf("first user creation failed: %v", err)
	}

	err := user2.Create()
	if err == nil {
		t.Error("should fail with duplicate username")
	}
}

func TestUser_Create_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	user1 := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "user1",
		Email:        "same@example.com",
		PasswordHash: "hash",
	}
	user2 := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "user2",
		Email:        "same@example.com",
		PasswordHash: "hash",
	}

	if err := user1.Create(); err != nil {
		t.Fatalf("first user creation failed: %v", err)
	}

	err := user2.Create()
	if err == nil {
		t.Error("should fail with duplicate email")
	}
}

func TestUser_SetPassword(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	user := &models.User{
		ID:       utils.GenerateRandomId(),
		Username: "passuser",
		Email:    "pass@example.com",
	}

	err := user.SetPassword("testpassword123")
	if err != nil {
		t.Fatalf("SetPassword failed: %v", err)
	}

	if user.PasswordHash == "testpassword123" {
		t.Error("password should not be stored in plaintext")
	}

	if len(user.PasswordHash) < 50 {
		t.Error("password hash should be properly hashed")
	}
}

func TestUser_MatchPassword_Success(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	user := &models.User{
		ID:       utils.GenerateRandomId(),
		Username: "matchuser",
		Email:    "match@example.com",
	}
	user.SetPassword("correctpassword")
	user.Create()

	if !user.MatchPassword("correctpassword") {
		t.Error("MatchPassword should return true for correct password")
	}
}

func TestUser_MatchPassword_WrongPassword(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	user := &models.User{
		ID:       utils.GenerateRandomId(),
		Username: "wrongpass",
		Email:    "wrong@example.com",
	}
	user.SetPassword("correctpassword")
	user.Create()

	if user.MatchPassword("wrongpassword") {
		t.Error("MatchPassword should return false for wrong password")
	}
}

func TestGetUserByID_Exists(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	origUser := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "getbyid",
		Email:        "getbyid@example.com",
		PasswordHash: "hash",
		Role:         "admin",
	}
	origUser.Create()

	user, err := models.GetUserByID(origUser.ID)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}

	if user.Username != origUser.Username {
		t.Error("retrieved user should match original")
	}
}

func TestGetUserByID_NotExists(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	_, err := models.GetUserByID(999999)
	if err == nil {
		t.Error("should return error for non-existent user")
	}
}

func TestGetUserByEmail_Exists(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	origUser := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "emailtest",
		Email:        "findme@example.com",
		PasswordHash: "hash",
	}
	origUser.Create()

	user, err := models.GetUserByEmail("findme@example.com")
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}

	if user.ID != origUser.ID {
		t.Error("retrieved user should match by email")
	}
}

func TestGetUserByUsername_Exists(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	origUser := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "uniquetestuser",
		Email:        "unique@example.com",
		PasswordHash: "hash",
	}
	origUser.Create()

	user, err := models.GetUserByUsername("uniquetestuser")
	if err != nil {
		t.Fatalf("GetUserByUsername failed: %v", err)
	}

	if user.ID != origUser.ID {
		t.Error("retrieved user should match by username")
	}
}

func TestUpdateUser(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	user := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "updateme",
		Email:        "update@example.com",
		PasswordHash: "hash",
	}
	user.Create()

	newUsername := "updatedname"
	newEmail := "newemail@example.com"
	user.Username = newUsername
	user.Email = newEmail

	err := models.UpdateUser(user)
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}

	updated, _ := models.GetUserByID(user.ID)
	if updated.Username != newUsername {
		t.Error("username should be updated")
	}
	if updated.Email != newEmail {
		t.Error("email should be updated")
	}
}

func TestUpdateUserPassword(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	user := &models.User{
		ID:       utils.GenerateRandomId(),
		Username: "passupdate",
		Email:    "passupdate@example.com",
	}
	user.SetPassword("oldpassword")
	user.Create()

	user.SetPassword("newpassword")
	err := user.UpdatePassword()
	if err != nil {
		t.Fatalf("UpdatePassword failed: %v", err)
	}

	if user.MatchPassword("oldpassword") {
		t.Error("old password should not match")
	}
	if !user.MatchPassword("newpassword") {
		t.Error("new password should match")
	}
}

func TestDeleteUser(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	user := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "deleteme",
		Email:        "delete@example.com",
		PasswordHash: "hash",
	}
	user.Create()

	err := models.DeleteUserByID(user.ID)
	if err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}

	_, err = models.GetUserByID(user.ID)
	if err == nil {
		t.Error("user should be deleted")
	}
}

func TestGetAllUsers(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	for i := 0; i < 5; i++ {
		user := &models.User{
			ID:           utils.GenerateRandomId(),
			Username:     "listuser" + string(rune('a'+i)),
			Email:        "listuser" + string(rune('a'+i)) + "@example.com",
			PasswordHash: "hash",
		}
		user.Create()
	}

	users, err := models.GetAllUsers()
	if err != nil {
		t.Fatalf("GetAllUsers failed: %v", err)
	}

	if len(users) < 5 {
		t.Error("should return all users")
	}
}

func TestGetUserCount(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	initial, _ := models.GetUserCount()

	for i := 0; i < 3; i++ {
		user := &models.User{
			ID:           utils.GenerateRandomId(),
			Username:     "countuser" + string(rune('a'+i)),
			Email:        "countuser" + string(rune('a'+i)) + "@example.com",
			PasswordHash: "hash",
		}
		user.Create()
	}

	count, err := models.GetUserCount()
	if err != nil {
		t.Fatalf("GetUserCount failed: %v", err)
	}

	if count != initial+3 {
		t.Error("count should include all created users")
	}
}

func TestUserRoleRetrieval(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	adminUser := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "adminuser",
		Email:        "admin@example.com",
		PasswordHash: "hash",
		Role:         "admin",
	}
	adminUser.Create()

	role, err := models.GetUserRole(adminUser.ID)
	if err != nil {
		t.Fatalf("GetUserRole failed: %v", err)
	}

	if role != "admin" {
		t.Error("should return correct role")
	}
}

func TestUser_OptionalFields(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	fullName := "Test Full Name"
	avatarURL := "https://example.com/avatar.png"
	bio := "Test bio"

	user := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "optionalfields",
		Email:        "optional@example.com",
		PasswordHash: "hash",
		FullName:     &fullName,
		AvatarURL:    &avatarURL,
		Bio:          &bio,
	}
	user.Create()

	retrieved, _ := models.GetUserByID(user.ID)
	if retrieved.FullName == nil || *retrieved.FullName != fullName {
		t.Error("FullName should be stored correctly")
	}
	if retrieved.AvatarURL == nil || *retrieved.AvatarURL != avatarURL {
		t.Error("AvatarURL should be stored correctly")
	}
	if retrieved.Bio == nil || *retrieved.Bio != bio {
		t.Error("Bio should be stored correctly")
	}
}

func TestProject_Insert(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "projectowner",
		Email:        "projectowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	desc := "Test project description"
	project := &models.Project{
		ID:          utils.GenerateRandomId(),
		Name:        "Test Project",
		Description: &desc,
		OwnerID:     owner.ID,
		Tags:        []string{"tag1", "tag2"},
	}

	err := project.InsertInDB()
	if err != nil {
		t.Fatalf("InsertInDB failed: %v", err)
	}

	if project.Owner == nil {
		t.Error("owner should be populated after insert")
	}
}

func TestProject_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "getprojectowner",
		Email:        "getprojectowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Get Project Test",
		OwnerID: owner.ID,
		Tags:    []string{"test"},
	}
	project.InsertInDB()

	retrieved, err := models.GetProjectByID(project.ID)
	if err != nil {
		t.Fatalf("GetProjectByID failed: %v", err)
	}

	if retrieved.Name != project.Name {
		t.Error("project name should match")
	}
}

func TestProject_Update(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "updateprojectowner",
		Email:        "updateprojectowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Original Name",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	newName := "Updated Name"
	newTags := []string{"newtag"}
	project.Name = newName
	project.Tags = newTags

	err := models.UpdateProject(project)
	if err != nil {
		t.Fatalf("UpdateProject failed: %v", err)
	}

	updated, _ := models.GetProjectByID(project.ID)
	if updated.Name != newName {
		t.Error("name should be updated")
	}
}

func TestProject_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "deleteprojectowner",
		Email:        "deleteprojectowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Delete Project Test",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	err := models.DeleteProjectByID(project.ID)
	if err != nil {
		t.Fatalf("DeleteProject failed: %v", err)
	}

	_, err = models.GetProjectByID(project.ID)
	if err == nil {
		t.Error("project should be deleted")
	}
}

func TestProject_HasUserAccess(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	user := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "accessuser",
		Email:        "accessuser@example.com",
		PasswordHash: "hash",
	}
	user.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Access Project",
		OwnerID: user.ID,
	}
	project.InsertInDB()

	hasAccess, err := models.HasUserAccessToProject(user.ID, project.ID)
	if err != nil {
		t.Fatalf("HasUserAccessToProject failed: %v", err)
	}

	if !hasAccess {
		t.Error("owner should have access to project")
	}
}

func TestProject_IsUserOwner(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	user := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "owneruser",
		Email:        "owneruser@example.com",
		PasswordHash: "hash",
	}
	user.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Owner Project",
		OwnerID: user.ID,
	}
	project.InsertInDB()

	isOwner, err := models.IsUserProjectOwner(user.ID, project.ID)
	if err != nil {
		t.Fatalf("IsUserProjectOwner failed: %v", err)
	}

	if !isOwner {
		t.Error("user should be project owner")
	}
}

func TestProject_UpdateMembers(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "memberowner",
		Email:        "memberowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	member1 := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "member1",
		Email:        "member1@example.com",
		PasswordHash: "hash",
	}
	member1.Create()

	member2 := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "member2",
		Email:        "member2@example.com",
		PasswordHash: "hash",
	}
	member2.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Members Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	err := models.UpdateProjectMembers(project.ID, []int64{member1.ID, member2.ID})
	if err != nil {
		t.Fatalf("UpdateProjectMembers failed: %v", err)
	}

	retrieved, _ := models.GetProjectByID(project.ID)
	if len(retrieved.ProjectMembers) < 2 {
		t.Error("project should have members")
	}
}

func TestProject_OwnerAlwaysIncludedInMembers(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "forcedowner",
		Email:        "forcedowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Forced Owner Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	err := models.UpdateProjectMembers(project.ID, []int64{})
	if err != nil {
		t.Fatalf("UpdateProjectMembers failed: %v", err)
	}

	retrieved, _ := models.GetProjectByID(project.ID)
	found := false
	for _, member := range retrieved.ProjectMembers {
		if member.ID == owner.ID {
			found = true
			break
		}
	}

	if !found {
		t.Error("owner should always be included in members")
	}
}

func TestProject_GetProjectsUserIsPartOf(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	user := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "partofuser",
		Email:        "partofuser@example.com",
		PasswordHash: "hash",
	}
	user.Create()

	project1 := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Project 1",
		OwnerID: user.ID,
	}
	project1.InsertInDB()

	projects, err := models.GetProjectsUserIsPartOf(user.ID)
	if err != nil {
		t.Fatalf("GetProjectsUserIsPartOf failed: %v", err)
	}

	if len(projects) < 1 {
		t.Error("should return projects user is part of")
	}
}

func TestApp_Insert(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "appowner",
		Email:        "appowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "App Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "Test App",
		CreatedBy: owner.ID,
		AppType:   models.AppTypeWeb,
		GitBranch: "main",
	}

	err := app.InsertInDB()
	if err != nil {
		t.Fatalf("InsertInDB failed: %v", err)
	}

	if app.ID == 0 {
		t.Error("app ID should be set")
	}
	if app.AppType != models.AppTypeWeb {
		t.Error("default app type should be web")
	}
	if app.Status != models.StatusStopped {
		t.Error("default status should be stopped")
	}
}

func TestApp_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "getappowner",
		Email:        "getappowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Get App Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "Get App Test",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	retrieved, err := models.GetApplicationByID(app.ID)
	if err != nil {
		t.Fatalf("GetApplicationByID failed: %v", err)
	}

	if retrieved.Name != app.Name {
		t.Error("retrieved app should match")
	}
}

func TestApp_GetByProjectID(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "listappowner",
		Email:        "listappowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "List App Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	for i := 0; i < 3; i++ {
		app := &models.App{
			ProjectID: project.ID,
			Name:      "App " + string(rune('a'+i)),
			CreatedBy: owner.ID,
		}
		app.InsertInDB()
	}

	apps, err := models.GetApplicationByProjectID(project.ID)
	if err != nil {
		t.Fatalf("GetApplicationByProjectID failed: %v", err)
	}

	if len(apps) < 3 {
		t.Error("should return all apps for project")
	}
}

func TestApp_Update(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "updateappowner",
		Email:        "updateappowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Update App Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "Original App Name",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	newName := "Updated App Name"
	port := int64(8080)
	app.Name = newName
	app.Port = &port

	err := app.UpdateApplication()
	if err != nil {
		t.Fatalf("UpdateApplication failed: %v", err)
	}

	updated, _ := models.GetApplicationByID(app.ID)
	if updated.Name != newName {
		t.Error("name should be updated")
	}
	if updated.Port == nil || *updated.Port != 8080 {
		t.Error("port should be updated")
	}
}

func TestApp_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "deleteappowner",
		Email:        "deleteappowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Delete App Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "Delete App Test",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	err := models.DeleteApplication(app.ID)
	if err != nil {
		t.Fatalf("DeleteApplication failed: %v", err)
	}

	_, err = models.GetApplicationByID(app.ID)
	if err == nil {
		t.Error("app should be deleted")
	}
}

func TestApp_IsUserOwner(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "appownerverify",
		Email:        "appownerverify@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "App Owner Verify Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "App Owner Verify Test",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	isOwner, err := models.IsUserApplicationOwner(owner.ID, app.ID)
	if err != nil {
		t.Fatalf("IsUserApplicationOwner failed: %v", err)
	}

	if !isOwner {
		t.Error("user should be app owner")
	}
}

func TestApp_CascadingDelete(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "cascadeappowner",
		Email:        "cascadeappowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Cascade App Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "Cascade App Test",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	deployment := &models.Deployment{
		ID:         utils.GenerateRandomId(),
		AppID:      app.ID,
		CommitHash: "abc123",
	}
	deployment.CreateDeployment()

	models.DeleteApplication(app.ID)

	var appCount int64
	db.Model(&models.App{}).Where("id = ?", app.ID).Count(&appCount)
	if appCount != 0 {
		t.Error("app should be deleted")
	}
}

func TestDeployment_Create(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "deployowner",
		Email:        "deployowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Deploy Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "Deploy App",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	deployment := &models.Deployment{
		ID:          utils.GenerateRandomId(),
		AppID:       app.ID,
		CommitHash:  "abc123def456",
		TriggeredBy: &owner.ID,
	}

	err := deployment.CreateDeployment()
	if err != nil {
		t.Fatalf("CreateDeployment failed: %v", err)
	}

	if deployment.DeploymentNumber == nil {
		t.Error("deployment number should be set")
	}
	if *deployment.DeploymentNumber != 1 {
		t.Error("first deployment should have number 1")
	}
	if deployment.Status != models.DeploymentStatusPending {
		t.Error("default status should be pending")
	}
}

func TestDeployment_AutoIncrementNumber(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "multideployowner",
		Email:        "multideployowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Multi Deploy Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "Multi Deploy App",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	for i := 0; i < 3; i++ {
		deployment := &models.Deployment{
			ID:         utils.GenerateRandomId(),
			AppID:      app.ID,
			CommitHash: fmt.Sprintf("commit_%d_%d", time.Now().UnixNano(), i),
		}
		err := deployment.CreateDeployment()
		if err != nil {
			t.Fatalf("deployment creation failed: %v", err)
		}

		if deployment.DeploymentNumber == nil {
			t.Errorf("deployment %d: deployment number is nil", i)
		} else {
			t.Logf("deployment %d: got number %d", i, *deployment.DeploymentNumber)
		}
	}

	deployments, _ := models.GetDeploymentsByAppID(app.ID)
	if len(deployments) != 3 {
		t.Errorf("expected 3 deployments, got %d", len(deployments))
	}
}

func TestDeployment_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "getdeployowner",
		Email:        "getdeployowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Get Deploy Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "Get Deploy App",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	deployment := &models.Deployment{
		ID:         utils.GenerateRandomId(),
		AppID:      app.ID,
		CommitHash: "getdeploy123",
	}
	deployment.CreateDeployment()

	retrieved, err := models.GetDeploymentByID(deployment.ID)
	if err != nil {
		t.Fatalf("GetDeploymentByID failed: %v", err)
	}

	if retrieved.CommitHash != deployment.CommitHash {
		t.Error("retrieved deployment should match")
	}
}

func TestDeployment_GetByAppID(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "listdeployowner",
		Email:        "listdeployowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "List Deploy Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "List Deploy App",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	for i := 0; i < 3; i++ {
		deployment := &models.Deployment{
			ID:         utils.GenerateRandomId(),
			AppID:      app.ID,
			CommitHash: "listdeploy" + string(rune('a'+i)),
		}
		deployment.CreateDeployment()
	}

	deployments, err := models.GetDeploymentsByAppID(app.ID)
	if err != nil {
		t.Fatalf("GetDeploymentsByAppID failed: %v", err)
	}

	if len(deployments) < 3 {
		t.Error("should return all deployments for app")
	}
}

func TestDeployment_UpdateStatus(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "statusdeployowner",
		Email:        "statusdeployowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Status Deploy Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "Status Deploy App",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	deployment := &models.Deployment{
		ID:         utils.GenerateRandomId(),
		AppID:      app.ID,
		CommitHash: "statusdeploy123",
	}
	deployment.CreateDeployment()

	err := models.UpdateDeploymentStatus(deployment.ID, string(models.DeploymentStatusSuccess), "completed", 100, nil)
	if err != nil {
		t.Fatalf("UpdateDeploymentStatus failed: %v", err)
	}

	updated, _ := models.GetDeploymentByID(deployment.ID)
	if updated.Status != models.DeploymentStatusSuccess {
		t.Error("status should be updated")
	}
	if updated.FinishedAt == nil {
		t.Error("finished_at should be set on completion")
	}
}

func TestDeployment_MarkActive(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "activedeployowner",
		Email:        "activedeployowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Active Deploy Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "Active Deploy App",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	deployment1 := &models.Deployment{
		ID:         utils.GenerateRandomId(),
		AppID:      app.ID,
		CommitHash: "active1",
	}
	deployment1.CreateDeployment()

	deployment2 := &models.Deployment{
		ID:         utils.GenerateRandomId(),
		AppID:      app.ID,
		CommitHash: "active2",
	}
	deployment2.CreateDeployment()

	err := models.MarkDeploymentActive(deployment1.ID, app.ID)
	if err != nil {
		t.Fatalf("MarkDeploymentActive failed: %v", err)
	}

	d1, _ := models.GetDeploymentByID(deployment1.ID)
	d2, _ := models.GetDeploymentByID(deployment2.ID)

	if !d1.IsActive {
		t.Error("first deployment should be active")
	}
	if d2.IsActive {
		t.Error("second deployment should not be active")
	}
}

func TestDeployment_GetIncomplete(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "incompletedeployowner",
		Email:        "incompletedeployowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Incomplete Deploy Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "Incomplete Deploy App",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	deployment1 := &models.Deployment{
		ID:         utils.GenerateRandomId(),
		AppID:      app.ID,
		CommitHash: "incomplete1",
	}
	deployment1.CreateDeployment()

	deployment2 := &models.Deployment{
		ID:         utils.GenerateRandomId(),
		AppID:      app.ID,
		CommitHash: "incomplete2",
	}
	deployment2.CreateDeployment()

	models.UpdateDeploymentStatus(deployment2.ID, string(models.DeploymentStatusBuilding), "building", 50, nil)

	incomplete, err := models.GetIncompleteDeployments()
	if err != nil {
		t.Fatalf("GetIncompleteDeployments failed: %v", err)
	}

	found := false
	for _, d := range incomplete {
		if d.ID == deployment2.ID && d.AppID == app.ID {
			found = true
			break
		}
	}

	if !found {
		t.Logf("Incomplete deployments found: %d", len(incomplete))
		for _, d := range incomplete {
			t.Logf("  - ID: %d, AppID: %d, Status: %s", d.ID, d.AppID, d.Status)
		}
		t.Error("should return incomplete deployments")
	}
}

func TestDomain_Create(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "domainowner",
		Email:        "domainowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Domain Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "Domain App",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	domain, err := models.CreateDomain(app.ID, "example.com")
	if err != nil {
		t.Fatalf("CreateDomain failed: %v", err)
	}

	if domain.Domain != "example.com" {
		t.Error("domain name should match")
	}
	if domain.AppID != app.ID {
		t.Error("app id should match")
	}
	if domain.SslStatus != models.SSLStatusPending {
		t.Error("default ssl status should be pending")
	}
}

func TestDomain_GetByAppID(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "listdomainowner",
		Email:        "listdomainowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "List Domain Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "List Domain App",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	models.CreateDomain(app.ID, "domain1.example.com")
	models.CreateDomain(app.ID, "domain2.example.com")

	domains, err := models.GetDomainsByAppID(app.ID)
	if err != nil {
		t.Fatalf("GetDomainsByAppID failed: %v", err)
	}

	if len(domains) < 2 {
		t.Error("should return all domains for app")
	}
}

func TestDomain_GetPrimary(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "primarydomainowner",
		Email:        "primarydomainowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Primary Domain Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "Primary Domain App",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	models.CreateDomain(app.ID, "primary.example.com")
	models.CreateDomain(app.ID, "secondary.example.com")

	primary, err := models.GetPrimaryDomainByAppID(app.ID)
	if err != nil {
		t.Fatalf("GetPrimaryDomainByAppID failed: %v", err)
	}

	if primary.Domain != "primary.example.com" {
		t.Error("first created domain should be primary")
	}
}

func TestDomain_Update(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "updatedomainowner",
		Email:        "updatedomainowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Update Domain Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "Update Domain App",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	domain, _ := models.CreateDomain(app.ID, "olddomain.com")

	err := models.UpdateDomain(domain.ID, "newdomain.com")
	if err != nil {
		t.Fatalf("UpdateDomain failed: %v", err)
	}

	updated, _ := models.GetDomainByID(domain.ID)
	if updated.Domain != "newdomain.com" {
		t.Error("domain should be updated")
	}
}

func TestDomain_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "deletedomainowner",
		Email:        "deletedomainowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Delete Domain Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "Delete Domain App",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	domain, _ := models.CreateDomain(app.ID, "deleteme.com")

	err := models.DeleteDomain(domain.ID)
	if err != nil {
		t.Fatalf("DeleteDomain failed: %v", err)
	}

	_, err = models.GetDomainByID(domain.ID)
	if err == nil {
		t.Error("domain should be deleted")
	}
}

func TestDomain_UniqueConstraint(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "uniquedomainowner",
		Email:        "uniquedomainowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Unique Domain Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app1 := &models.App{
		ProjectID: project.ID,
		Name:      "Unique Domain App 1",
		CreatedBy: owner.ID,
	}
	app1.InsertInDB()

	app2 := &models.App{
		ProjectID: project.ID,
		Name:      "Unique Domain App 2",
		CreatedBy: owner.ID,
	}
	app2.InsertInDB()

	_, err := models.CreateDomain(app1.ID, "same-domain.com")
	if err != nil {
		t.Fatalf("first domain creation failed: %v", err)
	}

	_, err = models.CreateDomain(app2.ID, "same-domain.com")
	if err == nil {
		t.Error("should fail with duplicate domain")
	}
}

func TestEnvVariable_Create(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "envarowner",
		Email:        "envarowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "EnvVar Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "EnvVar App",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	env, err := models.CreateEnvVariable(app.ID, "TEST_VAR", "test_value")
	if err != nil {
		t.Fatalf("CreateEnvVariable failed: %v", err)
	}

	if env.Key != "TEST_VAR" {
		t.Error("key should match")
	}
	if env.Value != "test_value" {
		t.Error("value should match")
	}
}

func TestEnvVariable_GetByAppID(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "listenvarowner",
		Email:        "listenvarowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "List EnvVar Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "List EnvVar App",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	models.CreateEnvVariable(app.ID, "VAR1", "val1")
	models.CreateEnvVariable(app.ID, "VAR2", "val2")
	models.CreateEnvVariable(app.ID, "VAR3", "val3")

	envs, err := models.GetEnvVariablesByAppID(app.ID)
	if err != nil {
		t.Fatalf("GetEnvVariablesByAppID failed: %v", err)
	}

	if len(envs) < 3 {
		t.Error("should return all env variables")
	}
}

func TestEnvVariable_Update(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "updateenvarowner",
		Email:        "updateenvarowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Update EnvVar Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "Update EnvVar App",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	env, _ := models.CreateEnvVariable(app.ID, "UPDATEME", "oldvalue")

	err := models.UpdateEnvVariable(env.ID, "UPDATEME", "newvalue")
	if err != nil {
		t.Fatalf("UpdateEnvVariable failed: %v", err)
	}

	updated, _ := models.GetEnvVariableByID(env.ID)
	if updated.Value != "newvalue" {
		t.Error("value should be updated")
	}
}

func TestEnvVariable_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "deleteenvarowner",
		Email:        "deleteenvarowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Delete EnvVar Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "Delete EnvVar App",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	env, _ := models.CreateEnvVariable(app.ID, "DELETEME", "value")

	err := models.DeleteEnvVariable(env.ID)
	if err != nil {
		t.Fatalf("DeleteEnvVariable failed: %v", err)
	}

	_, err = models.GetEnvVariableByID(env.ID)
	if err == nil {
		t.Error("env variable should be deleted")
	}
}

func TestEnvVariable_UniqueConstraint(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "uniqueenvarowner",
		Email:        "uniqueenvarowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Unique EnvVar Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "Unique EnvVar App",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	_, err := models.CreateEnvVariable(app.ID, "UNIQUE_VAR", "value1")
	if err != nil {
		t.Fatalf("first env var creation failed: %v", err)
	}

	_, err = models.CreateEnvVariable(app.ID, "UNIQUE_VAR", "value2")
	if err == nil {
		t.Error("should fail with duplicate key for same app")
	}
}

func TestVolume_Create(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "volumeowner",
		Email:        "volumeowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Volume Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "Volume App",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	vol, err := models.CreateVolume(app.ID, "data", "/host/data", "/container/data", false)
	if err != nil {
		t.Fatalf("CreateVolume failed: %v", err)
	}

	if vol.Name != "data" {
		t.Error("name should match")
	}
	if vol.HostPath != "/host/data" {
		t.Error("host path should match")
	}
	if vol.ContainerPath != "/container/data" {
		t.Error("container path should match")
	}
}

func TestVolume_GetByAppID(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "listvolumeowner",
		Email:        "listvolumeowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "List Volume Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "List Volume App",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	models.CreateVolume(app.ID, "vol1", "/h1", "/c1", false)
	models.CreateVolume(app.ID, "vol2", "/h2", "/c2", true)

	vols, err := models.GetVolumesByAppID(app.ID)
	if err != nil {
		t.Fatalf("GetVolumesByAppID failed: %v", err)
	}

	if len(vols) < 2 {
		t.Error("should return all volumes")
	}
}

func TestVolume_Update(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "updatevolumeowner",
		Email:        "updatevolumeowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Update Volume Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "Update Volume App",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	vol, _ := models.CreateVolume(app.ID, "oldname", "/old/host", "/old/container", false)

	err := models.UpdateVolume(vol.ID, "newname", "/new/host", "/new/container", true)
	if err != nil {
		t.Fatalf("UpdateVolume failed: %v", err)
	}

	updated, _ := models.GetVolumeByID(vol.ID)
	if updated.Name != "newname" {
		t.Error("name should be updated")
	}
	if updated.ReadOnly != true {
		t.Error("readOnly should be updated")
	}
}

func TestVolume_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	owner := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "deletevolumeowner",
		Email:        "deletevolumeowner@example.com",
		PasswordHash: "hash",
	}
	owner.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Delete Volume Project",
		OwnerID: owner.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "Delete Volume App",
		CreatedBy: owner.ID,
	}
	app.InsertInDB()

	vol, _ := models.CreateVolume(app.ID, "deleteme", "/host", "/container", false)

	err := models.DeleteVolume(vol.ID)
	if err != nil {
		t.Fatalf("DeleteVolume failed: %v", err)
	}

	result, err := models.GetVolumeByID(vol.ID)
	if err != nil {
		t.Fatalf("GetVolumeByID failed: %v", err)
	}
	if result != nil {
		t.Error("volume should be deleted")
	}
}

func TestTransaction_Success(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	user := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "transuser",
		Email:        "trans@example.com",
		PasswordHash: "hash",
	}
	user.Create()

	err := db.Transaction(func(tx *gorm.DB) error {
		user.Role = "admin"
		return tx.Save(user).Error
	})

	if err != nil {
		t.Fatalf("Transaction failed: %v", err)
	}

	updated, _ := models.GetUserByID(user.ID)
	if updated.Role != "admin" {
		t.Error("role should be updated within transaction")
	}
}

func TestTransaction_Rollback(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	user := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "rollbackuser",
		Email:        "rollback@example.com",
		PasswordHash: "hash",
	}
	user.Create()

	originalRole := user.Role

	err := db.Transaction(func(tx *gorm.DB) error {
		user.Role = "admin"
		if err := tx.Save(user).Error; err != nil {
			return err
		}
		return fmt.Errorf("force rollback")
	})

	if err == nil {
		t.Error("transaction should fail")
	}

	updated, _ := models.GetUserByID(user.ID)
	if updated.Role != originalRole {
		t.Error("changes should be rolled back")
	}
}

func TestSoftDelete(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	user := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "softdeleteuser",
		Email:        "softdelete@example.com",
		PasswordHash: "hash",
	}
	user.Create()

	db.Delete(user)

	var count int64
	db.Model(&models.User{}).Where("id = ?", user.ID).Count(&count)
	if count != 0 {
		t.Error("soft deleted user should not be found with regular query")
	}

	var deleted models.User
	db.Unscoped().Where("id = ?", user.ID).First(&deleted)
	if deleted.ID != user.ID {
		t.Error("soft deleted user should be retrievable with Unscoped")
	}
}

func TestConcurrentWrites(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	user := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "concurrentuser",
		Email:        "concurrent@example.com",
		PasswordHash: "hash",
	}
	user.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Concurrent Project",
		OwnerID: user.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "Concurrent App",
		CreatedBy: user.ID,
	}
	app.InsertInDB()

	results := make(chan error, 10)

	for i := 0; i < 10; i++ {
		go func(idx int) {
			deployment := &models.Deployment{
				ID:         utils.GenerateRandomId(),
				AppID:      app.ID,
				CommitHash: fmt.Sprintf("commit%d", idx),
			}
			results <- deployment.CreateDeployment()
		}(i)
	}

	for i := 0; i < 10; i++ {
		if err := <-results; err != nil {
			t.Errorf("concurrent write failed: %v", err)
		}
	}

	deployments, _ := models.GetDeploymentsByAppID(app.ID)
	if len(deployments) != 10 {
		t.Error("all concurrent deployments should be created")
	}
}

func TestBoundaryValues(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	user := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "boundaryuser",
		Email:        "boundary@example.com",
		PasswordHash: "hash",
	}
	user.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Boundary Project",
		OwnerID: user.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID:   project.ID,
		Name:        "Boundary App",
		CreatedBy:   user.ID,
		CPULimit:    float64Ptr(0.0),
		MemoryLimit: intPtr(0),
	}
	app.InsertInDB()

	retrieved, _ := models.GetApplicationByID(app.ID)
	if retrieved.CPULimit != nil && *retrieved.CPULimit != 0.0 {
		t.Error("zero CPU limit should be preserved")
	}
	if retrieved.MemoryLimit != nil && *retrieved.MemoryLimit != 0 {
		t.Error("zero memory limit should be preserved")
	}
}

func float64Ptr(f float64) *float64 {
	return &f
}

func intPtr(i int) *int {
	return &i
}

func TestLargePayload(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	user := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "largeuser",
		Email:        "large@example.com",
		PasswordHash: "hash",
	}
	user.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Large Project",
		OwnerID: user.ID,
		Tags:    make([]string, 100),
	}
	for i := range project.Tags {
		project.Tags[i] = fmt.Sprintf("tag%d", i)
	}
	project.InsertInDB()

	retrieved, _ := models.GetProjectByID(project.ID)
	if len(retrieved.Tags) != 100 {
		t.Error("all tags should be stored and retrieved")
	}
}

func TestNullFields(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	user := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "nulluser",
		Email:        "null@example.com",
		PasswordHash: "hash",
		FullName:     nil,
		AvatarURL:    nil,
		Bio:          nil,
	}
	user.Create()

	retrieved, _ := models.GetUserByID(user.ID)
	if retrieved.FullName != nil {
		t.Error("nil FullName should remain nil")
	}
	if retrieved.AvatarURL != nil {
		t.Error("nil AvatarURL should remain nil")
	}
	if retrieved.Bio != nil {
		t.Error("nil Bio should remain nil")
	}
}

func TestRelatedRecordsCascade(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	user := &models.User{
		ID:           utils.GenerateRandomId(),
		Username:     "cascadeuser",
		Email:        "cascade@example.com",
		PasswordHash: "hash",
	}
	user.Create()

	project := &models.Project{
		ID:      utils.GenerateRandomId(),
		Name:    "Cascade Project",
		OwnerID: user.ID,
	}
	project.InsertInDB()

	app := &models.App{
		ProjectID: project.ID,
		Name:      "Cascade App",
		CreatedBy: user.ID,
	}
	app.InsertInDB()

	models.CreateDomain(app.ID, "cascade.example.com")
	models.CreateVolume(app.ID, "data", "/host", "/container", false)

	var envCount, domainCount, volumeCount int64
	db.Model(&models.EnvVariable{}).Where("app_id = ?", app.ID).Count(&envCount)
	db.Model(&models.Domain{}).Where("app_id = ?", app.ID).Count(&domainCount)
	db.Model(&models.Volume{}).Where("app_id = ?", app.ID).Count(&volumeCount)

	if envCount != 0 || domainCount != 1 || volumeCount != 1 {
		t.Error("related records should be created")
	}

	models.DeleteApplication(app.ID)

	var appCount int64
	db.Model(&models.App{}).Where("id = ?", app.ID).Count(&appCount)
	if appCount != 0 {
		t.Error("app should be deleted")
	}

}
