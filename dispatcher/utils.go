package dispatcher

import "fmt"

func diffTimeStrapToShow(start int64, end int64) string {
	second := (end - start) / (1000 * 1000 * 1000)
	mSecond := (end - start - second*1000*1000*1000) / (1000 * 1000)
	uSecond := (end - start - second*1000*1000*1000 - mSecond*1000*1000) / 1000
	nSecond := end - start - second*1000*1000*1000 - mSecond*1000*1000 - uSecond*1000

	return fmt.Sprintf("%ds,%dms,%dus,%dns", second, mSecond, uSecond, nSecond)
}
