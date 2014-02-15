package reserve

import (
	"bytes"
	"encoding/hex"
	"log"
	"os"
	"path"
	"strconv"
	"text/template"
	"time"
)

var (
	DateFormat     = "20060102T1504"
	FilenameFormat = "{{.Start}}_{{.Type}}{{.Ch}}_{{.Duration}}_{{.Title}}.ts"
	SaveDirectory  = "."
	filenameTmpl   *template.Template
)

func init() {
	ChangeFilename(FilenameFormat)
}

func ChangeFilename(tmpl string) {
	filenameTmpl = template.Must(template.New("filename").Parse(tmpl))
}

type FileRecord struct {
	channel    *Channel
	program    *Program
	recordInfo *RecordInfo
	*os.File
}

func NewFileRecord(channel Channel, program Program) *FileRecord {
	return &FileRecord{
		channel: &channel,
		program: &program,
		recordInfo: &RecordInfo{
			Id:    hex.EncodeToString(program.Hash),
			Type:  channel.Type,
			Ch:    channel.Ch,
			Sid:   strconv.Itoa(channel.Sid),
			Start: program.Start,
			End:   program.End,
		},
	}
}

func (f *FileRecord) Info() *RecordInfo {
	return f.recordInfo
}

func (f *FileRecord) Write(p []byte) (int, error) {
	if f.File == nil {
		file, err := os.OpenFile(f.getFilename(), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Print(err)
			return 0, err
		}
		f.File = file
	}
	return f.File.Write(p)
}

func (f *FileRecord) getFilename() string {
	variables := struct {
		Type     string
		Ch       string
		Start    string
		End      string
		Duration int
		Title    string
	}{
		Type:     f.channel.Type,
		Ch:       f.channel.Ch,
		Start:    time.Unix(int64(f.program.Start), 0).Format(DateFormat),
		End:      time.Unix(int64(f.program.End), 0).Format(DateFormat),
		Duration: f.program.End - f.program.Start,
		Title:    f.program.Title,
	}

	buf := bytes.NewBuffer(nil)
	filenameTmpl.Execute(buf, variables)
	return path.Join(SaveDirectory, buf.String())
}
