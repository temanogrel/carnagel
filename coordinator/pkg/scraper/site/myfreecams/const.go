package myfreecams

const (
	wsOrigin    = "http://mobile.mfc.com"
	wsUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36)"

	greetMessage = "hello fcserver\n"
	loginMessage = "1 0 0 81 0 %7B%22err%22%3A0%2C%22start%22%3A1497188853157%2C%22stop%22%3A1497188853642%2C%22a%22%3A13230%2C%22time%22%3A1497188853%2C%22key%22%3A%22f7cc5a5f64e0d6cc192b26eb5dd0cb76%22%2C%22cid%22%3A%222b44de9c%22%7D\n"
	pingMessage  = "0 0 0 0 0"
)

type IncomingCommand int

const (
	CommandLogin  IncomingCommand = 1
	CommandRxData                 = 81

	// Copy paste from reverse engineered code
	FCTYPE_NULL           = 0
	FCTYPE_LOGIN          = 1
	FCTYPE_ADDFRIEND      = 2
	FCTYPE_PMESG          = 3
	FCTYPE_STATUS         = 4
	FCTYPE_DETAILS        = 5
	FCTYPE_TOKENINC       = 6
	FCTYPE_ADDIGNORE      = 7
	FCTYPE_PRIVACY        = 8
	FCTYPE_ADDFRIENDREQ   = 9
	FCTYPE_USERNAMELOOKUP = 10
	FCTYPE_ZBAN           = 11
	FCTYPE_BROADCASTNEWS  = 12
	FCTYPE_ANNOUNCE       = 13
	FCTYPE_MANAGELIST     = 14
	FCTYPE_INBOX          = 15
	FCTYPE_GWCONNECT      = 16
	FCTYPE_RELOADSETTINGS = 17
	FCTYPE_HIDEUSERS      = 18
	FCTYPE_RULEVIOLATION  = 19
	FCTYPE_SESSIONSTATE   = 20
	FCTYPE_REQUESTPVT     = 21
	FCTYPE_ACCEPTPVT      = 22
	FCTYPE_REJECTPVT      = 23
	FCTYPE_ENDSESSION     = 24
	FCTYPE_TXPROFILE      = 25
	FCTYPE_STARTVOYEUR    = 26
	FCTYPE_SERVERREFRESH  = 27
	FCTYPE_SETTING        = 28
	FCTYPE_BWSTATS        = 29
	FCTYPE_TKX            = 30
	FCTYPE_SETTEXTOPT     = 31
	FCTYPE_SERVERCONFIG   = 32
	FCTYPE_MODELGROUP     = 33
	FCTYPE_REQUESTGRP     = 34
	FCTYPE_STATUSGRP      = 35
	FCTYPE_GROUPCHAT      = 36
	FCTYPE_CLOSEGRP       = 37
	FCTYPE_UCR            = 38
	FCTYPE_MYUCR          = 39
	FCTYPE_SLAVECON       = 40
	FCTYPE_SLAVECMD       = 41
	FCTYPE_SLAVEFRIEND    = 42
	FCTYPE_SLAVEVSHARE    = 43
	FCTYPE_ROOMDATA       = 44
	FCTYPE_NEWSITEM       = 45
	FCTYPE_GUESTCOUNT     = 46
	FCTYPE_PRELOGINQ      = 47
	FCTYPE_MODELGROUPSZ   = 48
	FCTYPE_ROOMHELPER     = 49
	FCTYPE_CMESG          = 50
	FCTYPE_JOINCHAN       = 51
	FCTYPE_CREATECHAN     = 52
	FCTYPE_INVITECHAN     = 53
	FCTYPE_QUIETCHAN      = 55
	FCTYPE_BANCHAN        = 56
	FCTYPE_PREVIEWCHAN    = 57
	FCTYPE_SHUTDOWN       = 58
	FCTYPE_LISTBANS       = 59
	FCTYPE_UNBAN          = 60
	FCTYPE_SETWELCOME     = 61
	FCTYPE_CHANOP         = 62
	FCTYPE_LISTCHAN       = 63
	FCTYPE_TAGS           = 64
	FCTYPE_SETPCODE       = 65
	FCTYPE_SETMINTIP      = 66
	FCTYPE_UEOPT          = 67
	FCTYPE_HDVIDEO        = 68
	FCTYPE_METRICS        = 69
	FCTYPE_OFFERCAM       = 70
	FCTYPE_REQUESTCAM     = 71
	FCTYPE_MYWEBCAM       = 72
	FCTYPE_MYCAMSTATE     = 73
	FCTYPE_PMHISTORY      = 74
	FCTYPE_CHATFLASH      = 75
	FCTYPE_TRUEPVT        = 76
	FCTYPE_BOOKMARKS      = 77
	FCTYPE_EVENT          = 78
	FCTYPE_STATEDUMP      = 79
	FCTYPE_RECOMMEND      = 80
	FCTYPE_EXTDATA        = 81
	FCTYPE_NOTIFY         = 84
	FCTYPE_PUBLISH        = 85
	FCTYPE_ZGWINVALID     = 95
	FCTYPE_CONNECTING     = 96
	FCTYPE_CONNECTED      = 97
	FCTYPE_DISCONNECTED   = 98
	FCTYPE_LOGOUT         = 99
)

type wsCommand struct {
	Command IncomingCommand
	FromId  int64
	ToId    int64
	Arg1    int64
	Arg2    int64
	Payload string
}

type rxData struct {
	Msglen  int `json:"msglen"`
	Opts    int `json:"opts"`
	Respkey int `json:"respkey"`
	Serv    int `json:"serv"`
	Type    int `json:"type"`
}

type ServerConfig struct {
	Release          bool              `json:"release"`
	AjaxServers      []string          `json:"ajax_servers"`
	ChatServers      []string          `json:"chat_servers"`
	VideoServers     []string          `json:"video_servers"`
	H5VideoServers   map[string]string `json:"h5_video_servers"`
	WebsocketServers map[string]string `json:"websocket_servers"`
}
