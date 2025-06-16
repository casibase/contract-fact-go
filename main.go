// Copyright 2025 The Casibase Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

type FactContract struct {
}

func (f *FactContract) InitContract() protogo.Response {
	return sdk.Success([]byte("Init contract success"))
}

func (f *FactContract) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade contract success"))
}

func (f *FactContract) InvokeContract(method string) protogo.Response {
	switch method {
	case "save":
		return f.save()
	case "get":
		return f.get()
	default:
		return sdk.Error("invalid method")
	}
}

func (f *FactContract) save() protogo.Response {
	params := sdk.Instance.GetArgs()

	key, ok := params["key"]
	if !ok {
		return sdk.Error("arg key not exist")
	}
	field, ok := params["field"]
	if !ok {
		return sdk.Error("arg field not exist")
	}
	value, ok := params["value"]
	if !ok {
		return sdk.Error("arg value not exist")
	}

	sdk.Instance.EmitEvent("save", []string{string(key), string(field), string(value)})

	err := sdk.Instance.PutStateByte(string(key), string(field), value)
	if err != nil {
		return sdk.Error("fail to save:" + err.Error())
	}

	sdk.Instance.Infof("[save] key=" + string(key) + ",field=" + string(field) + ",value=" + string(value))

	return sdk.Success([]byte("Success"))
}

func (f *FactContract) get() protogo.Response {
	params := sdk.Instance.GetArgs()

	key, ok := params["key"]
	if !ok {
		return sdk.Error("arg key not exist")
	}
	field, ok := params["field"]
	if !ok {
		return sdk.Error("arg field not exist")
	}

	value, err := sdk.Instance.GetStateByte(string(key), string(field))
	if err != nil {
		return sdk.Error("failed to call get_state")
	}

	sdk.Instance.Infof("[get] key=" + string(key) + ",field=" + string(field) + ",value=" + string(value))

	return sdk.Success(value)
}

func main() {
	err := sandbox.Start(new(FactContract))
	if err != nil {
		sdk.Instance.Errorf(err.Error())
	}
}
