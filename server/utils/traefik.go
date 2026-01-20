package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

const (
	TraefikConfigDir   = "/var/lib/mist/traefik"
	TraefikDynamicFile = "dynamic.yml"
	TraefikStaticDir   = "/opt/mist"
	TraefikStaticFile  = "traefik-static.yml"
)

func InitializeTraefikConfig(wildcardDomain *string, mistAppName string) error {
	return GenerateDynamicConfig(wildcardDomain, mistAppName)
}

func GenerateDynamicConfig(wildcardDomain *string, mistAppName string) error {
	if err := os.MkdirAll(TraefikConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to create traefik config directory: %w", err)
	}

	dynamicConfigPath := filepath.Join(TraefikConfigDir, TraefikDynamicFile)
	content, err := generateDynamicYAML(wildcardDomain, mistAppName)
	if err != nil {
		return fmt.Errorf("failed to generate dynamic YAML: %w", err)
	}

	if err := os.WriteFile(dynamicConfigPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write dynamic config: %w", err)
	}

	log.Info().
		Str("path", dynamicConfigPath).
		Msg("Generated Traefik dynamic config")

	return nil
}

func generateDynamicYAML(wildcardDomain *string, mistAppName string) ([]byte, error) {
	cfg := map[string]any{
		"http": map[string]any{
			"routers":  map[string]any{},
			"services": map[string]any{},
			"middlewares": map[string]any{
				"https-redirect": map[string]any{
					"redirectScheme": map[string]any{
						"scheme":    "https",
						"permanent": true,
					},
				},
			},
		},
	}

	if wildcardDomain == nil || *wildcardDomain == "" {
		return yaml.Marshal(cfg)
	}

	domain := strings.TrimPrefix(*wildcardDomain, "*")
	domain = strings.TrimPrefix(domain, ".")
	
	mistDomain := mistAppName + "." + domain

	httpConfig := cfg["http"].(map[string]any)
	httpConfig["routers"] = map[string]any{
		"mist-dashboard": map[string]any{
			"rule":        fmt.Sprintf("Host(`%s`)", mistDomain),
			"entryPoints": []string{"websecure"},
			"service":     "mist-dashboard",
			"tls": map[string]any{
				"certResolver": "le",
			},
		},
		"mist-dashboard-http": map[string]any{
			"rule":        fmt.Sprintf("Host(`%s`)", mistDomain),
			"entryPoints": []string{"web"},
			"middlewares": []string{"https-redirect"},
			"service":     "mist-dashboard",
		},
	}

	httpConfig["services"] = map[string]any{
		"mist-dashboard": map[string]any{
			"loadBalancer": map[string]any{
				"servers": []map[string]any{
					{"url": "http://172.17.0.1:8080"},
				},
			},
		},
	}

	return yaml.Marshal(cfg)
}

func ChangeLetsEncryptEmail(email string) error {
	staticConfigPath := path.Join(TraefikStaticDir, TraefikStaticFile)

	content, err := os.ReadFile(staticConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read traefik-static.yml: %w", err)
	}

	var config yaml.Node
	if err := yaml.Unmarshal(content, &config); err != nil {
		return fmt.Errorf("failed to parse traefik-static.yml: %w", err)
	}

	emailUpdated := false
	if len(config.Content) > 0 {
		rootNode := config.Content[0]
		for i := 0; i < len(rootNode.Content); i += 2 {
			if rootNode.Content[i].Value == "certificatesResolvers" && i+1 < len(rootNode.Content) {
				certResolvers := rootNode.Content[i+1]
				for j := 0; j < len(certResolvers.Content); j += 2 {
					if certResolvers.Content[j].Value == "le" && j+1 < len(certResolvers.Content) {
						leNode := certResolvers.Content[j+1]
						for k := 0; k < len(leNode.Content); k += 2 {
							if leNode.Content[k].Value == "acme" && k+1 < len(leNode.Content) {
								acmeNode := leNode.Content[k+1]
								for l := 0; l < len(acmeNode.Content); l += 2 {
									if acmeNode.Content[l].Value == "email" && l+1 < len(acmeNode.Content) {
										acmeNode.Content[l+1].Value = email
										emailUpdated = true
										break
									}
								}
							}
							if emailUpdated {
								break
							}
						}
					}
					if emailUpdated {
						break
					}
				}
			}
			if emailUpdated {
				break
			}
		}
	}

	if !emailUpdated {
		return fmt.Errorf("email field not found in traefik-static.yml")
	}

	updatedContent, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("failed to marshal updated config: %w", err)
	}

	if err := os.WriteFile(staticConfigPath, updatedContent, 0644); err != nil {
		return fmt.Errorf("failed to write traefik-static.yml: %w", err)
	}

	log.Info().
		Str("email", email).
		Str("path", staticConfigPath).
		Msg("Updated Let's Encrypt email in traefik-static.yml")

	if err := RestartTraefik(); err != nil {
		return fmt.Errorf("failed to restart Traefik: %w", err)
	}

	return nil
}

func RestartTraefik() error {
	log.Info().Msg("Restarting Traefik container...")

	// NOTE: we still use the exec method here because moby doesn't support docker-compose for now
	cmd := exec.Command("docker", "compose", "-f", "/opt/mist/traefik-compose.yml", "restart", "traefik")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().
			Err(err).
			Str("output", string(output)).
			Msg("Failed to restart Traefik container")
		return fmt.Errorf("docker compose restart failed: %w", err)
	}

	log.Info().
		Str("output", string(output)).
		Msg("Traefik container restarted successfully")

	return nil
}
