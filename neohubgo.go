package neohubgo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
)

type Options struct {
	Username   string
	Password   string
	Url        string // defaults to DEFAULT_URL
	HttpClient *http.Client
}

type Device struct {
	Address    string `json:"address"`
	DeviceId   string `json:"deviceid"`
	DeviceName string `json:"devicename"`
	HubType    int    `json:"hub_type"`
	Online     bool   `json:"online"`
	Type       string `json:"type"`
	Version    int    `json:"version"`
}

type Data struct {
	Status     int            `json:"STATUS"`
	CacheValue DataCacheValue `json:"CACHE_VALUE"`
}

type DataCacheValue struct {
	LiveInfo  LiveInfo  `json:"live_info"`
	System    System    `json:"system"`
	Engineers Engineers `json:"engineers"`
}

type LiveInfo struct {
	TimestampEngineers      int              `json:"TIMESTAMP_ENGINEERS"`
	OpenDelay               int              `json:"OPEN_DELAY"`
	TimestampRecipes        int              `json:"TIMESTAMP_RECIPES"`
	HubHoliday              bool             `json:"HUB_HOLIDAY"`
	TimestampProfileTimers0 int              `json:"TIMESTAMP_PROFILE_TIMERS_0"`
	HolidayEnd              int              `json:"HOLIDAY_END"`
	HubTime                 int              `json:"HUB_TIME"`
	TimestampSystem         int              `json:"TIMESTAMP_SYSTEM"`
	Devices                 []LiveInfoDevice `json:"devices"`
}

type LiveInfoDevice struct {
	HeatOn                  bool     `json:"HEAT_ON"`
	CurrentFloorTemperature float32  `json:"CURRENT_FLOOR_TEMPERATURE"`
	Standby                 bool     `json:"STANDBY"`
	ManualOff               bool     `json:"MANUAL_OFF"`
	TimerOn                 bool     `json:"TIMER_ON"`
	WindowOpen              bool     `json:"WINDOW_OPEN"`
	AvailableModes          []string `json:"AVAILABLE_MODES"`
	WriteCount              int      `json:"WRITE_COUNT"`
	FloorLimit              bool     `json:"FLOOR_LIMIT"`
	Date                    string   `json:"DATE"`
	HeatMode                bool     `json:"HEAT_MODE"`
	Offline                 bool     `json:"OFFLINE"`
	Holiday                 bool     `json:"HOLIDAY"`
	ModeLock                bool     `json:"MODELOCK"`
	DeviceID                int      `json:"DEVICE_ID"`
	RecentTemps             []string `json:"RECENT_TEMPS"`
	CoolOn                  bool     `json:"COOL_ON"`
	RelativeHumidity        int      `json:"RELATIVE_HUMIDITY"`
	HoldCool                float32  `json:"HOLD_COOL"`
	Away                    bool     `json:"AWAY"`
	TimeClock               bool     `json:"TIMECLOCK"`
	TemporarySetFlag        bool     `json:"TEMPORARY_SET_FLAG"`
	PrgTimer                bool     `json:"PRG_TIMER"`
	Lock                    bool     `json:"LOCK"`
	ModulationLevel         int      `json:"MODULATION_LEVEL"`
	HcMode                  string   `json:"HC_MODE"`
	SetTemp                 string   `json:"SET_TEMP"`
	LowBattery              bool     `json:"LOW_BATTERY"`
	CoolTemp                float32  `json:"COOL_TEMP"`
	HoldOn                  bool     `json:"HOLD_ON"`
	HoldOff                 bool     `json:"HOLD_OFF"`
	ZoneName                string   `json:"ZONE_NAME"`
	HoldTime                string   `json:"HOLD_TIME"`
	CoolMode                bool     `json:"COOL_MODE"`
	HoldTemp                float32  `json:"HOLD_TEMP"`
	PreheatActive           bool     `json:"PREHEAT_ACTIVE"`
	ActiveProfile           int      `json:"ACTIVE_PROFILE"`
	SwitchDelayLeft         string   `json:"SWITCH_DELAY_LEFT"`
	PinNumber               string   `json:"PIN_NUMBER"`
	FanSpeed                string   `json:"FAN_SPEED"`
	ActualTemp              string   `json:"ACTUAL_TEMP"`
	Time                    string   `json:"TIME"`
	ActiveLevel             int      `json:"ACTIVE_LEVEL"`
	FanControl              string   `json:"FAN_CONTROL"`
	PrgTemp                 float32  `json:"PRG_TEMP"`
}

