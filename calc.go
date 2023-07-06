package fare

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
)

func CalcFare(dst io.Writer, src io.Reader) error {
	var (
		r = csv.NewReader(src)
		w = csv.NewWriter(dst)
	)
	r.ReuseRecord = true
	defer w.Flush()

	// 设置title
	if err := setFareTitle(w, r); err != nil {
		if errors.Is(err, io.EOF) {
			logger.Warnf("read empty file")
			return nil
		}
		return fmt.Errorf("添加'费用'列失败, %s", err)
	}

	// 计算运费
	var sumFare int64
	for {
		var (
			record, err = r.Read()
			line, _     = r.FieldPos(0)
		)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("读取数据错误, 行%d附近, %s", line, err)
		}

		if record[8] == "" {
			continue
		}
		weight, err := strconv.ParseFloat(record[8], 64)
		if err != nil {
			return fmt.Errorf("无效的结算重量, 行: %d", line)
		}

		fare := calcFare(record[7], weight)
		sumFare += fare
		record = append(record, strconv.FormatInt(fare, 10))
		logger.Debugf("append record, len: %d, cap: %d", len(record), cap(record))

		if err = w.Write(record); err != nil {
			return fmt.Errorf("写入数据失败, %s", err)
		}
	}

	lastRecord := make([]string, r.FieldsPerRecord+1)
	lastRecord[r.FieldsPerRecord] = strconv.FormatInt(sumFare, 10)
	if err := w.Write(lastRecord); err != nil {
		return fmt.Errorf("写入数据失败, %s", err)
	}
	return nil
}

func setFareTitle(w *csv.Writer, r *csv.Reader) error {
	record, err := r.Read()
	if err != nil {
		return err
	}

	record = append(record, "费用")
	return w.Write(record)
}

// 某地区运费单价
type areaFare struct {
	// 首重价格
	startingFare int64
	// 续重价格
	secondFare int64
}

var (
	defaultAreaFare = areaFare{
		startingFare: 5,
		secondFare:   2,
	}
	areaFareMap = map[string]areaFare{
		"黑龙江省": {
			startingFare: 7,
			secondFare:   4,
		},
		"吉林省": {
			startingFare: 7,
			secondFare:   4,
		},
		"辽宁省": {
			startingFare: 7,
			secondFare:   4,
		},
		"内蒙古自治区": {
			startingFare: 7,
			secondFare:   4,
		},
		"甘肃省": {
			startingFare: 7,
			secondFare:   4,
		},
		"宁夏回族自治区": {
			startingFare: 7,
			secondFare:   4,
		},
		"青海省": {
			startingFare: 7,
			secondFare:   4,
		},
		"新疆维吾尔自治区": {
			startingFare: 12,
			secondFare:   12,
		},
		"海南省": {
			startingFare: 12,
			secondFare:   12,
		},
		"西藏自治区": {
			startingFare: 12,
			secondFare:   12,
		},
	}
)

func getAreaFare(area string) areaFare {
	af, ok := areaFareMap[area]
	if ok {
		return af
	}
	return defaultAreaFare
}

func calcFare(area string, weight float64) int64 {
	if weight <= 0 {
		return 0
	}
	weight = math.Ceil(weight)
	af := getAreaFare(area)
	return af.startingFare + int64(weight-1)*af.secondFare
}
