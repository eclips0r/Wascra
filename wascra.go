/**************************************
***************************************
***************************************
@eclisp0r 2022
***************************************
***************************************
**************************************/

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// characteristics of a model
type Model struct {
    Power           string
    TroopSize       string
    Name            string
    Movement        string
    WeaponSkill     string
    BalisticSkill   string
    Strength        string
    Toughness       string
    Wounds          string
    Attacks         string
    Leadership      string
    Save            string
}

// prints model 
func printModel (m *Model) {
    fmt.Println(m.Name)
    fmt.Println("Power: "           + m.Power)
    fmt.Println("TroopSize: "       + m.TroopSize)
    fmt.Println("Movement: "        + m.Movement)
    fmt.Println("WeaponSkill: "     + m.WeaponSkill)
    fmt.Println("BalisticSkill: "   + m.BalisticSkill)
    fmt.Println("Strength: "        + m.Strength)
    fmt.Println("Toughness: "       + m.Toughness)
    fmt.Println("Wounds: "          + m.Wounds)
    fmt.Println("Attacks: "         + m.Attacks)
    fmt.Println("Leadership: "      + m.Leadership)
    fmt.Println("Saveroll: "        + m.Save)
    fmt.Println()
}

// returns formatted string of model
func returnModel (m *Model) string {
    ms := (m.Name         + "\n" +
        "Power: "         + m.Power         + "\n" +
        "TroopSize: "     + m.TroopSize     + "\n" +
        "Movement: "      + m.Movement      + "\n" +
        "WeaponSkill: "   + m.WeaponSkill   + "\n" +
        "BalisticSkill: " + m.BalisticSkill + "\n" +
        "Strength: "      + m.Strength      + "\n" +
        "Toughness: "     + m.Toughness     + "\n" +
        "Wounds: "        + m.Wounds        + "\n" +
        "Attacks: "       + m.Attacks       + "\n" +
        "Leadership: "    + m.Leadership    + "\n" +
        "Saveroll: "      + m.Save          + "\n" +
        "\n")

    return ms
}

func main() {
    // ++++ start of command-line argument definition ++++
    verbosePtr := flag.Bool("v", false, "verbose output")

    writePtr := flag.Bool("w", false, "writes to auto-generated file")

    writeverbPtr := flag.Bool("wv", false, "write verbose")

    // ++++ end of command-line argument definition ++++

    flag.Parse()

    args := flag.Args()

    // domains with protocol for easy change
    domain := "https://wahapedia.ru"
    domainFact := "https://wahapedia.ru/wh40k9ed/factions/"

    // allocate slice with starting length 0 and capacity for 100 entries
    models := make([]Model, 0, 100)
    

    // initialize new Collector and set specific options
    c := colly.NewCollector(
        colly.AllowedDomains("wahapedia.ru", "www.wahapedia.ru"),
        
        // prevents unnecessary download of same data even if collector is restarted
        colly.CacheDir("./cache"),
    )
    c.SetRequestTimeout(10 * time.Second)

    // initialize another collector for model information
    modelCollector := c.Clone()

    c.OnResponse(func(response *colly.Response) {
        fmt.Println("Landing page:", response.Request.URL)
    })

    // traverse army
    c.OnHTML(".NavColumns3", func(element *colly.HTMLElement) {
        element.ForEach("a", func(_ int, h *colly.HTMLElement) {
            // dont visit references on landing page
            if (!strings.Contains(h.Attr("href"), "#")) {
                modelCollector.Visit(domain + h.Attr("href"))
            }
        })
    })

    // error message
    c.OnError(func(response *colly.Response, err error) {
        fmt.Println("Requested URL:", response.Request.URL, "failed with response:", response, "\nError:", err)
    })


    // ++++ HTML selector strings ++++
    m_power := "div.pTable_no_dam > div:nth-child(1) > div"
    m_troopsize := "td.pTable2:nth-child(1)"
    m_name := "title"
    m_movement := "td.pTable1:nth-child(3)"
    m_ws := "td.pTable2:nth-child(4)"
    m_bs := "td.pTable2:nth-child(5)"
    m_strength := "td.pTable2:nth-child(6)"
    m_toughness := "td.pTable2:nth-child(7)"
    m_wounds := "td.pTable1:nth-child(8)"
    m_attacks := "td.pTable2:nth-child(9)"
    m_leadership := "td.pTable2:nth-child(10)"
    m_saveroll := "td.pTable2:nth-child(11)"

    // build model from site data, append to slice
    modelCollector.OnHTML("html", func(element *colly.HTMLElement) {
        temp := Model{
            Power: element.ChildText(m_power),
            TroopSize: element.ChildText(m_troopsize),
            Name: element.ChildText(m_name),
            Movement: element.ChildText(m_movement),
            WeaponSkill: element.ChildText(m_ws),
            BalisticSkill: element.ChildText(m_bs),
            Strength: element.ChildText(m_strength),
            Toughness: element.ChildText(m_toughness),
            Wounds: element.ChildText(m_wounds),
            Attacks: element.ChildText(m_attacks),
            Leadership: element.ChildText(m_leadership),
            Save: element.ChildText(m_saveroll),
        }

        models = append(models, temp)
    })

    // response message
    modelCollector.OnResponse(func(response *colly.Response) {
        if (*verbosePtr) {
            fmt.Println("Visited model URL:", response.Request.URL)
        }
    })

 
    // ++++ entrypoint of collector ++++
    for i := range args {
        faction := strings.ToLower(args[i])
        c.Visit(domainFact + faction)
    }

    // write scraped models to file (./faction/argx_argy.txt)
    if (*writePtr) {
        filePath := "./factions/"

        // concatenate arguments to filename
        for i := 0; i < len(args) - 1; i++ {
            filePath = filePath + args[i] + "_"
        }

        filePath = filePath + args[len(args) - 1] + ".txt"

        file, err := os.OpenFile(filePath, os.O_CREATE | os.O_WRONLY, 0644)

        if (err != nil) {
            panic(err)
        }
        
        for i := range models {
            l, err := file.WriteString(returnModel(&models[i]))

            if (err != nil) {
                panic(err)
            }

            if (*writeverbPtr) {
                fmt.Println("Model written with", l, "bytes")
            }
        }
        file.Close()
        fmt.Println("done")
        
    } else {
        for i := range models {
            printModel(&models[i])
        }
        fmt.Println("done")
    }
}
