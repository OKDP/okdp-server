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

package model

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

type Token struct {
	AccessToken string
}

func (t *Token) GetUserInfo(rolesAttributePath string, groupsAttributePath string) (UserInfo, error) {

	var userInfo UserInfo
	accessTokenPayload := strings.Split(t.AccessToken, ".")[1]
	accessTokenDecoded, err := base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(accessTokenPayload)
	if err != nil {
		return userInfo, err
	}
	err = json.Unmarshal(accessTokenDecoded, &userInfo)
	if err != nil {
		return userInfo, err
	}

	err = attributeFromPath(rolesAttributePath, string(accessTokenDecoded), &userInfo.Roles)
	if err != nil {
		return userInfo, fmt.Errorf("%s: %w", rolesAttributePath, err)
	}
	err = attributeFromPath(groupsAttributePath, string(accessTokenDecoded), &userInfo.Groups)
	if err != nil {
		return userInfo, fmt.Errorf("%s: %w", groupsAttributePath, err)
	}

	return userInfo, nil

}

func attributeFromPath(attribute string, accessTokenDecoded string,
	str *[]string) error {
	value := gjson.Get(accessTokenDecoded, attribute)
	if value.Exists() {
		err := json.Unmarshal([]byte(value.String()), &str)
		if err != nil {
			return err
		}
	} else {
		*str = []string{}
	}
	return nil
}
