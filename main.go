package main

import (
    "fmt"
    "net/http"
    "github.com/scrapli/scrapligo/driver/options"
    "github.com/scrapli/scrapligo/platform"
    "github.com/scrapli/scrapligo/channel"
)

func runScrapligo(ip, cmd string, vendoros string) (string, error) {
    p, err := platform.NewPlatform(
        vendoros,
        ip,
        options.WithAuthNoStrictKey(),
        options.WithAuthUsername("admin"),
        options.WithAuthPassword("admin"),
    )
    if err != nil {
        return "", fmt.Errorf("failed to create platform; error: %+v", err)
    }

    d, err := p.GetNetworkDriver()
    if err != nil {
        return "", fmt.Errorf("failed to fetch network driver from the platform; error: %+v", err)
    }

    err = d.Open()
    if err != nil {
        return "", fmt.Errorf("failed to open driver; error: %+v", err)
    }

	defer d.Close()

	// fetch the prompt
	prompt, err := d.Channel.GetPrompt()
	if err != nil {
		fmt.Printf("failed to get prompt; error: %+v\n", err)

		return"",err
	}

	fmt.Printf("found prompt: %s\n\n\n", prompt)

	// send some input
	output, err := d.Channel.SendInput(cmd)
	if err != nil {
		fmt.Printf("failed to send input to device; error: %+v\n", err)

		return"",err
	}

	fmt.Printf("output received (SendInput):\n %s\n\n\n", output)

	// send an interactive input
	// SendInteractive expects a slice of `SendInteractiveEvent` objects
	events := make([]*channel.SendInteractiveEvent, 2)
	events[0] = &channel.SendInteractiveEvent{
		ChannelInput:    "clear logging",
		ChannelResponse: "[confirm]",
		HideInput:       false,
	}
	events[1] = &channel.SendInteractiveEvent{
		ChannelInput:    "",
		ChannelResponse: "#",
		HideInput:       false,
	}

	interactiveOutput, err := d.SendInteractive(events)
	if err != nil {
		fmt.Printf("failed to send interactive input to device; error: %+v\n", err)
	}
	if interactiveOutput.Failed != nil {
		fmt.Printf("response object indicates failure: %+v\n", interactiveOutput.Failed)

		return"",err
	}

	fmt.Printf("output received (SendInteractive):\n %s\n\n\n", interactiveOutput.Result)

	// send a command -- as this is a driver created from a *platform* it will have some things
	// already done for us -- including disabling paging, so this command that would produce more
	// output than the default terminal lines will not cause any issues.
	r, err := d.SendCommand(cmd)
	if err != nil {
		fmt.Printf("failed to send command; error: %+v\n", err)
		return"",err
	}
	if r.Failed != nil {
		fmt.Printf("response object indicates failure: %+v\n", r.Failed)

		return"",err
	}

	fmt.Printf(
		"sent command '%s', output received (SendCommand):\n %s\n\n\n",
		r.Input,
		r.Result,
	)
	return r.Result, nil
}

func formHandler(w http.ResponseWriter, r *http.Request) {
    ip := r.FormValue("ip")
    cmd := r.FormValue("command")
	vendoros := r.FormValue("platform")
    output, err := runScrapligo(ip, cmd, vendoros)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "Output: %s", output)
}

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
    http.HandleFunc("/form", formHandler)
    http.ListenAndServe(":8080", nil)
	fmt.Printf("starting server at port 8080\n")
}
