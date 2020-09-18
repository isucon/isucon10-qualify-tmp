package scenario

import (
	"fmt"
	"sort"
	"time"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
)

func checkEstateEqualToAsset(e *asset.Estate) error {
	estate, err := asset.GetEstateFromID(e.ID)
	if err != nil {
		return err
	}

	if !estate.Equal(e) {
		return fmt.Errorf("物件の情報が不正です")
	}

	return nil
}

func checkEstatesOrderedByRent(e []asset.Estate) error {
	if len(e) == 0 {
		return nil
	}

	rent := e[0].Rent
	for _, estate := range e {
		r := estate.Rent

		if rent > r {
			return fmt.Errorf("物件が賃料順に並んでいません")
		}
		rent = r
	}
	return nil
}

func checkRecommendedEstates(estates []asset.Estate, chair *asset.Chair) error {
	lengths := [3]int64{chair.Width, chair.Height, chair.Depth}
	sort.Slice(lengths[:], func(i, j int) bool { return lengths[i] < lengths[j] })
	firstMin, secondMin := lengths[0], lengths[1]

	var popularity int64 = -1
	for i, estate := range estates {
		estate, err := asset.GetEstateFromID(estate.ID)
		if err != nil {
			return err
		}

		shorterDoorLen, longerDoorLen := estate.DoorWidth, estate.DoorHeight
		if shorterDoorLen > longerDoorLen {
			shorterDoorLen, longerDoorLen = longerDoorLen, shorterDoorLen
		}
		if firstMin > shorterDoorLen || secondMin > longerDoorLen {
			return fmt.Errorf("イスがドアを通過できない物件がおすすめされています")
		}

		p := estate.GetPopularity()
		if i > 0 && popularity < p {
			return fmt.Errorf("物件がpopularity順に並んでいません")
		}
		popularity = p
	}
	return nil
}

func checkEstatesOrderedByPopularity(e []asset.Estate) error {
	var popularity int64 = -1
	for i, estate := range e {
		e, err := asset.GetEstateFromID(estate.ID)
		if err != nil {
			return err
		}
		p := e.GetPopularity()
		if i > 0 && popularity < p {
			return fmt.Errorf("物件がpopularity順に並んでいません")
		}
		popularity = p
	}
	return nil
}

func checkChairEqualToAsset(c *asset.Chair) error {
	chair, err := asset.GetChairFromID(c.ID)
	if err != nil {
		return err
	}

	if !chair.Equal(c) {
		return fmt.Errorf("イスの情報が不正です")
	}

	return nil
}

func checkChairInStock(c *asset.Chair, t time.Time) error {
	if c.GetStock() > 0 {
		return nil
	}

	soldOutTime := c.GetSoldOutTime()
	if soldOutTime == nil {
		return nil
	}

	if t.After(*soldOutTime) {
		return fmt.Errorf("イスの在庫がありません")
	}

	return nil
}

func checkChairsOrderedByPrice(c []asset.Chair, t time.Time) error {
	if len(c) == 0 {
		return nil
	}

	price := c[0].Price
	for _, chair := range c {
		_chair, err := asset.GetChairFromID(chair.ID)
		if err != nil {
			return err
		}

		err = checkChairInStock(_chair, t)
		if err != nil {
			return err
		}

		p := _chair.Price

		if price > p {
			return fmt.Errorf("イスが価格順に並んでいません")
		}
		price = p
	}

	return nil
}

func checkChairsOrderedByPopularity(c []asset.Chair, t time.Time) error {
	var popularity int64 = -1
	for i, chair := range c {
		_chair, err := asset.GetChairFromID(chair.ID)
		if err != nil {
			return err
		}

		err = checkChairInStock(_chair, t)
		if err != nil {
			return err
		}

		p := _chair.GetPopularity()

		if i > 0 && popularity < p {
			return fmt.Errorf("イスがpopularity順に並んでいません")
		}
		popularity = p
	}

	return nil
}

func checkEstatesInBoundingBox(estates []asset.Estate, boundingBox [2]point) error {
	for _, estate := range estates {
		e, err := asset.GetEstateFromID(estate.ID)
		if err != nil || !e.Equal(&estate) {
			return fmt.Errorf("物件の情報が不正です")
		}

		if !(boundingBox[0].Latitude <= e.Latitude && boundingBox[1].Latitude >= e.Latitude) {
			return fmt.Errorf("バウンディングボックス外の物件があります: BoundingBox((%v, %v), (%v, %v)), Coordinate(%v, %v)",
				boundingBox[0].Latitude, boundingBox[0].Longitude, boundingBox[1].Latitude, boundingBox[1].Latitude, e.Latitude, e.Longitude)
		}

		if !(boundingBox[0].Longitude <= e.Longitude && boundingBox[1].Longitude >= e.Longitude) {
			return fmt.Errorf("バウンディングボックス外の物件があります: BoundingBox((%v, %v), (%v, %v)), Coordinate(%v, %v)",
				boundingBox[0].Latitude, boundingBox[0].Longitude, boundingBox[1].Latitude, boundingBox[1].Latitude, e.Latitude, e.Longitude)
		}
	}

	return nil
}
