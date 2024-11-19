package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Investor struct {
	Name       string
	UpdateDate time.Time
	URL        string
}

type Portfolio struct {
	Symbol          string
	Name            string
	PortfolioWeight float64
	Shares          float64
	CostPrice       float64
	Value           float64
}

// GetAllInvestors fetches the list of all investors from the homepage
func GetAllInvestors() ([]Investor, error) {
	url := "https://www.dataroma.com/m/home.php"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var investors []Investor
	doc.Find("#port_body li").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		parts := strings.Split(text, "Updated")
		if len(parts) != 2 {
			return
		}

		name := strings.TrimSpace(parts[0])
		dateStr := strings.TrimSpace(parts[1])
		
		// Parse date (assuming format DD/MM/YYYY)
		date, err := time.Parse("02/01/2006", dateStr)
		if err != nil {
			log.Printf("Failed to parse date %s: %v", dateStr, err)
			return
		}

		href, exists := s.Find("a").Attr("href")
		if !exists {
			return
		}

		investors = append(investors, Investor{
			Name:       name,
			UpdateDate: date,
			URL:        "https://www.dataroma.com" + href,
		})
	})

	return investors, nil
}

// GetInvestorPortfolio fetches the portfolio for a specific investor
func GetInvestorPortfolio(name string) ([]Portfolio, error) {
	// Get all investors first
	investors, err := GetAllInvestors()
	if err != nil {
		return nil, err
	}

	// Find matching investor
	name = strings.Title(strings.ToLower(name))
	var investorURL string
	for _, inv := range investors {
		if strings.Contains(inv.Name, name) {
			investorURL = inv.URL
			break
		}
	}

	if investorURL == "" {
		return nil, fmt.Errorf("investor %s not found", name)
	}

	// Fetch the portfolio page
	resp, err := http.Get(investorURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch portfolio: %v", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse portfolio HTML: %v", err)
	}

	var portfolio []Portfolio
	doc.Find("#grid tr").Each(func(i int, s *goquery.Selection) {
		// Skip header row
		if i == 0 {
			return
		}

		// Extract cells
		var cells []string
		s.Find("td").Each(func(_ int, cell *goquery.Selection) {
			cells = append(cells, strings.TrimSpace(cell.Text()))
		})

		if len(cells) < 5 {
			return
		}

		// Parse stock symbol and name
		stockParts := strings.Split(cells[0], "-")
		if len(stockParts) != 2 {
			return
		}

		// Helper function to parse numbers
		parseNumber := func(s string) float64 {
			// Remove any non-numeric characters except decimal point
			re := regexp.MustCompile(`[^0-9.]`)
			clean := re.ReplaceAllString(s, "")
			num, err := strconv.ParseFloat(clean, 64)
			if err != nil {
				return 0
			}
			return num
		}

		// Create portfolio entry
		entry := Portfolio{
			Symbol:          strings.TrimSpace(stockParts[0]),
			Name:            strings.TrimSpace(stockParts[1]),
			PortfolioWeight: parseNumber(cells[1]) / 100, // Convert percentage to decimal
			Shares:          parseNumber(cells[2]),
			CostPrice:       parseNumber(cells[3]),
			Value:           parseNumber(cells[4]),
		}

		portfolio = append(portfolio, entry)
	})

	return portfolio, nil
}

func main() {
	// Example usage
	investors, err := GetAllInvestors()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("All investors:")
	for _, inv := range investors {
		fmt.Printf("%s (Updated: %s)\n", inv.Name, inv.UpdateDate.Format("2006-01-02"))
	}

	fmt.Println("\nFetching Ackman's portfolio:")
	portfolio, err := GetInvestorPortfolio("Ackman")
	if err != nil {
		log.Fatal(err)
	}

	for _, holding := range portfolio {
		fmt.Printf("%s (%s): %.2f%%\n", 
			holding.Symbol, 
			holding.Name, 
			holding.PortfolioWeight * 100)
	}
}
