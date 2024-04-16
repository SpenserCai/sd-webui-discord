/*
 * @Author: ste-art
 * @Date: 2024-04-09 01:50:02
 * @version:
 * @LastEditors: ste-art
 * @LastEditTime: 2024-04-09 01:50:02
 * @Description: Feature gates that enables or disabled features depending on sd-webui configuration
 */
package global

import "github.com/SpenserCai/sd-webui-go/intersvc"

type FeatureGate struct {
	AwailableScripts *[]string
	Sag              bool
	FreeU            bool
}

func InitFeatureGate() *FeatureGate {
	fg := &FeatureGate{}
	scriptsSvc := &intersvc.SdapiV1Scripts{}
	scriptsSvc.Action(ClusterManager.GetNodeAuto().StableClient)
	if scriptsSvc.Error != nil {
		return fg
	}
	scripts := scriptsSvc.GetResponse()
	scriptsResult := make([]string, len(scripts.Txt2img))
	for i, script := range scripts.Txt2img {
		scriptsResult[i] = script.(string)
	}
	fg.AwailableScripts = &scriptsResult

	for _, script := range *fg.AwailableScripts {
		switch script {
		case "selfattentionguidance integrated":
			fg.Sag = true
		case "freeu integrated":
			fg.FreeU = true
		}
	}

	return fg
}
