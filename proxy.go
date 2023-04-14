/*
 *
 * Copyright 2023 puzzleproxy authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"inet.af/tcpproxy"
)

func main() {
	if godotenv.Overload() == nil {
		log.Println("Loaded .env file")
	}

	forwardHostsStr := os.Getenv("FORWARD_HOSTS")
	ports := os.Getenv("SERVICE_PORTS")
	if forwardHostsStr == "" || ports == "" {
		log.Fatalln("Missing environment variable")
	}

	var proxy tcpproxy.Proxy
	forwardHosts := strings.Split(forwardHostsStr, ",")
	for index, port := range strings.Split(ports, ",") {
		port = cleanPort(port)
		proxy.AddRoute(port, tcpproxy.To(strings.TrimSpace(forwardHosts[index])))
	}

	bgServeHttp()

	if err := proxy.Run(); err != nil {
		log.Println("An error occurred :", err)
	}
}

func bgServeHttp() {
	responseData := []byte("Hello World !")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(responseData)
	})

	go func() {
		http.ListenAndServe(cleanPort(os.Getenv("SERVICE_PORT")), nil)
	}()
}

func cleanPort(port string) string {
	if port = strings.TrimSpace(port); port[0] != ':' {
		port = ":" + port
	}
	return port
}
