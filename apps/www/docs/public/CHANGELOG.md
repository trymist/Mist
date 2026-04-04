# Changelog
All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog,
and this project adheres to Semantic Versioning.

## [unreleased]

### Added 
- github deployment
- expose port

### Changed
- better restart on app settings change

## [1.0.5] - 2026-01-31

### Added 
- public git url as source
- docker compose support (experimental)
- hide/show env variables

### Fixed
- deployment with dockerfile in children of root directory
- template entry in db on startup
- audit-log no username

## [1.0.4] - 2026-01-20

### Added 
- stop ongoing deployment
- DB tests

### Fixed
- update environment variable issue

### Changed
- switched to gorm

--- 

## [1.0.3] - 2026-01-12

### Added
- deployment cleanup on startup

### Fixed
- infinite logging when file not found

### Changed
- logs UI
- replace raw git commands with go-git
- replace raw docker commands with moby

---

## [1.0.2] - 2025-12-31

### Added
- individual docker container statistics
- owner's email as let's encrypt email on signup

---

## [1.0.1] - 2025-12-27

### Added
- Versioning in dashboard
- `git_clone_url` field for apps

### Fixed
- update script

---

## [1.0.0] - 2025-12-24

### Added
- Initial release
- Core paas functionalities
- web dashboard
- service deployment via docker
