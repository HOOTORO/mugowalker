package cfg

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func Parse(s string, out interface{}) error {
	f, err := os.ReadFile(s)
	if err != nil {
		wd, _ := os.Getwd()
		str := fmt.Sprintf("Error during search %s in %s\n%v", s, wd, err)
		panic(str)
	}
	err = yaml.Unmarshal(f, out)
	if err != nil {
		panic(fmt.Sprintf("UNMARSHAL WASTED: %v", err))
	}
	fmt.Printf("\n\nUNMARSHALLED: %v\n\n", out)
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
	fmt.Printf("MARSHALLED: %v\n\n", f)
}
