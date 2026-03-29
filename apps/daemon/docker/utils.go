package docker

import (
	"fmt"

	"github.com/corecollectives/mist/models"
)

type EnvironmentVariableSet struct {
	BuildTime map[string]string // required for image build
	Runtime   map[string]string // required at time of container build
}

func FetchDeploymentConfigurationForApp(app *models.App) (int, []string, *EnvironmentVariableSet, error) {
	port := 3000
	if app.Port != nil {
		port = int(*app.Port)
	}

	domains, err := models.GetDomainsByAppID(app.ID)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return 0, nil, nil, fmt.Errorf("get domains failed: %w", err)
	}

	var domainStrings []string
	for _, d := range domains {
		domainStrings = append(domainStrings, d.Domain)
	}

	envs, err := models.GetEnvVariablesByAppID(app.ID)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return 0, nil, nil, fmt.Errorf("get env variables failed: %w", err)
	}

	envSet := CategorizeEnvironmentVariables(envs)

	return port, domainStrings, envSet, nil
}

func CategorizeEnvironmentVariables(envs []models.EnvVariable) *EnvironmentVariableSet {
	buildTimeVars := make(map[string]string)
	runtimeVars := make(map[string]string)

	for _, env := range envs {
		isRuntime := env.IsRuntime()
		isBuildtime := env.IsBuildtime()

		// for reverse compatibility
		// if marked as both or neither, add to both
		// if marked only as runtime, add to runtime
		// if marked only as buildtime, add to buildtime
		if isBuildtime || (!isRuntime && !isBuildtime) {
			buildTimeVars[env.Key] = env.Value
		}
		if isRuntime || (!isRuntime && !isBuildtime) {
			runtimeVars[env.Key] = env.Value
		}
	}

	return &EnvironmentVariableSet{
		BuildTime: buildTimeVars,
		Runtime:   runtimeVars,
	}
}

func (e *EnvironmentVariableSet) GetTotalEnvVarCount() int {
	uniqueKeys := make(map[string]bool)
	for k := range e.BuildTime {
		uniqueKeys[k] = true
	}
	for k := range e.Runtime {
		uniqueKeys[k] = true
	}
	return len(uniqueKeys)
}

func (e *EnvironmentVariableSet) GetBuildTimeCount() int {
	return len(e.BuildTime)
}

func (e *EnvironmentVariableSet) GetRuntimeCount() int {
	return len(e.Runtime)
}

func GetDeploymentConfigForApp(app *models.App) (int, []string, map[string]string, error) {
	port, domains, envSet, err := FetchDeploymentConfigurationForApp(app)
	if err != nil {
		return 0, nil, nil, err
	}

	// merge both sets for backward compatibility
	merged := make(map[string]string)
	for k, v := range envSet.BuildTime {
		merged[k] = v
	}
	for k, v := range envSet.Runtime {
		merged[k] = v
	}

	return port, domains, merged, nil
}
