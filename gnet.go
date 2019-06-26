package gnet

var Common *commonST
var IsInit bool
var DB *dbST
var Format *formatST
var File *fileST
var Log *logST
var Config *configST
var Math *mathST
var Sys *sysST
var Date *dateST
var Web *webST
var Service *serviceST

func Init() bool {
	IsInit = false
	Common = newCommon()
	if Common == nil {
		return false
	}

	DB = newDB()
	if DB == nil {
		return false
	}

	Format = newFormat()
	if Format == nil {
		return false
	}

	File = newFile()
	if File == nil {
		return false
	}

	Log = newLog()
	if Log == nil {
		return false
	}

	Config = newConfig()
	if Config == nil {
		return false
	}

	Math = newMath()
	if Math == nil {
		return false
	}

	Sys = newSys()
	if Sys == nil {
		return false
	}

	Date = newDate()
	if Date == nil {
		return false
	}

	Web = newWeb()
	if Web == nil {
		return false
	}

	Service = newService()
	if Service == nil {
		return false
	}

	return true
}

func Start(handle IFIoservice, fps int) {
	Service.SetHandle(handle)
	Service.run(fps)
	go func() {
		handle.OnInit()
	}()

}

func Print(format string, a ...interface{}) {
	Log.Print(format, a...)
}

func Success(format string, a ...interface{}) {
	Log.Success(format, a...)
}

func Warning(format string, a ...interface{}) {
	Log.Warning(format, a...)
}

func Error(format string, a ...interface{}) {
	Log.Error(format, a...)
}
