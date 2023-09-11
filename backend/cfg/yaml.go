package cfg

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

func Parse(s string, out interface{}) error {
	f, err := os.ReadFile(s)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(f, out)
	if err != nil {
		log.Fatalf("UNMARSHAL WASTED: %v", err)
	}
	log.Tracef("\n\n\n\n\nUNMARSHALLED: %v\n\n", out)
	return err
}

func Save(name string, in interface{}) {
	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	b, err := yaml.Marshal(in)
	if err != nil {
		log.Fatalf("MARSHAL WASTED: %v", err)
	}
	_, err = f.Write(b)
	if err != nil {
		log.Errorf("write yaml (e): %v", err)
	}
	log.Tracef("MARSHALLED: %v\n\n", f)
}

func mostRecentModifiedYAML(dirs ...string) string {
	last := time.Time{}
	res := ""
	for _, d := range dirs {
		dir, e := os.ReadDir(d)
		if e != nil {
			log.Errorf("\nerr:%v\nduring run:%v", e, "lookout")
		}
		for _, entry := range dir {
			if !entry.IsDir() && filepath.Ext(entry.Name()) == ".yaml" {
				i, _ := entry.Info()
				if i.ModTime().After(last) {
					last = i.ModTime()
					res = filepath.Join(d, i.Name())
				}
			}
		}

	}
	return res
}
