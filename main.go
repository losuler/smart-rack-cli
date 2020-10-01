package main

import (
    "fmt"
    "log"
    "bufio"
    "os"
    "os/exec"
    "os/signal"
    "strings"
    "math/rand"
    "syscall"
    "regexp"
    "errors"
    "strconv"
    "time"
    "runtime"
    "golang.org/x/crypto/ssh/terminal"
    "github.com/gocolly/colly/v2"
    "github.com/briandowns/spinner"
    "github.com/cheynewallace/tabby"
)

const (
    // Disables printing date/time in logs
    LstdFlags = ""
)

type Login struct {
    Username string
    Password string
    RoomURL string
    BookedKit Kit
    BookedDuration string
}

type Kit struct {
    Name string
    Type string
    CurrentUser string
    BookURL string
}

type Device struct {
    Name string
    Model string
    VANPort string
    Power string
    Network string
    Number string
    // SSH login details
    Username string
    Password string
    Domain string
}

func getSession(c *colly.Collector, login Login) *colly.Collector {
    var sessionID string

    client := c.Clone()

    client.OnHTML("input[name='session_id']", func(e *colly.HTMLElement) {
        sessionID = e.Attr("value")
        err := client.Post(login.RoomURL, map[string]string {
                "action": "LOGIN",
                "username": login.Username,
                "password": login.Password,
                "session_id": sessionID,
        })
        if err != nil {
            log.Fatalln("Please check your login details.")
        }
    })

    client.OnError(func(r *colly.Response, err error) {
        log.Fatalln("There was an issue connecting to ", login.RoomURL)
    })

    client.Visit(login.RoomURL)

    return client
}

func getKits(c *colly.Collector, login Login) ([]Kit, error) {
    var kits []Kit

    // Clone retains cookies but excludes OnHTML (e.g. from getSession())
    client := c.Clone()

    client.OnHTML("body", func(e *colly.HTMLElement) {
        e.ForEach("table.esp tbody tr", func(_ int, e *colly.HTMLElement) {
            var kit Kit
            
            kit.Name = e.ChildText("td:nth-child(1)")
            kit.Type = e.ChildText("td:nth-child(2)")
            kit.CurrentUser = e.ChildText("td:nth-child(3)")
            kit.BookURL = e.ChildAttr("a:nth-child(1)", "href")

            kits = append(kits, kit)
        })
    })

    client.Visit(login.RoomURL + "/?action=SHOW_KITS")

    if len(kits) == 0 {
        err := errors.New("No kits were found.")
        return kits, err
    }
    
    return kits, nil
}

func pickKit(kits []Kit) (Kit, error) {
    var freeKits []Kit
    var freeKit Kit

    for _, kit := range kits {
        if kit.CurrentUser == "" && (kit.Type == "CCNP" || kit.Type == "CCNA") {
            freeKits = append(freeKits, kit)
        }
    }

    if len(freeKits) == 0 {
        err := errors.New("No free kits were found.")
        return freeKit, err
    }    

    randIndex := rand.Intn(len(freeKits))
    freeKit = freeKits[randIndex]

    return freeKit, nil
}

func getKit(c *colly.Collector, login Login) map[string]string {
    var req map[string]string
    client := c.Clone()

    client.OnHTML("body", func(e *colly.HTMLElement) {
        startTime := e.ChildAttr("select[name='start_time'] option", "value")

        data := map[string]string {
            "action": "REQUEST_BOOKING",
            "duration": login.BookedDuration,
            "start_time": startTime,
        }

        e.ForEach("input[type='checkbox']", func(_ int, e *colly.HTMLElement) {
            data[e.Attr("name")] = "true"
        })

        req = data
    })
    
    client.Visit(login.RoomURL + login.BookedKit.BookURL)

    return req
}

func bookKit(c *colly.Collector, login Login, req map[string]string) {
    client := c.Clone()

    err := client.Post(login.RoomURL + login.BookedKit.BookURL, req)
    if err != nil {
        log.Fatalln("Booking kit failed.")
    }

    // Needs to sleep as it takes a second to process the booking
    time.Sleep(5 * time.Second)
}

