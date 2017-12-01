package kkutil

import (
    "fmt"
    "os"
    "io"
    "log"
    "bufio"
    "strings"
    "io/ioutil"
    "os/exec"
    "encoding/json"
)

type DevInfo struct {
    Name              string      `json:"name"`
    Account           string      `json:"username"`
    Password          string      `json:"password"`
}

type AddDev struct {
    Uid               string      `json:"uuid"`
    DeviceInfo        DevInfo     `json:"info"`
}
/*************************************** file read ************************************/
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

func ReadJsonFile(filename string)(int, interface{}){
    raw, err := ioutil.ReadFile(filename)
    if err != nil {
        fmt.Println(err.Error())
        return -1,nil
    }

    var jsons interface{}
    json.Unmarshal(raw,&jsons)
    fmt.Println(jsons)
    return 0,jsons
}
/*************************************** file write ************************************/
func check(e error) {
    if e != nil {
        panic(e)
    }
}

func GenFileSimpleWithByte(filename string,bb []byte) {
    err := ioutil.WriteFile(filename, bb, 0644)
    check(err)
}

func GenFileByte(filename string,bb []byte) {
    f, err := os.Create(filename)
    check(err)
    defer f.Close()//It’s idiomatic to defer a Close immediately after opening a file.
    n, err := f.Write(bb)
    check(err)
    fmt.Printf("wrote %d bytes\n", n)
    f.Sync()
}

func GenFileString (filename string,s string) {
    f, err := os.Create(filename)
    check(err)
    defer f.Close()
    n, err := f.WriteString(s)
    check(err)
    fmt.Printf("wrote %d bytes\n", n)
    f.Sync()
}

func GenFileStringbufio(filename string,s string) {
    f, err := os.Create(filename)
    check(err)
    defer f.Close()
    w := bufio.NewWriter(f)
    n, err := w.WriteString(s)
    fmt.Printf("wrote %d bytes\n", n)
    //Use Flush to ensure all buffered operations have been applied to the underlying writer.
    w.Flush()
}
/*************************************** exec call ************************************/
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
/*************************************** handle data ************************************/
func GenJsonByte(mm map[string]string)[]byte {
    jm, _ := json.Marshal(mm)
    fmt.Println(string(jm))
    return []byte(jm)
}

func FilterStrings(ss []string, filterfun func(string) bool) (ret []string) {
    // use note:
    // _,srr:=interfun.MakeSimpleGetToStr(bb.(string))
    // st:=strings.Split(srr, "\n")
    // mytest := func(s string) bool { return strings.HasPrefix(s, "https") }
    // yy:=kkutil.FilterStrings(st, mytest)
    for _, s := range ss {
        if filterfun(s) {
            ret = append(ret, s)
        }
    }
    return
}

/*************************************** test function ************************************/
func TestJson() {
    dinfo:=DevInfo{
        Name:"MyDevice",
        Account:"hello",
        Password:"world",
    }
    adev:=AddDev{
        Uid:"ASDFGHJKLMNBVCXZ",
        DeviceInfo:dinfo,
    }

    jraw, err := json.Marshal(adev)
    if err != nil {
        fmt.Println("error:", err)
    }
    fmt.Println(string(jraw))

    var objmap map[string]*json.RawMessage
    json.Unmarshal(jraw, &objmap)

    fmt.Println(objmap)
    var s string
    json.Unmarshal(*objmap["uuid"], &s)
    fmt.Println(s)

    var d DevInfo
    json.Unmarshal(*objmap["info"], &d)
    fmt.Println(d)
}

func PrintJsonData(f interface{}) {
    for _,m :=range f.([]interface{}) {
        fmt.Println(m)
        for k, v := range m.(map[string]interface{}) {
            switch vv := v.(type) {
            case string:
                fmt.Println(k, "is string", vv)
            case float64:
                fmt.Println(k, "is float64", vv)
            case []interface{}:
                fmt.Println(k, "is an array:")
                for i, u := range vv {
                    fmt.Println(i, u)
                }
            default:
                fmt.Println(k, "is of a type I don't know how to handle")
            }
        }   
    }

}
