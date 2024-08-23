package util

import (
	"ch-common-package/logger"
	util2 "ch-common-package/util"
	"context"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func getLastHeightTimestamp() int64 {
	cnt := time.Now().Unix() / 30
	return cnt * 30
}

func GetDateList(lastTimestamp int64) []string {
	var dates []string

	for i := 29; i >= 0; i-- {
		last := lastTimestamp - int64(DaySeconds*i)
		if last >= GenesisTimestamp {
			dates = append(dates, time.Unix(last, 0).Format(GolangDayFormat))
		}
	}

	return dates
}

func GetTimestampList(lastTimestamp int64) []int64 {
	var timestamps []int64

	for i := 30; i > 0; i-- {
		last := lastTimestamp - int64(DaySeconds*i)
		if last >= GenesisTimestamp {
			timestamps = append(timestamps, last)
		}
	}

	return timestamps
}

func Offset(pageSize int, page int) int64 {
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return int64((page - 1) * pageSize)
}

func GinPostPagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func RegJsonData(s string) string {
	reg := regexp.MustCompile(`[\w]+`)
	s = reg.ReplaceAllStringFunc(s, func(k string) string {
		return "\"" + k + "\""
	})
	return s
}

func ReplaceLineFeed(s string) string {
	s = strings.ReplaceAll(string(s), "\\n", "")
	s = strings.ReplaceAll(s, "\"", "")
	return s
}

// 从program源码中提取出program id
func GetProgramId(code string) string {
	if code == "" {
		return ""
	}

	comma := strings.Index(code, "program")
	if comma == -1 {
		return ""
	}

	if len(code) > comma {
		b := code[comma:]
		single := strings.Split(b, ";")
		if len(single) > 0 {
			res := strings.TrimLeft(single[0], "program")
			if strings.Contains(res, ".aleo") {
				return strings.TrimSpace(res)
			}
		}
	}

	return ""
}

func GetLanguage(ctx context.Context) string {
	v := ctx.Value(LanguageKey)
	if v == nil {
		return Language_CN
	}

	language, ok := v.(string)
	if ok {
		return language
	}

	return Language_CN
}

func SaveLanguageToContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader(LanguageKey)
		if key == "" {
			key = c.Query(LanguageKey)
		}

		ctx := context.WithValue(context.Background(), LanguageKey, key)
		c.Set(LanguageKey, key)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// 返回的值需要退6位
func GetFloatInAleoNum(num string, unit string) float64 {
	str := strings.TrimSuffix(num, unit)
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		logger.Errorf("GetFloatInAleoNum ParseFloat | %v", err)
		return 0
	}
	return value
}

func ExitsInArrayInt(num int, array []int) bool {
	exist := false
	for _, v := range array {
		if v == num {
			exist = true
		}
	}
	return exist
}

func GenerateInviteCode(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}

func GetDateListByStart(start, genesis int64) []string {
	var res []string
	for i := start; i < start+DaySeconds*30; i += DaySeconds {
		//从创世区块开始计算
		if i < genesis {
			continue
		}
		res = append(res, time.Unix(i, 0).Format(GolangDateFormat))
	}
	return res
}

func GetDateListByTimeRange(start, end int64) []string {
	var res []string
	for i := start; i <= end; i += DaySeconds {
		res = append(res, time.Unix(i, 0).Format(GolangDateFormat))
	}
	return res
}

func GetTimestampListByStart(start, genesis int64) []int64 {
	var res []int64
	for i := start; i < start+DaySeconds*30; i += DaySeconds {
		if i < genesis {
			continue
		}
		res = append(res, i)
	}
	return res
}

func GetTimestampListByTimeRange(start, end, genesis int64) []int64 {
	var res []int64
	for i := start; i <= end; i += DaySeconds {
		if i < genesis {
			continue
		}
		res = append(res, i)
	}
	return res
}

func GetOneMonthTimeRange() (int64, int64) {
	end := util2.GetNullPoint(time.Now().Unix()) - DaySeconds
	start := end - 30*DaySeconds
	return start, end
}

func GetOneMonthUTCTimeRange() (int64, int64) {
	end := util2.GetNullUTCPoint(time.Now().Unix()) - DaySeconds
	start := end - 30*DaySeconds
	return start, end
}

func GetOneMonthAgoTime() int64 {
	return util2.GetNullPoint(time.Now().Unix()) - 30*DaySeconds
}

func GetOneMonthAgoUTCTimeStart() int64 {
	return util2.GetNullUTCPoint(time.Now().Unix()) - 30*DaySeconds
}

// 获取指定个数的随机字符串
// @param            size            int         指定的字符串个数
// @param            kind            int         0，纯数字；1，小写字母；2，大写字母；3，数字+大小写字母
// @return           string                      返回生成的随机字符串
func RandStr(size int, kind int) string {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	isAll := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return string(result)
}
