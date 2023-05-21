package util

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

type PerfTimer struct {
	label   string
	started time.Time
}

func StartTimer(label string) *PerfTimer {
	return &PerfTimer{
		label:   label,
		started: time.Now(),
	}
}

func (t *PerfTimer) End() {
	elapsed := time.Since(t.started)
	fmt.Printf("%s: %v\n", t.label, elapsed)
}
