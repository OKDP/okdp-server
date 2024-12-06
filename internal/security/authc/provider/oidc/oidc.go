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

package oidc

import (
	"encoding/gob"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/okdp/okdp-server/internal/config"
	"github.com/okdp/okdp-server/internal/constants"
	"github.com/okdp/okdp-server/internal/errors"
	log "github.com/okdp/okdp-server/internal/logging"
	"github.com/okdp/okdp-server/internal/security/authc/model"
	"github.com/okdp/okdp-server/internal/utils"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type Provider struct {
	*oidc.Provider
	*oauth2.Config
	*oidc.IDTokenVerifier
	context.Context
	sessions.Store
}

var (
	rolesAttributePath  string
	groupsAttributePath string
)

func init() {
	gob.Register(model.UserInfo{})
}

func NewProvider(oidcConf config.OpenIDAuth) (*Provider, error) {
	ctx := context.Background()
	store := cookie.NewStore([]byte(oidcConf.CookieSecret))
	oidcProvider, err := oidc.NewProvider(ctx, oidcConf.IssuerURI)
	if err != nil {
		return &Provider{}, err
	}
	oidcConfig := &oidc.Config{
		ClientID: oidcConf.ClientID,
	}
	verifier := oidcProvider.Verifier(oidcConfig)

	config := &oauth2.Config{
		ClientID:     oidcConf.ClientID,
		ClientSecret: oidcConf.ClientSecret,
		Endpoint:     oidcProvider.Endpoint(),
		RedirectURL:  oidcConf.RedirectURI,
		Scopes:       strings.Split(oidcConf.Scope, "+"),
	}

	rolesAttributePath = oidcConf.RolesAttributePath
	groupsAttributePath = oidcConf.GroupsAttributePath

	return &Provider{oidcProvider, config, verifier, ctx, store}, nil
}

func (p *Provider) AuthLogin(c *gin.Context) {
	state, err := utils.RandomString()
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.OfType(errors.OkdpServer).GenericError(http.StatusUnauthorized, "Failed to to create OAuth2 state"))
		return
	}
	nonce, err := utils.RandomString()
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.OfType(errors.OkdpServer).GenericError(http.StatusUnauthorized, "Failed to to create OAuth2 nonce"))
		return
	}
	url := p.Config.AuthCodeURL(state, oidc.Nonce(nonce))
	session := sessions.Default(c)
	session.Set(constants.OAuth2State, state)
	session.Set(constants.OAuth2Nonce, nonce)
	if err = session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, errors.OfType(errors.OkdpServer).GenericError(http.StatusInternalServerError, "Failed to save user session in cookie"))
	}
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Auth returns a middleware which authenticates the user with a oidc authentication
// and returns a second middleware which stores the user info in a secure cookie.
// It also propagates the user info (roles/groups) into the autorization provider.
func (p *Provider) Auth() []gin.HandlerFunc {
	return []gin.HandlerFunc{p.authenticate(), p.cookieSessionStore()}
}

func (p *Provider) authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			ok       bool
			userInfo model.UserInfo
		)

		session := sessions.Default(c)
		maybeUserInfo := session.Get(constants.OAuth2UserInfo)
		if userInfo, ok = maybeUserInfo.(model.UserInfo); ok {
			c.Set(constants.OAuth2UserInfo, userInfo)
			log.Debug("The user (Email: %s, Subject: %s) was already authenticated", userInfo.Email, userInfo.Subject)
			c.Next()
			return
		}

		state := session.Get(constants.OAuth2State)
		if c.Query("state") != state {
			log.Warn("Invalid authentication OAuth2 state")
			c.AbortWithStatusJSON(http.StatusBadRequest, errors.OfType(errors.OkdpServer).GenericError(http.StatusBadRequest, "Invalid authentication OAuth2 'state': "+state.(string)))
			return
		}
		// Exchange the authorization code for an access token
		token, err := p.Config.Exchange(p.Context, c.Query("code"))
		if err != nil {
			log.Warn("Failed to exchange the authorization code with an access token: %w", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.OfType(errors.OkdpServer).GenericError(http.StatusUnauthorized, "Failed to exchange authorization code with an access token: "+err.Error()))
			return
		}

		rawIDToken, ok := token.Extra("id_token").(string)
		if !ok {
			log.Warn("No id_token field found in the OAuth2 token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.OfType(errors.OkdpServer).GenericError(http.StatusUnauthorized, "No id_token field found in the OAuth2 token"))
			return
		}
		idToken, err := p.Verify(p.Context, rawIDToken)
		if err != nil {
			log.Warn("Failed to verify the ID Token: %w", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.OfType(errors.OkdpServer).GenericError(http.StatusUnauthorized, "Failed to verify the ID Token: "+err.Error()))
			return
		}

		nonce := session.Get(constants.OAuth2Nonce)
		if idToken.Nonce != nonce {
			log.Warn("Invalid authentication OAuth2 'nonce': %s", nonce.(string))
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.OfType(errors.OkdpServer).GenericError(http.StatusUnauthorized, "Invalid authentication OAuth2 'nonce': "+nonce.(string)))
			return
		}
		userInfo, err = p.getUserInfo(token.AccessToken)
		if err != nil {
			log.Warn("Unable to get user roles/groups from the access token: %w", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.OfType(errors.OkdpServer).GenericError(http.StatusUnauthorized, "Unable to get user roles/groups from access token: "+err.Error()))
			return
		}
		// Retrieve the user information from the access token
		// client := provider.Config.Client(provider.Ctx, token)
		// resp, err := client.Get(provider.Provider.UserInfoEndpoint())
		// if err != nil {
		// 	c.JSON(http.StatusUnauthorized, errors.OfType(errors.OKDP_SERVER).GenericError(http.StatusUnauthorized, "Failed to retrieve user information"))
		// 	return
		// }
		// defer resp.Body.Close()
		// err = json.NewDecoder(resp.Body).Decode(&userInfo)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, errors.OfType(errors.OKDP_SERVER).GenericError(http.StatusUnauthorized, "Failed to parse user information" + err.Error()))
		// 	return
		// }
		// if err := session.Save(); err != nil {
		// 	c.JSON(http.StatusInternalServerError, errors.OfType(errors.OKDP_SERVER).GenericError(http.StatusUnauthorized, "Failed to save user information" + err.Error()))
		// }
		c.Set(constants.OAuth2UserInfo, userInfo)
		session.Set(constants.OAuth2UserInfo, userInfo)
		log.Debug("Successfully authenticated user : %s", userInfo.AsJSONString())
		c.Next()
	}
}

func (p *Provider) getUserInfo(accessToken string) (model.UserInfo, error) {
	token := &model.Token{AccessToken: accessToken}
	return token.GetUserInfo(rolesAttributePath, groupsAttributePath)
}

func (p *Provider) cookieSessionStore() gin.HandlerFunc {
	return sessions.Sessions(constants.OAuth2SessionName, p.Store)
}
