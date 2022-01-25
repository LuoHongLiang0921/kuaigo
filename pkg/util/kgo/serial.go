// @Description 串行执行迭代器

package kgo

import (
	"go.uber.org/multierr"
)

// SerialWithError
// 	@Description 创建串行迭代器，返回所有函数错误，不中断执行
//	@Param fns 执行函数列表
// 	@return func() error 新的迭代器
func SerialWithError(fns ...func() error) func() error {
	return func() error {
		var errs error
		for _, fn := range fns {
			errs = multierr.Append(errs, try(fn, nil))
		}
		return errs
	}
}

// SerialUntilError
// 	@Description: 创建串行迭代器，返回所有函数错误，中断执行，返回第一个错误
//	@Param fns 执行函数列表
// 	@return func() error 新的迭代器
func SerialUntilError(fns ...func() error) func() error {
	return func() error {
		for _, fn := range fns {
			if err := try(fn, nil); err != nil {
				return err
				// return errors.Wrap(err, xstring.FunctionName(fn))
			}
		}
		return nil
	}
}

// WhenError 策略注入
type WhenError int

var (

	// ReturnWhenError ...
	ReturnWhenError WhenError = 1

	// ContinueWhenError ...
	ContinueWhenError WhenError = 2

	// PanicWhenError ...
	PanicWhenError WhenError = 3

	// LastErrorWhenError ...
	LastErrorWhenError WhenError = 4
)

// SerialWhenError
// 	@Description 创建串行执行迭代器，并带有错误处理逻辑
//	@Param we 错误类型
// 	@return func(fn ...func() error) func() error
func SerialWhenError(we WhenError) func(fn ...func() error) func() error {
	return func(fns ...func() error) func() error {
		return func() error {
			var errs error
			for _, fn := range fns {
				if err := try(fn, nil); err != nil {
					switch we {
					case ReturnWhenError: // 直接退出
						return err
					case ContinueWhenError: // 继续执行
						errs = multierr.Append(errs, err)
					case PanicWhenError: // panic
						panic(err)
					case LastErrorWhenError: // 返回最后一个错误
						errs = err
					}
				}
			}
			return errs
		}
	}
}
