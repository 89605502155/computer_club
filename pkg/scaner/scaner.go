package scaner

type Scaner interface {
	Read()
	Close()
}

type ScanerStruct struct {
	Scaner
}

func NewScaner(fileName string) *ScanerStruct {
	return &ScanerStruct{
		Scaner: NewScanConsole(fileName),
	}
}
