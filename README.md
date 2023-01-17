# Command-Runner

This is simple learning exercise to use a golang web api to serve up a form to run commands against an IOS-XE device.

Scrapligo is used to execute the CLI command with the output returned back to the browser.

Note this is purely for learning purposes in terms of integrating the API with Scrapligo and is not intended for a Production use case

## How to use this repo ##

1) Clone the repo
2) Ensure you have installed go and the associated scrapligo dependencies
3) `go build` 
4) `go run main.go` 
5) Fire up your `localost:8080/form.html` on your browser which is served from the static folder
6) Enter the IP address and command as per screenshot below


![Screenshot from 2023-01-17 20-20-53](https://user-images.githubusercontent.com/63735312/213004449-2eeceb82-f30a-4a5e-a074-4ce2ad3975a8.png)
