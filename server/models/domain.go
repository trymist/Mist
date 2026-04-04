package models

import (
	"time"

	"github.com/corecollectives/mist/utils"
)

type sslStatus string
type sslProvider string
type acmeChallengeType string

const (
	SSLStatusPending  sslStatus = "pending"
	SSLStatusActive   sslStatus = "active"
	SSLStatusFailed   sslStatus = "failed"
	SSLStatusDisabled sslStatus = "disabled"
	SSLStatusExpired  sslStatus = "expired"

	SSLProvider       sslProvider = "letsencrypt"
	SSLProviderCustom sslProvider = "custom"
	SSLProviderNone   sslProvider = "none"

	AcmeChallengeTypeHttp01    acmeChallengeType = "http-01"
	AcmeChallengeTypeDns01     acmeChallengeType = "dns-01"
	AcmeChallengeTypeTlsAlpn01 acmeChallengeType = "tls-alpn-01"
)

type Domain struct {
	ID int64 `gorm:"primaryKey;autoIncrement:false" json:"id"`

	AppID int64 `gorm:"index;not null;constraint:OnDelete:CASCADE" json:"appId"`

	Domain string `gorm:"column:domain_name;uniqueIndex;not null" json:"domain"`

	SslStatus   sslStatus   `gorm:"default:'pending';index" json:"sslStatus"`
	SslProvider sslProvider `gorm:"default:'letsencrypt'" json:"sslProvider,omitempty"`

	CertificatePath *string `json:"certificatePath,omitempty"`
	CertificateData *string `json:"-"`
	KeyPath         *string `json:"keyPath,omitempty"`
	KeyData         *string `json:"-"`
	ChainPath       *string `json:"chainPath,omitempty"`

	AcmeAccountUrl    *string           `json:"acmeAccountUrl,omitempty"`
	AcmeChallengeType acmeChallengeType `json:"acmeChallengeType,omitempty"`

	Issuer             *string    `json:"issuer,omitempty"`
	IssuedAt           *time.Time `json:"issuedAt,omitempty"`
	ExpiresAt          *time.Time `gorm:"index" json:"expiresAt,omitempty"`
	LastRenewalAttempt *time.Time `json:"lastRenewalAttempt,omitempty"`
	RenewalError       *string    `json:"renewalError,omitempty"`

	AutoRenew         bool `gorm:"default:true" json:"autoRenew"`
	ForceHttps        bool `gorm:"default:false" json:"forceHttps"`
	HstsEnabled       bool `gorm:"default:false" json:"hstsEnabled"`
	HstsMaxAge        int  `gorm:"default:31536000" json:"hstsMaxAge"`
	RedirectWww       bool `gorm:"default:false" json:"redirectWww"`
	RedirectWwwToRoot bool `gorm:"default:true" json:"redirectWwwToRoot"`

	DnsConfigured bool       `gorm:"default:false" json:"dnsConfigured"`
	DnsVerifiedAt *time.Time `json:"dnsVerifiedAt,omitempty"`
	LastDnsCheck  *time.Time `json:"lastDnsCheck,omitempty"`
	DnsCheckError *string    `json:"dnsCheckError,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func CreateDomain(appID int64, domain string) (*Domain, error) {
	var d Domain
	id := utils.GenerateRandomId()
	d.ID = id
	d.AppID = appID
	d.Domain = domain
	result := db.Create(&d)
	if result.Error != nil {
		return nil, result.Error
	}
	return &d, nil
}
func GetDomainsByAppID(appID int64) ([]Domain, error) {
	var domains []Domain
	result := db.Where("app_id = ?", appID).Order("created_at ASC").Find(&domains)
	if result.Error != nil {
		return nil, result.Error
	}
	return domains, nil
}

func GetPrimaryDomainByAppID(appID int64) (*Domain, error) {
	var d Domain
	result := db.Where("app_id = ?", appID).Order("created_at ASC").First(&d)
	if result.Error != nil {
		return nil, result.Error
	}
	return &d, nil
}

func UpdateDomain(id int64, domain string) error {
	result := db.Model(&Domain{}).Where("id=?", id).Update("domain_name", domain)
	return result.Error
}

func DeleteDomain(id int64) error {
	result := db.Delete(&Domain{}, id)
	return result.Error
}

func GetDomainByID(id int64) (*Domain, error) {
	var d Domain
	result := db.First(&d, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &d, nil
}

func UpdateDomainDnsStatus(id int64, configured bool, errorMsg *string) error {
	updates := map[string]interface{}{
		"last_dns_check": time.Now(),
	}

	if configured {
		updates["dns_configured"] = true
		updates["dns_verified_at"] = time.Now()
		updates["dns_check_error"] = nil
	} else {
		updates["dns_configured"] = false
		updates["dns_check_error"] = errorMsg
	}

	return db.Model(&Domain{ID: id}).Updates(updates).Error
}

//#############################################################################################################################################
//ARCHIVED CODE BELOW------->

// package models

// import (
// 	"time"

// 	"github.com/corecollectives/mist/utils"
// )

// type Domain struct {
// 	ID            int64      `json:"id"`
// 	AppID         int64      `json:"appId"`
// 	Domain        string     `json:"domain"`
// 	SslStatus     string     `json:"sslStatus"`
// 	DnsConfigured bool       `json:"dnsConfigured"`
// 	DnsVerifiedAt *time.Time `json:"dnsVerifiedAt"`
// 	LastDnsCheck  *time.Time `json:"lastDnsCheck"`
// 	DnsCheckError *string    `json:"dnsCheckError"`
// 	CreatedAt     time.Time  `json:"createdAt"`
// 	UpdatedAt     time.Time  `json:"updatedAt"`
// }

// func CreateDomain(appID int64, domain string) (*Domain, error) {
// 	id := utils.GenerateRandomId()
// 	query := `
// 		INSERT INTO domains (id, app_id, domain_name, ssl_status, dns_configured, created_at, updated_at)
// 		VALUES (?, ?, ?, 'pending', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
// 		RETURNING id, app_id, domain_name, ssl_status, dns_configured, dns_verified_at, last_dns_check, dns_check_error, created_at, updated_at
// 	`
// 	var d Domain
// 	err := db.QueryRow(query, id, appID, domain).Scan(
// 		&d.ID, &d.AppID, &d.Domain, &d.SslStatus, &d.DnsConfigured, &d.DnsVerifiedAt, &d.LastDnsCheck, &d.DnsCheckError, &d.CreatedAt, &d.UpdatedAt,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &d, nil
// }

// func GetDomainsByAppID(appID int64) ([]Domain, error) {
// 	query := `
// 		SELECT id, app_id, domain_name, ssl_status, dns_configured, dns_verified_at, last_dns_check, dns_check_error, created_at, updated_at
// 		FROM domains
// 		WHERE app_id = ?
// 		ORDER BY created_at ASC
// 	`
// 	rows, err := db.Query(query, appID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var domains []Domain
// 	for rows.Next() {
// 		var d Domain
// 		err := rows.Scan(&d.ID, &d.AppID, &d.Domain, &d.SslStatus, &d.DnsConfigured, &d.DnsVerifiedAt, &d.LastDnsCheck, &d.DnsCheckError, &d.CreatedAt, &d.UpdatedAt)
// 		if err != nil {
// 			return nil, err
// 		}
// 		domains = append(domains, d)
// 	}
// 	return domains, nil
// }

// func GetPrimaryDomainByAppID(appID int64) (*Domain, error) {
// 	query := `
// 		SELECT id, app_id, domain_name, ssl_status, dns_configured, dns_verified_at, last_dns_check, dns_check_error, created_at, updated_at
// 		FROM domains
// 		WHERE app_id = ?
// 		ORDER BY created_at ASC
// 		LIMIT 1
// 	`
// 	var d Domain
// 	err := db.QueryRow(query, appID).Scan(&d.ID, &d.AppID, &d.Domain, &d.SslStatus, &d.DnsConfigured, &d.DnsVerifiedAt, &d.LastDnsCheck, &d.DnsCheckError, &d.CreatedAt, &d.UpdatedAt)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &d, nil
// }

// func UpdateDomain(id int64, domain string) error {
// 	query := `
// 		UPDATE domains
// 		SET domain_name = ?, updated_at = CURRENT_TIMESTAMP
// 		WHERE id = ?
// 	`
// 	_, err := db.Exec(query, domain, id)
// 	return err
// }

// func DeleteDomain(id int64) error {
// 	query := `DELETE FROM domains WHERE id = ?`
// 	_, err := db.Exec(query, id)
// 	return err
// }

// func GetDomainByID(id int64) (*Domain, error) {
// 	query := `
// 		SELECT id, app_id, domain_name, ssl_status, dns_configured, dns_verified_at, last_dns_check, dns_check_error, created_at, updated_at
// 		FROM domains
// 		WHERE id = ?
// 	`
// 	var d Domain
// 	err := db.QueryRow(query, id).Scan(&d.ID, &d.AppID, &d.Domain, &d.SslStatus, &d.DnsConfigured, &d.DnsVerifiedAt, &d.LastDnsCheck, &d.DnsCheckError, &d.CreatedAt, &d.UpdatedAt)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &d, nil
// }

// func UpdateDomainDnsStatus(id int64, configured bool, errorMsg *string) error {
// 	var query string
// 	var err error

// 	if configured {
// 		query = `
// 			UPDATE domains
// 			SET dns_configured = 1,
// 				dns_verified_at = CURRENT_TIMESTAMP,
// 				last_dns_check = CURRENT_TIMESTAMP,
// 				dns_check_error = NULL,
// 				updated_at = CURRENT_TIMESTAMP
// 			WHERE id = ?
// 		`
// 		_, err = db.Exec(query, id)
// 	} else {
// 		query = `
// 			UPDATE domains
// 			SET dns_configured = 0,
// 				last_dns_check = CURRENT_TIMESTAMP,
// 				dns_check_error = ?,
// 				updated_at = CURRENT_TIMESTAMP
// 			WHERE id = ?
// 		`
// 		_, err = db.Exec(query, errorMsg, id)
// 	}

// 	return err
// }
