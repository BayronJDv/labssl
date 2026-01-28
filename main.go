package main

import (
    "log"
    // Importa el paquete UI
    // Reemplaza "tu-modulo" con el nombre que está en tu go.mod
    "github.com/BayronJDv/labssl/bubbletea" 

    tea "github.com/charmbracelet/bubbletea"
)

func main() {
    // Llamamos a ui.InitialModel()
    p := tea.NewProgram(bubbletea.InitialModel())
    if _, err := p.Run(); err != nil {
        log.Fatal(err)
    }
}