type System struct {
	DstAuto          bool    `json:"DST_AUTO"`
	Format           int     `json:"FORMAT"`
	HubVersion       int     `json:"HUB_VERSION"`
	CoolboxOverride  *string `json:"COOLBOX_OVERRIDE"`
	DeviceID         string  `json:"DEVICE_ID"`
	AltTimerFormat   *string `json:"ALT_TIMER_FORMAT"`
	TimeZone         float32 `json:"TIME_ZONE"`
	ExtendedHistory  string  `json:"EXTENDED_HISTORY"`
	GlobalSystemType string  `json:"GLOBAL_SYSTEM_TYPE"`
	ZigbeeChannel    string  `json:"ZIGBEE_CHANNEL"`
	LegacyLocalPort  bool    `json:"LEGACY_LOCAL_PORT"`
	HeatingLevels    int     `json:"HEATING_LEVELS"`
	HubType          int     `json:"HUB_TYPE"`
	Utc              int     `json:"UTC"`
	Corf             string  `json:"CORF"`
	TimezoneStr      string  `json:"TIMEZONESTR"`
	Partition        string  `json:"PARTITION"`
	NtpOn            string  `json:"NTP_ON"`
	GlobalHcMode     string  `json:"GLOBAL_HC_MODE"`
	DstOn            bool    `json:"DST_ON"`
	CoolboxPresent   int     `json:"COOLBOX_PRESENT"`
	NoRfBroadcast    bool    `json:"NO_RF_BROADCAST"`
	Timestamp        int     `json:"TIMESTAMP"`
	Gdevlist         []int   `json:"GDEVLIST"`
	Coolbox          string  `json:"COOLBOX"`
}

type Engineer struct {
	SwitchDelay           int     `json:"SWITCH_DELAY"`
	OutputDelay           int     `json:"OUTPUT_DELAY"`
	UltraVersion          int     `json:"ULTRA_VERSION"`
	DewPoint              bool    `json:"DEW_POINT"`
	SensorMode            string  `json:"SENSOR_MODE"`
	CoolEnable            bool    `json:"COOL_ENABLE"`
	Timestamp             int     `json:"TIMESTAMP"`
	UserLimit             int     `json:"USER_LIMIT"`
	StatFailsafe          int     `json:"STAT_FAILSAFE"`
	StatVersion           int     `json:"STAT_VERSION"`
	MaxPreheat            int     `json:"MAX_PREHEAT"`
	PumpDelay             int     `json:"PUMP_DELAY"`
	RfSensorMode          string  `json:"RF_SENSOR_MODE"`
	SystemType            int     `json:"SYSTEM_TYPE"`
	SwitchingDifferential int     `json:"SWITCHING DIFFERENTIAL"`
	DeviceType            int     `json:"DEVICE_TYPE"`
	FrostTemp             float32 `json:"FROST_TEMP"`
	FloorLimit            int     `json:"FLOOR_LIMIT"`
	WindowSwitchOpen      bool    `json:"WINDOW_SWITCH_OPEN"`
	Deadband              int     `json:"DEADBAND"`
	DeviceID              int     `json:"DEVICE_ID"`
}

type Engineers map[string]Engineer

const (
	DEFAULT_URL                 = "https://neohub.co.uk/"
	USER_LOGIN_ENDPOINT         = "hm_user_login"
	CACHE_VALUE_ENDPOINT        = "hm_cache_value"
	DEFAULT_CACHE_VALUE_REQUEST = "engineers,comfort,profile0,timeclock0,system,device_list,timeclock,live_info"
)

type NeoHub interface {
	Login() ([]Device, error)
	GetData(deviceId string) (*Data, error)
}

func NewNeoHub(options *Options) NeoHub {
	client := resty.New()
	if options == nil {
		options = &Options{}
	}
	if options.HttpClient != nil {
		client = resty.NewWithClient(options.HttpClient)
	}
	if options.Url == "" {
		options.Url = DEFAULT_URL
	}

	return &neoHub{
		client:  client,
		options: options,
	}
}

type neoHub struct {
	options *Options
	token   string
	client  *resty.Client
}

type cacheValueRequest struct {
	CacheValue string `url:"cache_value"`
	DeviceId   string `url:"device_id"`
	Token      string `url:"token"`
}

func (n *neoHub) GetData(deviceId string) (*Data, error) {
	var data Data
	req := cacheValueRequest{DEFAULT_CACHE_VALUE_REQUEST, deviceId, n.token}
	err := n.formPostRequest(CACHE_VALUE_ENDPOINT, req, &data)
	if err != nil {
		return nil, err
	}
	if data.Status != 1 {
		return nil, fmt.Errorf("Unexpected status %d", data.Status)
	}

	return &data, nil
}

type loginRequest struct {
	Username string `url:"USERNAME"`
	Password string `url:"PASSWORD"`
}

type loginResponse struct {
	Status  int      `json:"STATUS"`
	Token   string   `json:"TOKEN"`
	UID     int      `json:"UID"`
	Devices []Device `json:"devices"`
}

func (n *neoHub) Login() ([]Device, error) {
	var loginResponse loginResponse
	err := n.formPostRequest(USER_LOGIN_ENDPOINT, loginRequest{n.options.Username, n.options.Password}, &loginResponse)
	if err != nil {
		return nil, err
	}
	if loginResponse.Status != 1 {
		return nil, fmt.Errorf("Unexpected login response: %d", loginResponse.Status)
	}
	n.token = loginResponse.Token
	return loginResponse.Devices, nil
}

func (n *neoHub) formPostRequest(endpoint string, request interface{}, response interface{}) error {
	form, err := query.Values(request)
	if err != nil {
		return fmt.Errorf("Parse error on request: %s", err.Error())
	}

	resp, err := n.client.R().SetFormDataFromValues(form).Post(n.options.Url + endpoint)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("Unexpected status code in request to %s: %s", endpoint, resp.StatusCode())
	}
	err = json.Unmarshal(resp.Body(), response)
	if err != nil {
		return fmt.Errorf("JSON unmarshal error: %s", err.Error())
	}

	return nil
}
