package service

import (
	"swim-tools/enum"
	"testing"
)

func TestConvert(t *testing.T) {
	err := Convert("../assets/sdm26.dsv7", enum.FileTypeDsv, enum.FileTypeLenex)
	if err != nil {
		return
	}
}
