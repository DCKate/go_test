package kkutil

import (
    "fmt"
    "os"
    "io"
    "log"
    "bufio"
    "strings"
    "os/exec"
    "encoding/json"
)

// type DevInfo struct {
//     Name              string      `json:"uuname"`
//     Account           string      `json:"username"`
//     Password          string      `json:"password"`
// }

// type AddDev struct {
//     Uid               string      `json:"uuid"`
//     DeviceType        int         `json:"type"`
//     DeviceInfo        DevInfo     `json:"info"`
// }

func ReadLine(filename string) {
    f, err := os.Open(filename)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer f.Close()
    r := bufio.NewReaderSize(f, 4*1024)
    line, isPrefix, err := r.ReadLine()
    for err == nil && !isPrefix {
        s := string(line)
        fmt.Println(s)
        line, isPrefix, err = r.ReadLine()
    }
    if isPrefix {
        fmt.Println("buffer size to small")
        return
    }
    if err != io.EOF {
        fmt.Println(err)
        return
    }
}

func ReadString(filename string)[]string {
    rstrs := make([]string, 0)
    f, err := os.Open(filename)
    if err != nil {
        fmt.Println(err)
        return rstrs
    }
    defer f.Close()
    r := bufio.NewReader(f)
    line, err := r.ReadString('\n')
    for err == nil {
        rstrs = append(rstrs,strings.Split(line,"\n")[0])
        line, err = r.ReadString('\n')
        // fmt.Println(strings.Split(line,"\n")[0])
    }
    if err != io.EOF {
        fmt.Println(err)
        return rstrs
    }
    rstrs = append(rstrs,strings.Split(line,"\n")[0])
    return rstrs
}

func AsyncCmd(program string, args ...string)*exec.Cmd {
  cmd := exec.Command(program,args...)
  err := cmd.Start()
  if err != nil {
      log.Fatal(err)
  }
  return cmd
}

func ExecCmd(program string, args ...string) {
  // Start starts the specified command but does not wait for it to complete.
  // The Wait method will return the exit code and release associated resources once the command exits.
    cmd := exec.Command(program, args...)
    // cmd.Stdin = os.Stdin;
    cmd.Stdout = os.Stdout;
    // cmd.Stderr = os.Stderr;
    
    err := cmd.Run() 
    if err != nil {
        log.Printf("%v\n", err)
    }
    //err = cmd.Wait()
    log.Printf("Command finished with error: %v", err)
}

func GenJsonByte(mm map[string]string)[]byte {
    jm, _ := json.Marshal(mm)
    fmt.Println(string(jm))
    return []byte(jm)

    // dinfo:=DevInfo{
    //     Name:uid,
    //     Account:acc,
    //     Password:pwd,
    // }
    // add_req:=AddDev{
    //     Uid:uid,
    //     DeviceType:1,
    //     DeviceInfo:dinfo,
    // }

    // add_payload, err := json.Marshal(add_req)
    // if err != nil {
    //     fmt.Println("error:", err)
    // }
}