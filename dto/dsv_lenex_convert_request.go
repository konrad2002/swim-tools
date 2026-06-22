package dto

import "swim-tools/enum"

type DsvLenexConvertRequest struct {
	Filename string        `json:"filename"`
	Origin   enum.FileType `json:"origin"`
	Target   enum.FileType `json:"target"`
}
