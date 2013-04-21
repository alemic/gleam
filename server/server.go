package main

import (
    "os"
    "flag"
    "time"
    "github.com/mikespook/golib/log"
    "github.com/mikespook/z-node/node"
    "github.com/mikespook/golib/signal"
)

const (
    SCRIPT_ROOT = "Z_NODE_SCRIPT_ROOT"
)

var (
    uri = flag.String("doozer", "doozer:?ca=127.0.0.1:8046", "address of the doozerd")
    buri = flag.String("dzns", "", "address of the DzNS")
    region = flag.String("region", "z-node", "a region of the z-node located in")
    scriptPath = flag.String("script", "", "default script path(as the enviroment variable $Z_NODE_SCRIPT_ROOT)")
)

func init() {
    if !flag.Parsed() {
        flag.Parse()
    }
    node.ErrHandler = func(err error) {
        log.Error(err)
    }

    if *scriptPath == "" {
        *scriptPath = os.Getenv(SCRIPT_ROOT)
    }
    if *scriptPath == "" {
        var err error
        *scriptPath, err = os.Getwd()
        if err != nil {
            log.Error(err)
        }
    }
    node.SetDefaultPath(*scriptPath)
}

func main() {
    defer func() {
        log.Message("Exit.")
        time.Sleep(time.Second)
    }()
    log.Message("Starting...")
    hostname, err := os.Hostname()
    if err != nil {
        log.Error(err)
        return
    }
    n := node.New(*region, hostname)
    n.ErrHandler = node.ErrHandler
    n.Bind("Stop", node.Stop)
    n.Bind("Restart", node.Restart)
    n.Bind("Python", node.Python)
    if err := n.Start(*uri, *buri); err != nil {
        log.Error(err)
        return
    }
    defer n.Close()
    go func() {
        n.Wait()
        signal.Send(os.Getpid(), os.Interrupt)
    }()
    // signal handler
    sh := signal.NewHandler()
    sh.Bind(os.Interrupt, func() bool {return true})
    sh.Loop()
}
