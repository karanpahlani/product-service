package logging

import(
	"log/slog"
	"os"
)

func New() *slog.Logger{
	return *slog.New(slog.NewJsonHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelInfo,
})) }