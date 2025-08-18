package log

import (
	"flag"

	zzap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

const (
	DebugLevel = int(zapcore.DebugLevel)
)

// Options stores controller-runtime (zap) log config
var Options = &zap.Options{
	Development: true,
	// we dont log with panic level, so this essentially
	// disables stacktrace, for now, it avoids un-needed clutter in logs
	StacktraceLevel: zapcore.DPanicLevel,
	TimeEncoder:     zapcore.RFC3339NanoTimeEncoder,
	Level:           zzap.NewAtomicLevelAt(zapcore.InfoLevel),
	// log caller (file and line number) in "caller" key
	EncoderConfigOptions: []zap.EncoderConfigOption{func(ec *zapcore.EncoderConfig) { ec.CallerKey = "caller" }},
	ZapOpts:              []zzap.Option{zzap.AddCaller()},
}

// BindFlags binds controller-runtime logging flags to provided flag Set
func BindFlags(fs *flag.FlagSet) {
	Options.BindFlags(fs)
}

// InitLog initializes controller-runtime log (zap log)
// this should be called once Options have been initialized
// either by parsing flags or directly modifying Options.
func InitLog() {
	log.SetLogger(zap.New(zap.UseFlagOptions(Options)))
}

// SetLogLevel sets current logging level to the provided lvl
func SetLogLevel(lvl string) error {
	newLevel, err := zapcore.ParseLevel(lvl)
	if err != nil {
		return err
	}

	currLevel := Options.Level.(zzap.AtomicLevel).Level()

	if newLevel != currLevel {
		log.Log.Info("set log verbose level", "newLevel", newLevel.String(), "currentLevel", currLevel.String())
		Options.Level.(zzap.AtomicLevel).SetLevel(newLevel)
	}
	return nil
}

// GetLogLevel returns the current logging level
func GetLogLevel() string {
	return Options.Level.(zzap.AtomicLevel).Level().String()
}
