package xiaoai

import "net/http"

type MiClient struct {
	HTTP          *http.Client
	Config        *XiaoaiConfig
	UserAgentMiot string
	UserAgentMina string
}

type XiaoaiConfig struct {
	UserID    string     `json:"userId,omitempty"`
	PassToken string     `json:"passToken,omitempty"`
	Miot      MiotCookie `json:"miot"`
	Mina      MinaCookie `json:"mina"`
}

type MiotCookie struct {
	Ssecurity    string `json:"ssecurity,omitempty"`
	ServiceToken string `json:"serviceToken,omitempty"`
}

type MinaCookie struct {
	MicoapiSlh   string `json:"micoapi_slh,omitempty"`
	MicoapiPh    string `json:"micoapi_ph,omitempty"`
	ServiceToken string `json:"serviceToken,omitempty"`
}

type MiotDeviceListResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Result  struct {
		List []struct {
			AdminFlag int    `json:"adminFlag,omitempty"`
			Bssid     string `json:"bssid,omitempty"`
			Desc      string `json:"desc,omitempty"`
			Did       string `json:"did,omitempty"`
			Extra     struct {
				FwVersion         string `json:"fw_version,omitempty"`
				IsPasswordEncrypt int    `json:"isPasswordEncrypt,omitempty"`
				IsSetPincode      int    `json:"isSetPincode,omitempty"`
				NeedVerifyCode    int    `json:"needVerifyCode,omitempty"`
				PincodeType       int    `json:"pincodeType,omitempty"`
			} `json:"extra,omitempty"`
			FamilyID    int    `json:"family_id,omitempty"`
			InternetIP  string `json:"internet_ip,omitempty"`
			IsOnline    bool   `json:"isOnline,omitempty"`
			Latitude    string `json:"latitude,omitempty"`
			Localip     string `json:"localip,omitempty"`
			Longitude   string `json:"longitude,omitempty"`
			Mac         string `json:"mac,omitempty"`
			Model       string `json:"model,omitempty"`
			Name        string `json:"name,omitempty"`
			P2PID       string `json:"p2p_id,omitempty"`
			ParentID    string `json:"parent_id,omitempty"`
			ParentModel string `json:"parent_model,omitempty"`
			Password    string `json:"password,omitempty"`
			PdID        int    `json:"pd_id,omitempty"`
			PermitLevel int    `json:"permitLevel,omitempty"`
			Pid         string `json:"pid,omitempty"`
			ResetFlag   int    `json:"reset_flag,omitempty"`
			Rssi        int    `json:"rssi,omitempty"`
			ShareFlag   int    `json:"shareFlag,omitempty"`
			ShowMode    int    `json:"show_mode,omitempty"`
			Ssid        string `json:"ssid,omitempty"`
			Token       string `json:"token,omitempty"`
			UID         int64  `json:"uid,omitempty"`
		} `json:"list,omitempty"`
	} `json:"result,omitempty"`
}

type MinaDeviceListResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    []struct {
		DeviceID     string `json:"deviceID,omitempty"`
		SerialNumber string `json:"serialNumber,omitempty"`
		Name         string `json:"name,omitempty"`
		Alias        string `json:"alias,omitempty"`
		Current      bool   `json:"current,omitempty"`
		Presence     string `json:"presence,omitempty"`
		Address      string `json:"address,omitempty"`
		MiotDID      string `json:"miotDID,omitempty"`
		Hardware     string `json:"hardware,omitempty"`
		RomVersion   string `json:"romVersion,omitempty"`
		RomChannel   string `json:"romChannel,omitempty"`
		Capabilities struct {
			MultiroomMusic        int `json:"multiroom_music,omitempty"`
			MultiroomMiplay       int `json:"multiroom_miplay,omitempty"`
			ContentBlacklist      int `json:"content_blacklist,omitempty"`
			NightModeV2           int `json:"night_mode_v2,omitempty"`
			WeakupFeedbackRecord  int `json:"weakup_feedback_record,omitempty"`
			StoreDemoMode         int `json:"store_demo_mode,omitempty"`
			SchoolTimetable       int `json:"school_timetable,omitempty"`
			UserNickName          int `json:"user_nick_name,omitempty"`
			NightMode             int `json:"night_mode,omitempty"`
			PlayerPauseTimer      int `json:"player_pause_timer,omitempty"`
			DialogH5              int `json:"dialog_h5,omitempty"`
			ChildMode2            int `json:"child_mode_2,omitempty"`
			StereoModeV2          int `json:"stereo_mode_v2,omitempty"`
			Dlna                  int `json:"dlna,omitempty"`
			ReportTimes           int `json:"report_times,omitempty"`
			AiInstruction         int `json:"ai_instruction,omitempty"`
			AlarmVolume           int `json:"alarm_volume,omitempty"`
			CustomTts             int `json:"custom_tts,omitempty"`
			ClassifiedAlarm       int `json:"classified_alarm,omitempty"`
			LoadmoreV2            int `json:"loadmore_v2,omitempty"`
			Mesh                  int `json:"mesh,omitempty"`
			AiProtocol30          int `json:"ai_protocol_3_0,omitempty"`
			VoicePrintMultidevice int `json:"voice_print_multidevice,omitempty"`
			NightModeDetail       int `json:"night_mode_detail,omitempty"`
			ChildMode             int `json:"child_mode,omitempty"`
			BabySchedule          int `json:"baby_schedule,omitempty"`
			DidiAuth              int `json:"didi_auth,omitempty"`
			ToneSetting           int `json:"tone_setting,omitempty"`
			Earthquake            int `json:"earthquake,omitempty"`
			NearbyWakeupV2        int `json:"nearby_wakeup_v2,omitempty"`
			AlarmRepeatOptionV2   int `json:"alarm_repeat_option_v2,omitempty"`
			XiaomiVoip            int `json:"xiaomi_voip,omitempty"`
			NearbyWakeupCloud     int `json:"nearby_wakeup_cloud,omitempty"`
			FamilyVoice           int `json:"family_voice,omitempty"`
			BluetoothOptionV2     int `json:"bluetooth_option_v2,omitempty"`
			CustomIr              int `json:"custom_ir,omitempty"`
			Yueyu                 int `json:"yueyu,omitempty"`
			Yunduantts            int `json:"yunduantts,omitempty"`
			StereoMode            int `json:"stereo_mode,omitempty"`
			MicoCurrent           int `json:"mico_current,omitempty"`
			DtsSoundEffect        int `json:"dts_sound_effect,omitempty"`
			VoipUsedTime          int `json:"voip_used_time,omitempty"`
		} `json:"capabilities,omitempty"`
		RemoteCtrlType  string `json:"remoteCtrlType,omitempty"`
		DeviceSNProfile string `json:"deviceSNProfile,omitempty"`
		DeviceProfile   string `json:"deviceProfile,omitempty"`
		BrokerEndpoint  string `json:"brokerEndpoint,omitempty"`
		BrokerIndex     int    `json:"brokerIndex,omitempty"`
		Mac             string `json:"mac,omitempty"`
		Ssid            string `json:"ssid,omitempty"`
	} `json:"data,omitempty"`
}

type MinaGetLatestConversationResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    string `json:"data,omitempty"` // 先用 string 接收
}

type MinaConversationData struct {
	BitSet  []int `json:"bitSet,omitempty"`
	Records []struct {
		BitSet  []int `json:"bitSet,omitempty"`
		Answers []struct {
			BitSet []int  `json:"bitSet,omitempty"`
			Type   string `json:"type,omitempty"`
			Tts    struct {
				BitSet []int  `json:"bitSet,omitempty"`
				Text   string `json:"text,omitempty"`
			} `json:"tts,omitempty"`
		} `json:"answers,omitempty"`
		Time      int64  `json:"time,omitempty"`
		Query     string `json:"query,omitempty"`
		RequestID string `json:"requestId,omitempty"`
	} `json:"records,omitempty"`
	NextEndTime int64 `json:"nextEndTime,omitempty"`
}
