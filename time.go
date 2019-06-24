package gnet

import "time"

//GetDate ex:2019-01-14 22:39:50 +0800 CST
func GetDate() time.Time {
	t1 := time.Now().Year()
	t2 := time.Now().Month()
	t3 := time.Now().Day()
	t4 := time.Now().Hour()
	t5 := time.Now().Minute()
	t6 := time.Now().Second()
	//t7 := time.Now().Nanosecond()
	return time.Date(t1, t2, t3, t4, t5, t6, 0, time.Local)
}
