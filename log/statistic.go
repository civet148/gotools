package log

/*
  统计每个文件的方法被调用的总次数、总调用时间、平均执行时间（毫秒）
*/

type Statistic struct {
	SticMap map[string]interface{}
}

//create a new statistic object
func NewStatistic() *Statistic {
	return &Statistic{
		SticMap: make(map[string]interface{}, 1),
	}
}

//进入方法(enter function)
func (s *Statistic) Enter() {

}

//退出方法(leave function)
func (s *Statistic) Leave() {

}

//统计信息汇总(statistic summary)
func (s *Statistic) summary(args...interface{}) (strSummary string) {

	return
}
