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

const SrcSymbol = "s"
const SrcWebSymbol = "w"
const SrcTestSymbol = "t"
const ClobberSymbol = "c"

var SrcAliases = map[string]string{SrcSymbol: SRC, SrcTestSymbol: SrcTest, SrcWebSymbol: SrcWeb}
var TestSources = []string{SRC, SrcTest, SrcWeb}

const BuildCommandFormat = "ant -f %s/%s/build.xml"
const ClobberCommandFormat = "ant clobber -f %s/%s/build.xml"
const TestCommandFormat = "ant %s -f %s/%s/build.xml"
const SpecificTestCommandFormat = " -Dtest.includes=**/%s"