func getDevices(c *colly.Collector, login Login) ([]Device, error) {
    var devices []Device

    client := c.Clone()

    client.OnHTML("body", func(e *colly.HTMLElement) {
        e.ForEach("table.esp tbody tr", func(_ int, e *colly.HTMLElement) {
            var device Device

            device.Name = e.ChildText("td:nth-child(1)")
            device.Model = e.ChildText("td:nth-child(2)")
            device.VANPort = e.ChildText("td:nth-child(4)")
            device.Power = e.ChildText("td:nth-child(5)")
            device.Network = e.ChildText("td:nth-child(6)")
            device.Number = e.ChildAttr("input[type='checkbox']", "name")

            ssh := e.ChildAttr("td:nth-child(7) a:nth-child(4)", "href")
            sshRegex := regexp.MustCompile(`^ssh://(.+);password=(.+)@(.+)/$`)
            sshRegexGroups := sshRegex.FindStringSubmatch(ssh)

            device.Username = sshRegexGroups[1]
            device.Password = sshRegexGroups[2]
            device.Domain = sshRegexGroups[3]

            devices = append(devices, device)
        })
    })

    client.Visit(login.RoomURL)

    if len(devices) == 0 {
        err := errors.New("No Devices were found.")
        return devices, err
    }

    return devices, nil
}

func pickDevice(devices []Device, alreadyBooked bool) (device Device, skip bool) {
    // The kit name gets printed on first booking
    if alreadyBooked {
        // Use first device as they all have the same name
        kitDevice := strings.Split(devices[0].Name, " ")

        // This prints the same as Kit.Name
        // TODO: Can this be done better?
        if len(kitDevice) == 8 {
            fmt.Printf("Name: %s %s %s %s %s %s\n", kitDevice[0], kitDevice[1], kitDevice[2],
                                                    kitDevice[3], kitDevice[4], kitDevice[5])
        } else if len(kitDevice) == 7 {
            fmt.Printf("Name: %s %s %s %s %s\n", kitDevice[0], kitDevice[1], kitDevice[2],
                                                    kitDevice[3], kitDevice[4])
        }
    }

    table := tabby.New()
    for i, d := range devices {
        // Number to string
        num := strconv.Itoa(i + 1)

        // Print only the device type (router/switch) and number
        kitDevice := strings.Split(d.Name, " ")
        deviceName := fmt.Sprintf("%s %s", kitDevice[6], kitDevice[7])

        table.AddLine(num + ": " + deviceName, d.Model, d.VANPort, d.Power, d.Network)
    }
    table.Print()

    tries := 1
    for {
        msg := fmt.Sprintf("Choose a device to connect to from 1-%d (0 will skip): ",
                           len(devices))
        resp := readInput(msg)

        // Needs to be int to compare to len()
        respNum, _ := strconv.Atoi(resp)
        // Since this is + 1 in the prompt
        respNum -= 1

        if respNum <= len(devices) && respNum >= 1 {
            skip = false
            return devices[respNum], skip
        } else if respNum < 1 {
            skip = true
            return devices[0], skip
        } else if tries > 3 {
            log.Fatalln("A device in the range was not chosen.") 
        }
    }
}

func powerOnDevice(c *colly.Collector, login Login, device Device) error {
    client := c.Clone()

    data := map[string]string {
        "action": "POWER_ON",
    }

    data[device.Number] = "true"

    err := client.Post(login.RoomURL, data)
    if err != nil {
        return err
    }

    time.Sleep(5 * time.Second)

    return nil
}

func getBookedDevices(c *colly.Collector, login Login) map[string]string {
    var req map[string]string
    client := c.Clone()

    client.OnHTML("body", func(e *colly.HTMLElement) {
        data := map[string]string {
            "action": "DO_RELEASE_BOOKING",
        }

        e.ForEach("input[type='checkbox']", func(_ int, e *colly.HTMLElement) {
            data[e.Attr("name")] = "true"
        })

        req = data
    })
    
    client.Visit(login.RoomURL + "/?action=RELEASE_BOOKING")

    return req  
}

