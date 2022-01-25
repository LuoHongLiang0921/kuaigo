package console

import (
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/zap"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap/buffer"
)

var (
	_pool = buffer.NewPool()
	// Get retrieves a buffer from the pool, creating one if necessary.
	Get = _pool.Get
)

type FormatEncoder struct {
	*zapcore.MapObjectEncoder
}

func NewFormatEncoder() zap.Encoder {
	return FormatEncoder{
		MapObjectEncoder: zapcore.NewMapObjectEncoder(),
	}
}

func (enc FormatEncoder) Clone() zapcore.Encoder {
	return FormatEncoder{
		MapObjectEncoder: zapcore.NewMapObjectEncoder(),
	}
}

func (enc FormatEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	myEnc := enc.Clone().(FormatEncoder)
	buf := _pool.Get()

	buf.AppendString(entry.Message)
	buf.AppendString(" ")

	for _, field := range fields {
		field.AddTo(myEnc)
		value := myEnc.MapObjectEncoder.Fields[field.Key]
		buf.AppendString(field.Key)
		buf.AppendString("=")
		if value == "" {
			buf.AppendString(" ''")
		} else {
			buf.AppendString(fmt.Sprintf("%v ", value))
		}
	}

	buf.AppendByte('\n')

	if entry.Stack != "" {
		buf.AppendString(entry.Stack)
		buf.AppendByte('\n')
	}
	return buf, nil
}
