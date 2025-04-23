/*
 *    Copyright 2024 okdp.io
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package config

import (
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/okdp/okdp-server/internal/model"
	"github.com/spf13/viper"
)

// Application configuration
type ApplicationConfig struct {
	Server   Server           `mapstructure:"server"`
	Security Security         `mapstructure:"security"`
	Logging  Logging          `mapstructure:"logging"`
	Swagger  Swagger          `mapstructure:"swagger"`
	Catalogs []*model.Catalog `mapstructure:"catalog"`
}

// Server configuration
type Server struct {
	ListenAddress string `mapstructure:"listenAddress"`
	Port          int    `mapstructure:"port"`
	Mode          string `mapstructure:"mode"`
}

// Logging configuration
type Logging struct {
	Level  string `yaml:"provider"`
	Format string `yaml:"format"`
}

// Security configuration
type Security struct {
	AuthN   AuthN             `yaml:"authN"`
	AuthZ   AuthZ             `yaml:"authZ"`
	Cors    Cors              `yaml:"cors"`
	Headers map[string]string `yaml:"headers"`
}

// Cors configuration
type Cors struct {
	AllowedOrigins   []string `json:"allowedOrigins"`
	AllowedMethods   []string `json:"allowedMethods"`
	AllowedHeaders   []string `json:"allowedHeaders"`
	ExposedHeaders   []string `json:"exposedHeaders"`
	AllowCredentials bool     `json:"allowCredentials"`
	MaxAge           int64    `json:"maxAge"`
}

// Authentication configuration
type AuthN struct {
	Provider []string    `yaml:"provider"`
	OpenID   OpenIDAuth  `yaml:"openid"`
	Bearer   BearerAuth  `yaml:"bearer"`
	Basic    []BasicAuth `yaml:"basic"`
}

// Basic auth based authentication configuration
type BasicAuth struct {
	Login     string   `json:"login"`
	Password  string   `json:"password"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
}

// OpenID based authentication configuration
type OpenIDAuth struct {
	ClientID            string `yaml:"clientId"`
	ClientSecret        string `yaml:"clientSecret"`
	IssuerURI           string `yaml:"issuerUri"`
	RedirectURI         string `yaml:"redirectUri"`
	CookieSecret        string `yaml:"cookieSecret"`
	Scope               string `yaml:"scope"`
	RolesAttributePath  string `yaml:"rolesAttributePath"`
	GroupsAttributePath string `yaml:"groupsAttributePath"`
}

// Bearer based authentication configuration
type BearerAuth struct {
	IssuerURI           string `yaml:"issuerUri"`
	JwksURL             string `yaml:"jwksURL"`
	RolesAttributePath  string `yaml:"rolesAttributePath"`
	GroupsAttributePath string `yaml:"groupsAttributePath"`
	SkipIssuerCheck     bool   `yaml:"skipIssuerCheck"`
	SkipSignatureCheck  bool   `yaml:"skipSignatureCheck"`
}

// AuthZ configuration options
type AuthZ struct {
	Provider string      `yaml:"provider"`
	File     FileAuthZ   `yaml:"file"`
	Database DBAuthZ     `yaml:"database"`
	InLine   InLineAuthZ `yaml:"inline"`
}

// File-based authorization
type FileAuthZ struct {
	ModelPath  string `yaml:"modelPath"`
	PolicyPath string `yaml:"policyPath"`
}

// Database-based authorization
type DBAuthZ struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

// Inline-based authorization
type InLineAuthZ struct {
	Policy string `yaml:"policy"`
	Model  string `yaml:"model"`
}

type Swagger struct {
	SecuritySchemes map[string]*openapi3.SecurityScheme `yaml:"securitySchemes,omitempty"`
	Security        openapi3.SecurityRequirements       `yaml:"security,omitempty"`
}

var (
	instance *ApplicationConfig
	once     sync.Once
)

// GetAppConfig returns a singleton instance of the application configuration.
// It reads the yaml file provided in the argument (--config=/path/to/app-config.yaml) at the startup of the application
// into the ApplicationConfig struct
func GetAppConfig() *ApplicationConfig {
	once.Do(func() {
		instance = &ApplicationConfig{}
		configFile := viper.GetString("config")
		viper.SetConfigFile(configFile)
		fmt.Println("Loading configuration from config file: ", configFile)

		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("failed to read the configuration file")
			panic(err)
		}

		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("Config file changed:", e.Name)
			if err := viper.Unmarshal(&instance); err != nil {
				fmt.Println("failed to register config change watcher")
				panic(err)
			}
		})

		if err := viper.Unmarshal(&instance); err != nil {
			fmt.Println("failed to parse the configuration file")
			panic(err)
		}
	})
	return instance
}

// Visible for testing
func resetAppConfig() {
	instance = nil
	once = sync.Once{}
}
