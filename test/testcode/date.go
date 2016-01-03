package main

import (
"fmt"
"time"
"strings"
)


func main() {
    // s := "10.0.0.1 Jan 11 2014 10:00:00 hello"
    // r := regexp.MustCompile("^([^/w]+) ([a-zA-Z]+ [0-9]{1,2} [0-9]{4} [0-9]{1,2}:[0-9]{2}:[0-9]{2}) (.*)")
    // m := r.FindStringSubmatch(s)
    // if len(m) >= 4 {
    //     fmt.Println("IP:", m[1])
    //     fmt.Println("Timestamp:", m[2])
    //     fmt.Println("Message:", m[3])
    //     t, err := time.Parse("Jan 02 2006 15:04:05", m[2])
    //     if err != nil {
    //         fmt.Println(err.Error())
    //     } else {
    //         fmt.Println("Parsed Time:",t)
    //     }
    // } else {
    //        fmt.Println("Regexp mismatch!")
    // }

   // str := "20150901";
   //  size := len(str);
   //  fmt.Printf(" size = %d \n", size);

   //  MachineDate  := strings.Replace(time.Now().String()[0:10], "-", "", -1)
   //  fmt.Println(MachineDate)

    package main
 
import (
    "fmt"
    "bufio"
    "os"
    "regexp"
    "strconv"
)
 
func main() {
    year := input("year", "^[0-9]{1}[0-9]{3}$")
    month := input("month", "^(0{1}[0-9]{1}|1{1}[0-2]{1})$")
    count(year, month)
    fmt.Println("Press Enter button to continue ...")
    reader := bufio.NewReader(os.Stdin)
    lastInput, _, err := reader.ReadRune()
    if err != nil {
        fmt.Fprintln(os.Stderr, "Occur error when input (last) '", lastInput, "':", err)
    }
    return
}
 
func count(year int, month int) (days int) {
    if month != 2 {
        if month == 4 || month == 6 || month == 9 || month == 11 {
            days = 30
 
        } else {
            days = 31
            fmt.Fprintln(os.Stdout, "The month has 31 days");
        }
    } else {
        if (((year % 4) == 0 && (year % 100) != 0) || (year % 400) == 0) {
            days = 29
        } else {
            days = 28
        }
    }
    fmt.Fprintf(os.Stdout, "The %d-%d has %d days.\n", year, month, days)
    return
}
 
func input(name string, regexpText string) (number int) {
    var validNumber = false
    for !validNumber {
        fmt.Println("Please input a", name, ": ")
        reader := bufio.NewReader(os.Stdin)
        inputBytes, _, err := reader.ReadLine()
        if err != nil {
            fmt.Fprintln(os.Stderr, "Occur error when input", name, ":", err)
            continue
        }
        inputText := string(inputBytes)
        validNumber, err = regexp.MatchString(regexpText, inputText)
        if err != nil {
            fmt.Fprintln(os.Stderr, "Occur error when match", name, "(", inputText, "):",err)
            continue
        }
        if validNumber {
            number, err = strconv.Atoi(inputText)
            if err != nil {
                fmt.Fprintln(os.Stderr, "Occur error when convert", name, "(", inputText, "):", err)
                continue
            }
        } else {
            fmt.Fprintln(os.Stdout, "The", name, "(", inputText, ") does not have the correct format!")
        }
    }
    fmt.Println("The input", name, ": ", number)
    return
}



}