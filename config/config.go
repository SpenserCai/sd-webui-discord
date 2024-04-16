/*
 * @Author: SpenserCai
 * @Date: 2023-08-16 11:05:40
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-13 13:49:25
 * @Description: file content
 */
package config

type ServerItem struct {
	Name          string `json:"name"`
	Host          string `json:"host"`
	MaxConcurrent int    `json:"max_concurrent"`
	MaxQueue      int    `json:"max_queue"`
	MaxVRAM       string `json:"max_vram"`
}

type StableConfig struct {
	Model          string  `json:"sd_model_checkpoint"`
	Vae            string  `json:"sd_vae"`
	Height         int64   `json:"height"`
	Width          int64   `json:"width"`
	Steps          int64   `json:"steps"`
	CfgScale       float64 `json:"cfg_scale"`
	NegativePrompt string  `json:"negative_prompt"`
	Sampler        string  `json:"sampler"`
	ClipSkip       int64   `json:"clip_skip"`
	NIter          int64   `json:"n_iter"`
	Seed           int64   `json:"seed"`
	MaxLoad        int64   `json:"max_load"`

	Sag          string  `json:"sag"`
	SagScale     float64 `json:"sag_scale"`
	SagBlurSigma float64 `json:"sag_blur_sigma"`

	FreeU   string  `json:"freeu"`
	FreeUB1 float64 `json:"freeu_b1"`
	FreeUB2 float64 `json:"freeu_b2"`
	FreeUS1 float64 `json:"freeu_s1"`
	FreeUS2 float64 `json:"freeu_s2"`
}

type DbConfig struct {
	Type string `json:"type"`
	DSN  string `json:"dsn"`
}

type UserCenter struct {
	Enable       bool     `json:"enable"`
	DbConfig     DbConfig `json:"db_config"`
	MustRegister bool     `json:"must_register"`
	EnableCredit bool     `json:"enable_credit"`
}

type Preset struct {
	Name     string       `json:"name"`
	Settings StableConfig `json:"settings"`
}

type WebSite struct {
	Api struct {
		Host            string `json:"host"`
		Port            int    `json:"port"`
		JwtSecret       string `json:"jwt_secret"`
		AuthCallbackUrl string `json:"auth_callback_url"`
	} `json:"api"`
	Web struct {
		Host            string `json:"host"`
		Port            int    `json:"port"`
		StartWithServer bool   `json:"start_with_server"`
		OpenBrowser     bool   `json:"open_browser"`
	} `json:"web"`
}

type Config struct {
	SDWebUi struct {
		Servers        []ServerItem `json:"servers"`
		DefaultSetting StableConfig `json:"default_setting"`
		Presets        []Preset     `json:"presets"`
	} `json:"sd_webui"`
	Discord struct {
		AppId             string `json:"app_id"`
		Token             string `json:"token"`
		ClientSecret      string `json:"client_secret"`
		OAuth2RedirectUrl string `json:"oauth2_redirect_base_url"`
		ServerId          string `json:"server_id"`
		BotName           string `json:"bot_name"`
		BotAvatar         string `json:"bot_avatar"`
		BotUrl            string `json:"bot_url"`
		DiscordServerUrl  string `json:"discord_server_url"`
	} `json:"discord"`
	UserCenter           UserCenter `json:"user_center"`
	DisableReturnGenInfo bool       `json:"disable_return_gen_info"`
	WebSite              WebSite    `json:"website"`
}
