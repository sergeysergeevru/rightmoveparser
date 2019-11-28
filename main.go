package main

import (
	"encoding/json"
	"fmt"
	"github.com/sergeysergeevru/rightmoveparser/types"
	"io/ioutil"
	"net/http"
)

const URL = "https://www.rightmove.co.uk/api/_search?locationIdentifier=REGION%5E87490&maxBedrooms=3&numberOfPropertiesPerPage=24&radius=0.0&sortType=6&index=%d&includeLetAgreed=false&viewType=LIST&furnishTypes=furnished&channel=RENT&areaSizeUnit=sqft&currencyCode=GBP&isFetching=false"
const STEP  = 24


func main()  {
	index := 0
	list := make(map[string]int64)
	minPrice := make(map[string]int)
	maxPrice := make(map[string]int)

	uniqueStorage := make(map[uint32]bool)

	count := make(map[string]int64)
	var orderIndex []string
	var avrCommon float64
	for {
		fmt.Println("Offset ", index, " page is ", index/STEP)
		response, err := http.Get(fmt.Sprint(URL, index))
		if err != nil {
			panic(err)
		}
		if response.StatusCode != 200 {
			fmt.Println("Status code is not 200: ", response.StatusCode)
			break
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}
		file := &types.RightMoveResponse{}
		err = json.Unmarshal(data, file)
		if err != nil {
			panic(err)
		}


		for _, v := range file.Properties {
			_, wasIt := uniqueStorage[v.Id]
			if wasIt {
				continue
			}
			uniqueStorage[v.Id] = true
			avrCommon += float64(v.Price.GetMonthPrice())
			key := fmt.Sprintf("property_%d", v.Bedrooms)
			if v.IsShare() {
				key = "share"
			}
			_, isInArray := list[key]
			if !isInArray {
				orderIndex = append(orderIndex, key)
				count[key] = 1
				minPrice[key] = v.Price.GetMonthPrice()
				maxPrice[key] = v.Price.GetMonthPrice()
			} else {
				count[key] ++
				if minPrice[key] > v.Price.GetMonthPrice() {
					minPrice[key] = v.Price.GetMonthPrice()
				}
				if maxPrice[key] < v.Price.GetMonthPrice() {
					maxPrice[key] = v.Price.GetMonthPrice()
				}
			}
			list[key] += int64(v.Price.GetMonthPrice())

			response.Body.Close()
		}
		fmt.Println(avrCommon/float64(len(uniqueStorage)), " len ", len(uniqueStorage))
		for _, key := range orderIndex {
			fmt.Printf("%s price is %d count %d min is %d max is %d \n",
				key,
				list[key]/count[key],
				count[key],
				minPrice[key],
				maxPrice[key],
				)
		}
		index += STEP
	}



}