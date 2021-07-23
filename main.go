package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const (
	base       = `\\?\GLOBALROOT\Device\HarddiskVolumeShadowCopy`
	timeFormat = "2006-01-02T15_04_05Z07_00"
)

func processFile(path string) ([]byte, time.Time, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, time.Now(), fmt.Errorf("error opening file: %+v", err)
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil, time.Now(), fmt.Errorf("error getting file info: %+v", err)
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, time.Now(), fmt.Errorf("error reading file content: %+v", err)
	}
	return content, info.ModTime(), nil
}

func checkFile(path string) ([]byte, time.Time) {
	var lastmodify time.Time
	var content []byte
	for i := 1; i <= 20; i++ {
		fullPath := fmt.Sprintf(`%s%d\%s`, base, i, path)
		fileContent, fileMod, err := processFile(fullPath)
		if err != nil {
			// fmt.Println(err)
			continue
		}
		if fileMod.After(lastmodify) {
			lastmodify = fileMod
			content = fileContent
		}
	}
	return content, lastmodify
}

func main() {
	content, lastMod := checkFile(`Windows\System32\config\SAM`)
	filename := fmt.Sprintf("hive_sam_%s", lastMod.Format(timeFormat))
	if err := ioutil.WriteFile(filename, content, 0644); err != nil {
		fmt.Printf("could not write %s: %v\n", filename, err)
	}
	fmt.Printf("Saved a copy of SAM to %s with last modify date of %s\n", filename, lastMod)

	content, lastMod = checkFile(`Windows\System32\config\SECURITY`)
	filename = fmt.Sprintf("hive_security_%s", lastMod.Format(timeFormat))
	if err := ioutil.WriteFile(filename, content, 0644); err != nil {
		fmt.Printf("could not write %s: %v\n", filename, err)
	}
	fmt.Printf("Saved a copy of SECURITY to %s with last modify date of %s\n", filename, lastMod)

	content, lastMod = checkFile(`Windows\System32\config\SYSTEM`)
	filename = fmt.Sprintf("hive_system_%s", lastMod.Format(timeFormat))
	if err := ioutil.WriteFile(filename, content, 0644); err != nil {
		fmt.Printf("could not write %s: %v\n", filename, err)
	}
	fmt.Printf("Saved a copy of SYSTEM to %s with last modify date of %s\n", filename, lastMod)
}
