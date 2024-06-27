package config

const ErrColor = "\033[0;31m"
const WarningColor = "\033[0;43m"
const OkColor = "\033[0;32m"
const NoColor = "\033[0m"

const CommandSize = 128
const CmdFiller = "-"

const SRC = "src"
const SrcTest = "src_test"
const SrcWeb = "src_web"
const SrcSelenium = "src_selenium"
const SrcUpgrade = "src_upgrade"
const SrcHybrid = "src_hybrid"

const SrcSymbol = "s"
const SrcWebSymbol = "w"
const SrcTestSymbol = "t"
const ClobberSymbol = "c"
const SeleniumSymbol = "f"
const UpgradeSymbol = "u"
const HybridSymbol = "h"

var SrcAliases = map[string]string{
	SrcSymbol:      SRC,
	SrcTestSymbol:  SrcTest,
	SrcWebSymbol:   SrcWeb,
	SeleniumSymbol: SrcSelenium,
	UpgradeSymbol:  SrcUpgrade,
	HybridSymbol:   SrcHybrid,
}
var ClobberableSources = map[string]string{
	SrcSymbol:      SRC,
	SrcTestSymbol:  SrcTest,
	SeleniumSymbol: SrcSelenium,
	UpgradeSymbol:  SrcUpgrade,
	HybridSymbol:   SrcHybrid,
}
var NonClobberableSources = map[string]string{
	SrcWebSymbol: SrcWeb,
}
var TestSources = []string{SRC, SrcTest, SrcWeb}

const BuildCommandFormat = "ant -f %s/%s/build.xml"
const ClobberCommandFormat = "ant clobber -f %s/%s/build.xml"
const TestCommandFormat = "ant %s -f %s/%s/build.xml"
const SpecificTestCommandFormat = " -Dtest.includes=**/%s"
const NumKeyBuildCommandFormat = `ant -v -f /opt/wnc/tools_vs/build/commonUtils.xml darjeeling.start_dbserver
ant -f %s/%s/build.xml clean clobber all -Ddarjeeling.updnumkey=true
ant -v -f /opt/wnc/tools_vs/build/commonUtils.xml darjeeling.stop_dbserver`
