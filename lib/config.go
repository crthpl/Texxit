package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)


type Config struct {
	BreakThresh int32	`json:"breakThresh"`
	WalkThresh int16	`json:"walkThresh"`
	WalkSpeed int16		`json:"walkSpeed"`
}

type CraftConfig struct {
	Crafts []CraftRecipe`json:"crafts"`
}

type CraftRecipe struct {
	Reqs []ItemStack	`json:"reqs"`
	Result ItemStack	`json:"result"`
}

func ReadJson(file string) (Config) {
	//open the file
	jsonFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	//read + decrypt the file


	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result Config
	json.Unmarshal(byteValue, &result)

	return result
}

func ReadCrafts(file string) (CraftConfig) {
	//open the file
	jsonFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	//read + decrypt the file


	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result CraftConfig
	json.Unmarshal(byteValue, &result)
	fmt.Println(string(byteValue))
	return result
}

