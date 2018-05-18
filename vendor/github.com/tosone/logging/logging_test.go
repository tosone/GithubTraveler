package logging

import (
	"crypto/rand"
	"testing"
)

func TestInfo(t *testing.T) {
	Info("info level")
	WithFields(Fields{"test": "test"}).Info("info level")
}

func TestDebug(t *testing.T) {
	Debug("debug level")
	WithFields(Fields{"test": "test"}).Debug("info level")
}

func TestWarn(t *testing.T) {
	Warn("warn level")
	WithFields(Fields{"test": "test"}).Warn("info level")
}

func TestError(t *testing.T) {
	Error("error level")
	WithFields(Fields{"test": "test"}).Warn("error level")
}

func TestFatal(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	Fatal("fatal level")
	WithFields(Fields{"test": "test"}).Fatal("info level")
}

func TestUnknownLevel(t *testing.T) {
	var level uint = 10
	t.Log(Level(level).String())
}

func TestPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	Panic("panic level")
	WithFields(Fields{"test": "test"}).Panic("info level")
}

func TestLevel_String(t *testing.T) {
	var level Level
	level = DebugLevel
	t.Log(level.String())
	level = InfoLevel
	t.Log(level.String())
	level = WarnLevel
	t.Log(level.String())
	level = ErrorLevel
	t.Log(level.String())
	level = FatalLevel
	t.Log(level.String())
	level = PanicLevel
	t.Log(level.String())
}

func TestFileCannotBeWrite(t *testing.T) {
	var gb = 1024 * 1024
	var maxSize = 1
	var conf = Config{
		LogLevel:   DebugLevel,
		Filename:   "test.log",
		MaxSize:    maxSize,
		MaxBackups: 2,
		MaxAge:     30,
		LocalTime:  true,
		Compress:   true,
	}
	Setting(conf)
	var arr = make([]byte, gb*maxSize)
	rand.Read(arr)
	Info(arr)
}

func TestSetting(t *testing.T) {
	var conf = Config{
		LogLevel:   DebugLevel,
		Filename:   "test.log",
		MaxSize:    10,
		MaxBackups: 2,
		MaxAge:     30,
		LocalTime:  true,
		Compress:   true,
	}
	Setting(conf)
	Info("after setting info level")
	WithFields(Fields{"test": "test"}).Info("after setting info level")
}

func TestLevelLower(t *testing.T) {
	var conf = Config{
		LogLevel:   InfoLevel,
		Filename:   "test.log",
		MaxSize:    10,
		MaxBackups: 2,
		MaxAge:     30,
		LocalTime:  true,
		Compress:   true,
	}
	Setting(conf)
	Debug("after setting info level")
	WithFields(Fields{"test": "test"}).Debug("after setting info level")
}

func TestRotate(t *testing.T) {
	var conf = Config{
		LogLevel:   DebugLevel,
		Filename:   "test.log",
		MaxSize:    10,
		MaxBackups: 2,
		MaxAge:     30,
		LocalTime:  true,
		Compress:   true,
	}
	Setting(conf)
	Rotate()
	Info("after setting info level")
	WithFields(Fields{"test": "test"}).Info("after setting info level")
}

func BenchmarkInfo(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Info("info level")
	}
}

func BenchmarkInfoFile(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Info("write info level to file")
	}
}