func powerOffDevices(c *colly.Collector, login Login, devices []Device) error {
    client := c.Clone()

    data := map[string]string {
        "action": "POWER_OFF",
    }

    poweredOn := 0
    for _, device := range devices {
        if device.Power == "On" {
            data[device.Number] = "true"
            poweredOn += 1
        }  
    }
    
    if poweredOn != 0 {
       err := client.Post(login.RoomURL, data)
        if err != nil {
            return err
        }
    }

    time.Sleep(5 * time.Second)

    return nil
}

func releaseDevices(c *colly.Collector, login Login, bookedDevices map[string]string) error {
    client := c.Clone()

    err := client.Post(login.RoomURL, bookedDevices)
    if err != nil {
        return err
    }

    return nil
}

func readInput(msg string) string {
    tries := 0
    for {
        var input string

        fmt.Print(msg)
        scanner := bufio.NewScanner(os.Stdin)
        scanner.Scan()
        input = scanner.Text()

        tries += 1
        if input != "" {
            return input
        } else if tries > 3 {
            // TODO: Fix doesn't print exitMsg()
            log.Fatalln("No input was provided.") 
        }

        fmt.Println("Please check your input.")
    }
}

func readPass(msg string) string {
    tries := 0
    for {
        var pass string
        
        fmt.Print(msg)
        bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
        fmt.Print("\n")
        pass = string(bytePassword)

        tries += 1
        if pass != "" {
            return pass
        } else if tries > 3 { 
            log.Fatalln("No input was provided.")
        }
        
        fmt.Println("Please check your input.")
    }
}

func readInputDuration() string {
    tries := 0
    for {
        bookingTimes := regexp.MustCompile(`^(30|90|120|150|180|210|240)$`)
        
        input := readInput("Book for 30, 90, 120, 150, 180, 210 or 240 minutes: ")

        tries += 1
        if bookingTimes.MatchString(input) {
            return input
        } else if tries > 3 {
            log.Fatalln("No input was provided.")
        }
    }
}

func readInputRoom() string {
    tries := 0
    for {
        var roomURL string
    
        input := readInput("ATC328, ATC329 or ATC330: ")

        switch strings.ToUpper(input) {
        case "ATC328":
            roomURL = "https://ictencsvr2.ict.swin.edu.au"
        case "ATC329":
            roomURL = "https://ictencsvr10.ict.swin.edu.au"
        case "ATC330":
            roomURL = "https://ictencsvr12.ict.swin.edu.au"
        default:
            roomURL = ""
            fmt.Println("Please check you entered a room name.")
            tries += 1
        }

        if roomURL != "" {
            return roomURL
        } else if tries > 3 {
            log.Fatalln("Correct room name was not provided.")
        }
    }
}

func runSSH(device Device) error {
    host := fmt.Sprintf("%s@%s", device.Username, device.Domain)
    system := runtime.GOOS

    switch system {
        case "linux", "darwin":
            cmd := exec.Command("sshpass", "-p", device.Password, "ssh", 
                                "-o", "StrictHostKeyChecking=no", host)
    
            cmd.Stdout = os.Stdout
        	cmd.Stderr = os.Stderr
            cmd.Stdin = os.Stdin
    
            if err := cmd.Run(); err != nil {
                return err
            }
        case "windows":
            cmd := exec.Command("putty.exe", "-ssh", host, "-pw", device.Password)
                
            cmd.Stdout = os.Stdout
        	cmd.Stderr = os.Stderr
            cmd.Stdin = os.Stdin
    
            if err := cmd.Run(); err != nil {
                return err
            }
        default:
            err := errors.New("Operating system not identified.")
            return err
    }

    return nil
}

func alreadyBooked(user string, kits []Kit) bool {
    // Kit selection appends an "s" to student IDs
    user = "s" + user

    for _, kit := range kits {
        if kit.CurrentUser == user {
            return true
        }
    }
    return false
}

