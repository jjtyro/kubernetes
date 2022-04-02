/*
Wrap try-catch-finally interface. so that program not panic
*/

package cache_version

type Block struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

type Exception interface{}

func Throw(up Exception) {
	//exist another func to replace panic ???
	panic(up)
}

func (tcf Block) Do() {
	if tcf.Finally != nil {
		defer tcf.Finally()
	}
	if tcf.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				tcf.Catch(r)
			}
		}()
	}
	tcf.Try()
}

var ExceptionTriggered = false

func IsExceptioned() bool {
	return ExceptionTriggered
}