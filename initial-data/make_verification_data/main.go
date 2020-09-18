package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
)

const (
	NumOfChairDetailData                = 100
	NumOfChairSearchData                = 100
	NumOfEstateDetailData               = 100
	NumOfEstateSearchData               = 100
	NumOfRecommendedEstateWithChairData = 100
	NumOfEstatesNazotteData             = 100
)

func init() {
	rand.Seed(19700101)
}

func writeSnapshotDataToFile(path string, snapshot Snapshot) {
	bytes, err := json.Marshal(snapshot)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path, bytes, os.FileMode(0777))
	if err != nil {
		panic(err)
	}
}

func main() {
	flags := flag.NewFlagSet("isucon10-qualify", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	var TargetServer string
	var DestDirectoryPath string
	var FixtureDirectoryPath string

	flags.StringVar(&TargetServer, "target-url", "http://127.0.0.1:1323", "target url")
	flags.StringVar(&FixtureDirectoryPath, "fixture-dir", "../../webapp/fixture", "fixture directory")
	flags.StringVar(&DestDirectoryPath, "dest-dir", "./result/verification_data", "destination directory")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}

	var chairSearchCondition ChairSearchCondition
	var estateSearchCondition EstateSearchCondition

	wg.Add(1)
	go func() {
		defer wg.Done()
		jsonText, err := ioutil.ReadFile(filepath.Join(FixtureDirectoryPath, "chair_condition.json"))
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(jsonText, &chairSearchCondition)
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		jsonText, err := ioutil.ReadFile(filepath.Join(FixtureDirectoryPath, "estate_condition.json"))
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(jsonText, &estateSearchCondition)
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()

	MkdirIfNotExists(DestDirectoryPath)

	// chair detail
	MkdirIfNotExists(filepath.Join(DestDirectoryPath, "chair_detail"))
	for i := 0; i < NumOfChairDetailData; i++ {
		wg.Add(1)
		go func(id int) {
			req := Request{
				Method:   "GET",
				Resource: fmt.Sprintf("/api/chair/%d", id),
				Query:    "",
				Body:     "",
			}

			snapshot := getSnapshotFromRequest(TargetServer, req)

			filename := fmt.Sprintf("%d.json", id)
			writeSnapshotDataToFile(filepath.Join(DestDirectoryPath, "chair_detail", filename), snapshot)
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Println("Done generating verification data of /api/chair/:id")

	// chair search condition
	MkdirIfNotExists(filepath.Join(DestDirectoryPath, "chair_search_condition"))
	wg.Add(1)
	go func() {
		req := Request{
			Method:   "GET",
			Resource: "/api/chair/search/condition",
			Query:    "",
			Body:     "",
		}

		snapshot := getSnapshotFromRequest(TargetServer, req)
		writeSnapshotDataToFile(filepath.Join(DestDirectoryPath, "chair_search_condition", "0.json"), snapshot)
		wg.Done()
	}()
	wg.Wait()
	log.Println("Done generating verification data of /api/chair/search/condition")

	// chair search
	MkdirIfNotExists(filepath.Join(DestDirectoryPath, "chair_search"))
	for i := 0; i < NumOfChairSearchData; i++ {
		wg.Add(1)
		go func(id int) {
			req := Request{
				Method:   "GET",
				Resource: "/api/chair/search",
				Query:    createRandomChairSearchQuery(chairSearchCondition).Encode(),
				Body:     "",
			}

			snapshot := getSnapshotFromRequest(TargetServer, req)

			filename := fmt.Sprintf("%d.json", id)
			writeSnapshotDataToFile(filepath.Join(DestDirectoryPath, "chair_search", filename), snapshot)
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Println("Done generating verification data of /api/chair/search")

	// estate detail
	MkdirIfNotExists(filepath.Join(DestDirectoryPath, "estate_detail"))
	for i := 0; i < NumOfEstateDetailData; i++ {
		wg.Add(1)
		go func(id int) {
			req := Request{
				Method:   "GET",
				Resource: fmt.Sprintf("/api/estate/%d", id),
				Query:    "",
				Body:     "",
			}

			snapshot := getSnapshotFromRequest(TargetServer, req)

			filename := fmt.Sprintf("%d.json", id)
			writeSnapshotDataToFile(filepath.Join(DestDirectoryPath, "estate_detail", filename), snapshot)
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Println("Done generating verification data of /api/estate/:id")

	// estate search condition
	MkdirIfNotExists(filepath.Join(DestDirectoryPath, "estate_search_condition"))
	wg.Add(1)
	go func() {
		req := Request{
			Method:   "GET",
			Resource: "/api/estate/search/condition",
			Query:    "",
			Body:     "",
		}

		snapshot := getSnapshotFromRequest(TargetServer, req)
		writeSnapshotDataToFile(filepath.Join(DestDirectoryPath, "estate_search_condition", "0.json"), snapshot)
		wg.Done()
	}()
	wg.Wait()
	log.Println("Done generating verification data of /api/estate/search/condition")

	// estate search
	MkdirIfNotExists(filepath.Join(DestDirectoryPath, "estate_search"))
	for i := 0; i < NumOfEstateSearchData; i++ {
		wg.Add(1)
		go func(id int) {
			req := Request{
				Method:   "GET",
				Resource: "/api/estate/search",
				Query:    createRandomEstateSearchQuery(estateSearchCondition).Encode(),
				Body:     "",
			}

			snapshot := getSnapshotFromRequest(TargetServer, req)
			filename := fmt.Sprintf("%d.json", id)
			writeSnapshotDataToFile(filepath.Join(DestDirectoryPath, "estate_search", filename), snapshot)
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Println("Done generating verification data of /api/estate/search")

	// chair/low_priced
	MkdirIfNotExists(filepath.Join(DestDirectoryPath, "chair_low_priced"))
	wg.Add(1)
	go func() {
		req := Request{
			Method:   "GET",
			Resource: "/api/chair/low_priced",
			Query:    "",
			Body:     "",
		}

		snapshot := getSnapshotFromRequest(TargetServer, req)
		writeSnapshotDataToFile(filepath.Join(DestDirectoryPath, "chair_low_priced", "0.json"), snapshot)
		wg.Done()
	}()
	wg.Wait()
	log.Println("Done generating verification data of /api/chair/low_priced")

	// estate/low_priced
	MkdirIfNotExists(filepath.Join(DestDirectoryPath, "estate_low_priced"))
	wg.Add(1)
	go func() {
		req := Request{
			Method:   "GET",
			Resource: "/api/estate/low_priced",
			Query:    "",
			Body:     "",
		}

		snapshot := getSnapshotFromRequest(TargetServer, req)
		writeSnapshotDataToFile(filepath.Join(DestDirectoryPath, "estate_low_priced", "0.json"), snapshot)
		wg.Done()
	}()
	wg.Wait()
	log.Println("Done generating verification data of /api/estate/low_priced")

	// recommended_estate/:id
	MkdirIfNotExists(filepath.Join(DestDirectoryPath, "recommended_estate_with_chair"))
	for i := 0; i < NumOfRecommendedEstateWithChairData; i++ {
		wg.Add(1)
		go func(id int) {
			req := Request{
				Method:   "GET",
				Resource: fmt.Sprintf("/api/recommended_estate/%d", id),
				Query:    "",
				Body:     "",
			}

			snapshot := getSnapshotFromRequest(TargetServer, req)
			fileName := fmt.Sprintf("%d.json", id)
			writeSnapshotDataToFile(filepath.Join(DestDirectoryPath, "recommended_estate_with_chair", fileName), snapshot)
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Println("Done generating verification data of /api/recommended_estate/:id")

	// estate nazotte
	MkdirIfNotExists(filepath.Join(DestDirectoryPath, "estate_nazotte"))
	for i := 0; i < NumOfEstatesNazotteData; i++ {
		wg.Add(1)
		go func(id int) {
			req := Request{
				Method:   "POST",
				Resource: "/api/estate/nazotte",
				Query:    "",
				Body:     createRandomConvexhull(),
			}

			snapshot := getSnapshotFromRequest(TargetServer, req)
			fileName := fmt.Sprintf("%d.json", id)
			writeSnapshotDataToFile(filepath.Join(DestDirectoryPath, "estate_nazotte", fileName), snapshot)
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Println("Done generating verification data of /api/estate/nazotte")
}
