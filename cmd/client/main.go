package main

import (
	"leveling/internal/client/ui"
)

func main() {
	newUi := *ui.NewUi()
	newUi.Run("Brian")
}
