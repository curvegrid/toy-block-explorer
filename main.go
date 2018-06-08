package main

import (
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"html/template"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"path/filepath"
	"time"
)

// GetBlockchainInfo retrieve top-level information about the blockchain
func GetBlockchainInfo() *BlockchainInfo {
	blockchainInfo := &BlockchainInfo{}

	// RPC call to retrieve the latest block

	// translate from string (hex probably) to *big.Int

	// retrieve last 10 blocks, ensuring not to attempt to access an invalid block
	maxBlock := 10

	for i := 0; i < maxBlock; i++ {
		// retrieve the block, which includes all of the transactions

		// store the block info in a struct

		// append the block info to the blockchain info struct
	}

	return blockchainInfo
}

// ShortHex returns a shortened version of a hex string
func ShortHex(long string) string {
	if len(long) < 19 {
		return long
	}

	return long[0:8] + "..." + long[len(long)-8:]
}

// HandleTemplates is an HTTP handler for templated files
func HandleTemplates(next http.Handler) http.Handler {
	templatedExtension := ".html"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request from %v: %v", r.RemoteAddr, r.URL.Path)

		// check if file extension matches templated extension
		if filepath.Ext(r.URL.Path) == templatedExtension || r.URL.Path == "/" {
			// populate template data
			blockchainInfo := GetBlockchainInfo()

			if blockchainInfo == nil {
				log.Println("Unable to retrieve blockchain info")
				return
			}

			// parse any form and query parameters
			if err := r.ParseForm(); err != nil {
				log.Println("Error parsing form parameters")
			}

			// requested filename on disk
			diskFilename := filepath.Base(r.URL.Path)
			if r.URL.Path == "/" {
				diskFilename = "index.html"
			}

			// pull out just the filename and add to the www root directory
			diskFilepath := filepath.Join(Options.WWWRoot, diskFilename)

			// clone the template components
			newTemplates, err := Templates.Clone()
			if err != nil {
				log.Printf("Unable to clone templates: %v", err.Error())
				return
			}

			// read the requested file
			templateData, err := ioutil.ReadFile(diskFilepath)
			if err != nil {
				log.Printf("Unable to read the requested file: %v", err.Error())
				return
			}

			// custom template functions
			funcMap := template.FuncMap{
				"ShortHex": ShortHex,
			}

			// add the top requested file
			newTemplates, err = newTemplates.New("main").Funcs(funcMap).Parse(string(templateData))
			if err != nil {
				log.Printf("Unable to parse template for requested file: %v", err.Error())
				return
			}

			// execute and write the composed template to the HTTP response writer
			newTemplates.Execute(w, blockchainInfo)
			if err != nil {
				log.Printf("Unable to execute template for requested file: %v", err.Error())
				return
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}

// InitTemplates reads in the HTML template files for later use
func InitTemplates() {
	templateFiles, err := filepath.Glob(Options.TemplatesGlob)
	if err != nil {
		log.Fatal(err)
	}

	Templates = template.Must(template.New("base").Parse(""))

	for _, templateFile := range templateFiles {
		templateData, err := ioutil.ReadFile(templateFile)
		if err != nil {
			log.Fatal(err)
		}

		Templates = template.Must(Templates.New(filepath.Base(templateFile)).Parse(string(templateData)))
	}
}

type BlockchainInfo struct {
	LastBlockNum *big.Int
	ThisBlockNum *big.Int
	Blocks       []BlockInfo
}

type BlockInfo struct {
	Num              *big.Int
	Timestamp        time.Time
	Hash             string
	TransactionCount int
	Miner            string
}

// OptionsStruct contains global program options
type OptionsStruct struct {
	Host          string
	Port          int
	WWWRoot       string
	TemplatesGlob string
	EthEndpoint   string
}

var Options OptionsStruct

// templates
var Templates *template.Template

// blockchain client connection
var Client *ethclient.Client
var RPCClient *rpc.Client

// cached blockchain data
var MaxBlockNum int64

func main() {
	// command line options
	flag.StringVar(&Options.Host, "host", "", "Hostname to bind web server to")
	flag.IntVar(&Options.Port, "port", 8080, "Port to bind web server to")
	flag.StringVar(&Options.WWWRoot, "www", "www", "Directory to serve")
	flag.StringVar(&Options.TemplatesGlob, "templates", "templates/*", "Templates glob")
	flag.StringVar(&Options.EthEndpoint, "ethendpoint", "http://localhost:8545", "Ethereum node endpoint")
	flag.Parse()

	// setup templates
	InitTemplates()

	log.Printf("Connecting to Ethereum node at %v", Options.EthEndpoint)

	// connect to RPC via the Eth client

	// connect to RPC via the RPC client

	// start web server
	http.Handle("/", HandleTemplates(http.FileServer(http.Dir(Options.WWWRoot))))

	bind := fmt.Sprintf("%v:%d", Options.Host, Options.Port)
	log.Printf("Web server started on %v", bind)

	log.Fatal(http.ListenAndServe(bind, nil))
}
