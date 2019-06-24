package gnet

var Common *CommonST
var IsInit bool
var DB *DBST
var Format *FormatST
var File *FileST
var Log *LogST
var Config *ConfigST
var Math *mathST
var Sys *sysST
var Date *dateST
var Web *webST
var Service *serviceST

func Init() bool {
	IsInit = false
	Common = NewCommon()
	if Common == nil {
		return false
	}

	DB = NewDB()
	if DB == nil {
		return false
	}

	Format = NewFormat()
	if Format == nil {
		return false
	}

	File = NewFile()
	if File == nil {
		return false
	}

	Log = NewLog()
	if Log == nil {
		return false
	}

	Config = NewConfig()
	if Config == nil {
		return false
	}

	Math = NewMath()
	if Math == nil {
		return false
	}

	Sys = NewSys()
	if Sys == nil {
		return false
	}

	Date = NewDate()
	if Date == nil {
		return false
	}

	Web = NewWeb()
	if Web == nil {
		return false
	}

	Service = NewService()
	if Service == nil {
		return false
	}

	return true
}

func Start(handle IFIoservice, fps int) {
	Service.SetHandle(handle)
	Service.Run(fps)
}
