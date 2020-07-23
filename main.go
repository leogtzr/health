package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/marcusolsson/tui-go"
)

var (
	configFile = flag.String("config", "", "configuration file")
)

func main() {

	var sb strings.Builder
	for i := 0; i < 100; i++ {
		sb.WriteString(fmt.Sprintf("Txt: %d\n", i))
	}

	flag.Parse()

	if len(*configFile) == 0 {
		panic(fmt.Errorf("config option missing"))
	}

	executablePath, err := binaryPath()
	if err != nil {
		panic(err)
	}

	envConfig, err := readConfig(*configFile, executablePath, map[string]interface{}{
		MaxIdleConnections:        20,
		RequestTimeout:            7,
		InstancesFileName:         "instances.csv",
		PrintStatusEverySeconds:   3,
		CheckServicesEverySeconds: 5,
	})
	if err != nil {
		panic(err)
	}

	httpClient := createHTTPClient(envConfig)
	httpsClient := createHTTPSClient(envConfig)

	config := Configuration{
		envConfig:   envConfig,
		httpClient:  httpClient,
		httpsClient: httpsClient,
	}

	logFile, err := os.OpenFile("calls.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	file, err := os.Open(filepath.Join(executablePath, envConfig.GetString("instancesFileName")))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	services, err := readInstancesFile(file)
	if err != nil {
		panic(err)
	}

	txtArea := tui.NewTextEdit()
	txtArea.SetSizePolicy(tui.Expanding, tui.Expanding)
	txtArea.SetText("")
	txtArea.SetFocused(true)
	txtArea.SetWordWrap(true)
	txtAreaScroll := tui.NewScrollArea(txtArea)
	// txtAreaScroll.SetAutoscrollToBottom(true)

	txtAreaBox := tui.NewVBox(txtAreaScroll)
	txtAreaBox.SetBorder(true)

	inputCommand := newInputCommandEntry()
	inputCommandBox := newInputCommandBox(inputCommand)

	txtReader := tui.NewVBox(txtAreaBox, inputCommandBox)
	txtReader.SetSizePolicy(tui.Expanding, tui.Expanding)

	root := tui.NewHBox(txtReader)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	addCloseApplicationKeyBinding(ui)

	go updateScreen(&services, &config, txtArea, ui)

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}

}
