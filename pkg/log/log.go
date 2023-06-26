package log

import (
	"fmt"
	"time"
)

func ErrPrint(err error) {
	fmt.Printf("[LEN]-[%s] Error: %s\n", time.Now(), err)
}

func Print(log string) {
	fmt.Printf("[LEN]-[%s]  %s\n", time.Now(), log)
}
