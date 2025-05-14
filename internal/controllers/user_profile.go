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

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/okdp/okdp-server/internal/common/constants"
	"github.com/okdp/okdp-server/internal/model"
	authc "github.com/okdp/okdp-server/internal/security/authc/model"
)

type IUserProfileController struct {
}

func UserProfileController() *IUserProfileController {
	return &IUserProfileController{}
}

func (r IUserProfileController) GetMyProfile(c *gin.Context) {

	if maybeUserInfo, found := c.Get(constants.OAuth2UserInfo); found {
		c.JSON(http.StatusOK, maybeUserInfo.(*authc.UserInfo))
	}

}

func GetUserInfo(c *gin.Context) (*authc.UserInfo, *model.ServerResponse) {
	maybeUserInfo, found := c.Get(constants.OAuth2UserInfo)
	err := model.
		NewServerResponse(model.OkdpServerResponse).
		Unauthorized("Unauthorized")
	if !found {
		return nil, err
	}

	userInfo, ok := maybeUserInfo.(*authc.UserInfo)
	if !ok {
		return nil, err
	}

	return userInfo, nil
}
