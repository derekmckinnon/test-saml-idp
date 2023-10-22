package idp

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

type zerologAdapter struct {
}

func (l zerologAdapter) Printf(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}

func (l zerologAdapter) Print(v ...interface{}) {
	log.Info().Msg(fmt.Sprint(v...))
}

func (l zerologAdapter) Println(v ...interface{}) {
	log.Info().Msg(fmt.Sprint(v...))
}

func (l zerologAdapter) Fatal(v ...interface{}) {
	log.Fatal().Msg(fmt.Sprint(v...))
}

func (l zerologAdapter) Fatalf(format string, v ...interface{}) {
	log.Fatal().Msgf(format, v...)
}

func (l zerologAdapter) Fatalln(v ...interface{}) {
	log.Fatal().Msg(fmt.Sprint(v...))
}

func (l zerologAdapter) Panic(v ...interface{}) {
	log.Panic().Msg(fmt.Sprint(v...))
}

func (l zerologAdapter) Panicf(format string, v ...interface{}) {
	log.Panic().Msgf(format, v...)
}

func (l zerologAdapter) Panicln(v ...interface{}) {
	log.Panic().Msg(fmt.Sprint(v...))
}
