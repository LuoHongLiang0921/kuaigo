package console

import (
	"fmt"
	"sync"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap/buffer"
)

var _sliceEncoderPool = sync.Pool{
	New: func() interface{} {
		return &sliceArrayEncoder{elems: make([]interface{}, 0, 2)}
	},
}

func getSliceEncoder() *sliceArrayEncoder {
	return _sliceEncoderPool.Get().(*sliceArrayEncoder)
}

func putSliceEncoder(e *sliceArrayEncoder) {
	e.elems = e.elems[:0]
	_sliceEncoderPool.Put(e)
}

type consoleEncoder struct {
	*jsonEncoder
}

//NewConsoleEncoder 创建了一个编码器，其输出是为人类而非机器消费而设计的。
//它以纯文本格式序列化核心日志条目数据（消息、级别、时间戳等），
//并将结构化上下文保留为 JSON。请注意，尽管控制台编码器不使用编码器配置中指定的键，但它会忽略任何键设置为空字符串的元素。
func NewConsoleEncoder(cfg EncoderConfig) zapcore.Encoder {
	if len(cfg.ConsoleSeparator) == 0 {
		// Use a default delimiter of '\t' for backwards compatibility
		cfg.ConsoleSeparator = "\t"
	}
	return consoleEncoder{newJSONEncoder(cfg, true)}
}

func (c consoleEncoder) Clone() zapcore.Encoder {
	return consoleEncoder{c.jsonEncoder.Clone().(*jsonEncoder)}
}

func (c consoleEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	line := _pool.Get()

	// We don't want the entry's metadata to be quoted and escaped (if it's
	// encoded as strings), which means that we can't use the JSON encoder. The
	// simplest option is to use the memory encoder and fmt.Fprint.
	//
	// If this ever becomes a performance bottleneck, we can implement
	// ArrayEncoder for our plain-text format.
	arr := getSliceEncoder()
	if c.TimeKey != "" && c.EncodeTime != nil {
		c.EncodeTime(ent.Time, arr)
	}
	if c.LevelKey != "" && c.EncodeLevel != nil {
		c.EncodeLevel(ent.Level, arr)
	}
	if ent.LoggerName != "" && c.NameKey != "" {
		nameEncoder := c.EncodeName

		if nameEncoder == nil {
			// Fall back to FullNameEncoder for backward compatibility.
			nameEncoder = zapcore.FullNameEncoder
		}

		nameEncoder(ent.LoggerName, arr)
	}
	if ent.Caller.Defined {
		if c.CallerKey != "" && c.EncodeCaller != nil {
			c.EncodeCaller(ent.Caller, arr)
		}
		if c.FunctionKey != "" {
			arr.AppendString(ent.Caller.Function)
		}
	}
	for i := range arr.elems {
		if i > 0 {
			line.AppendString(c.ConsoleSeparator)
		}
		fmt.Fprint(line, arr.elems[i])
	}
	putSliceEncoder(arr)

	// Add the message itself.
	if c.MessageKey != "" {
		c.addSeparatorIfNecessary(line)
		line.AppendString(ent.Message)
	}

	fer := NewFormatter(Format(c.Format))
	format := fer.String(fields...)
	if format != "" {
		line.AppendString(c.ConsoleSeparator)
	}
	line.AppendString(format)
	// Add any structured context. 去掉
	//c.writeContext(line, fields)

	// If there's no stacktrace key, honor that; this allows users to force
	// single-line output.
	if ent.Stack != "" && c.StacktraceKey != "" {
		line.AppendByte('\n')
		line.AppendString(ent.Stack)
	}

	if c.LineEnding != "" {
		line.AppendString(c.LineEnding)
	} else {
		line.AppendString(zapcore.DefaultLineEnding)
	}
	return line, nil
}

func (c consoleEncoder) writeContext(line *buffer.Buffer, extra []zapcore.Field) {
	context := c.jsonEncoder.Clone().(*jsonEncoder)
	defer func() {
		// putJSONEncoder assumes the buffer is still used, but we write out the buffer so
		// we can free it.
		context.buf.Free()
		putJSONEncoder(context)
	}()

	addFields(context, extra)
	context.closeOpenNamespaces()
	if context.buf.Len() == 0 {
		return
	}

	c.addSeparatorIfNecessary(line)
	line.AppendByte('{')
	line.Write(context.buf.Bytes())
	line.AppendByte('}')
}

func (c consoleEncoder) addSeparatorIfNecessary(line *buffer.Buffer) {
	if line.Len() > 0 {
		line.AppendString(c.ConsoleSeparator)
	}
}