func exitMsg(login Login) {
    fmt.Println("Please shutdown and release the devices manually.")
    fmt.Println(login.RoomURL)
}

func main() {
    var login Login

    log.SetFlags(0)
    client := colly.NewCollector(
        colly.AllowURLRevisit(),
    )

    login.RoomURL = readInputRoom()

    login.Username = readInput("Student ID: ")
    login.Password = readPass("Password: ")

    spin := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
    spin.Suffix = " Logging in..."
    spin.Start()

    client = getSession(client, login)

    spin.Suffix = " Getting kits..."

    kits, err := getKits(client, login)
    if err != nil {
        log.Fatalln(err)
    }

    // True if a kit is actively booked
    alreadyBooked := alreadyBooked(login.Username, kits)

    if !alreadyBooked {
        freeKit, err := pickKit(kits)
        if err != nil {
            log.Fatalln(err)
        }

        spin.FinalMSG = "✔ Getting kits complete.\n"
        spin.Stop()

        fmt.Printf("Name: %s\nType: %s\n", freeKit.Name, freeKit.Type)

        resp := readInput("Would you like to continue with booking this kit? (y/n) ")

        if strings.ToLower(resp) == "y" {
            login.BookedKit = freeKit

            login.BookedDuration = readInputDuration()

            spin.Suffix = " Booking kit..."
            spin.Start()

            req := getKit(client, login)
            bookKit(client, login, req)

            spin.Suffix = " Getting devices..."
        } else {
            os.Exit(0)
        }
    } else {
        spin.Suffix = " Getting devices..."
    }

    // Once booked, capture Ctrl+C and print exit message
    go func() {
        sigchan := make(chan os.Signal)
        signal.Notify(sigchan, os.Interrupt)
        <-sigchan
        exitMsg(login)
        os.Exit(0)
    }()

    devices, err := getDevices(client, login)
    if err != nil {
        spin.FinalMSG = "✔ Getting devices complete.\n"
        spin.Stop()
        exitMsg(login)
        log.Fatalln(err)
    }

    spin.FinalMSG = "✔ Getting devices complete.\n"
    spin.Stop()

    pickedDevice, skip := pickDevice(devices, alreadyBooked)

    if !skip {
        spin.Suffix = " Powering on device..."
        spin.Start()

        powerOnDevice(client, login, pickedDevice)

        spin.FinalMSG = "✔ Powering on device complete.\n"
        spin.Stop()

        fmt.Println(pickedDevice.Name)

        if err := runSSH(pickedDevice); err != nil {
            tries := 0
            for {
                resp := readInput("The session appears to have disconnected. Retry? (y/n) ")

                if resp == "y" {
                    runSSH(pickedDevice)
                    tries += 1
                } else if resp == "n" {
                    break
                } else if tries > 3 {
                    log.Println("Retries exceeded.")
                    break
                }
            }
        }
    }

    resp := readInput("Do you want to shutdown and release all devices now? (y/n) ")

    if resp == "y" {
        spin.Suffix = " Powering off devices..."
        spin.Start()

        // Get latest device status
        devices, err = getDevices(client, login)
        if err != nil {
            spin.Stop()
            exitMsg(login)
            log.Fatalln(err)
        }

        err := powerOffDevices(client, login, devices)
        if err != nil {
            spin.Stop()
            exitMsg(login)
            log.Fatalln("Unable to poweroff devices.")
        }

        bookedDevices := getBookedDevices(client, login)

        spin.Suffix = " Releasing devices..."

        err = releaseDevices(client, login, bookedDevices)
        if err != nil {
            spin.Stop()
            exitMsg(login)
            log.Fatalln("Unable to release devices.")
        }

        spin.FinalMSG = "✔ Releasing devices complete.\n"
        spin.Stop()
    } else {
        log.Println("User aborted.")
        exitMsg(login)
        os.Exit(0)
    }
}
