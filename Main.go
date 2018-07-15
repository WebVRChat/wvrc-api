package main

import (
    "log"
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

var app Application

func main() {
    // Parse configuration
    configFile, err := ioutil.ReadFile("config.yaml")
    if err != nil {
        log.Fatal(err)
    }
    config := GlobalConfiguration{}
    yaml.Unmarshal(configFile, &config)

    // Initialize & run API
    var app Application
    app.Initialize(config)
    app.Run()
}