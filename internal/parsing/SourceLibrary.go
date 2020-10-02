package parsing

import (
	"github.com/aleferri/casmeleon/internal/text"
	"path/filepath"
	"strings"
)

type sourceList []*text.SourceLine

func expandIfFile(name string, isFile bool) (string, error) {
	if !isFile {
		return name, nil
	}
	return filepath.Abs(name)
}

//SourceLibrary list all parsed sourced
type SourceLibrary struct {
	sourcesName []string
	sources     []sourceList
}

//NewSourceLibrary return a new SourceLibrary struct
func NewSourceLibrary() *SourceLibrary {
	return &SourceLibrary{}
}

//IndexOf return the index of the source or -1
func (p *SourceLibrary) IndexOf(sourceName string) int {
	for i, s := range p.sourcesName {
		if strings.EqualFold(s, sourceName) {
			return i
		}
	}
	return -1
}

func (p *SourceLibrary) createNewEntry(name string, source []*text.SourceLine) {
	p.sourcesName = append(p.sourcesName, name)
	p.sources = append(p.sources, source)
}

//AddSource add a source line to the specified source
//Return the inserted name and an error if occurred
func (p *SourceLibrary) AddSource(sourceName string, isFile bool, source []*text.SourceLine) (string, error) {
	fullName, err := expandIfFile(sourceName, isFile)
	if err != nil {
		return sourceName, err
	}
	index := p.IndexOf(fullName)
	if index == -1 {
		p.createNewEntry(fullName, source)
	} else {
		p.sources[index] = source
	}
	return fullName, nil
}

//GetSource return the full source of the sourceName
func (p *SourceLibrary) GetSource(fullName string) []*text.SourceLine {
	for i, s := range p.sourcesName {
		if strings.EqualFold(s, fullName) {
			return p.sources[i]
		}
	}
	return nil
}
