package sl

import "log/slog"

const (
	ModuleKey = "module"
	MethodKey = "op"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func Module(name string) slog.Attr {
	return slog.String(ModuleKey, name)
}

func Method(name string) slog.Attr {
	return slog.String(MethodKey, name)
}
