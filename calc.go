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

		fare := calcFare(weight)
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

const (
	// 首重价格5rmb
	startingFare = 5
	// 首重1kg
	startingWeight = 1
	// 续重价格2rmb
	secondFare = 2
)

func calcFare(weight float64) int64 {
	if weight <= 0 {
		return 0
	}
	weight = math.Ceil(weight)
	return int64(startingFare + (weight-startingWeight)*secondFare)
}
