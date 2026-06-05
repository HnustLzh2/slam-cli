package authflow

import (
	"os"

	qrterminal "github.com/mdp/qrterminal/v3"
)

func PrintLoginQR(content string) {
	cfg := qrterminal.Config{
		Level:          qrterminal.L,
		Writer:         os.Stdout,
		HalfBlocks:     true,
		BlackChar:      qrterminal.BLACK_BLACK,
		WhiteBlackChar: qrterminal.WHITE_BLACK,
		WhiteChar:      qrterminal.WHITE_WHITE,
		BlackWhiteChar: qrterminal.BLACK_WHITE,
		QuietZone:      1,
	}
	qrterminal.GenerateWithConfig(content, cfg)
}
