# Dataroma Portfolio Scraper

A Go-based web scraper that fetches investment portfolio data from Dataroma.com. This tool allows you to retrieve and analyze the portfolio holdings of famous investors like Bill Ackman, Warren Buffett, and others.

## Features

- Fetch complete list of investors from Dataroma
- Retrieve detailed portfolio holdings for specific investors
- Clean data parsing for portfolio weights, share counts, and pricing
- Structured output with proper Go types
- Error handling and input validation

## Prerequisites

- Go 1.16 or higher
- Internet connection to access Dataroma.com

## Installation

1. Clone the repository:
```bash
git clone https://github.com/makalin/dataroma-scraper.git
cd dataroma-scraper
```

2. Install the required dependency:
```bash
go get github.com/PuerkitoBio/goquery
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
)

func main() {
    // Get all investors
    investors, err := GetAllInvestors()
    if err != nil {
        log.Fatal(err)
    }

    // Print all investors
    for _, inv := range investors {
        fmt.Printf("%s (Updated: %s)\n", inv.Name, inv.UpdateDate.Format("2006-01-02"))
    }

    // Get specific investor's portfolio
    portfolio, err := GetInvestorPortfolio("Ackman")
    if err != nil {
        log.Fatal(err)
    }

    // Print portfolio holdings
    for _, holding := range portfolio {
        fmt.Printf("%s (%s): %.2f%%\n", 
            holding.Symbol, 
            holding.Name, 
            holding.PortfolioWeight * 100)
    }
}
```

### Data Structures

The scraper uses two main structs for organizing data:

```go
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
```

### Available Functions

1. `GetAllInvestors() ([]Investor, error)`
   - Fetches and returns a list of all investors from Dataroma's homepage
   - Returns a slice of Investor structs and any error encountered

2. `GetInvestorPortfolio(name string) ([]Portfolio, error)`
   - Fetches portfolio holdings for a specific investor
   - Takes investor name as input (case insensitive)
   - Returns a slice of Portfolio structs and any error encountered

## Example Output

```
All investors:
Bill Ackman (Updated: 2024-03-15)
Warren Buffett (Updated: 2024-03-15)
...

Ackman's portfolio:
GS (Goldman Sachs): 15.20%
GOOG (Alphabet Inc): 12.50%
...
```

## Error Handling

The scraper includes comprehensive error handling for common scenarios:

- Network connectivity issues
- Invalid investor names
- HTML parsing errors
- Date parsing errors
- Number parsing errors

Errors are returned using Go's standard error interface and include descriptive messages.

## Limitations

- Depends on Dataroma.com's HTML structure (may break if the website changes)
- Rate limiting is not implemented
- No caching mechanism
- No direct plotting/visualization capabilities

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This tool is for educational purposes only. Please respect Dataroma.com's terms of service and implement appropriate rate limiting if using this tool in production. The author is not responsible for any misuse of this tool.

## Acknowledgments

- Data provided by [Dataroma.com](https://www.dataroma.com)
- Inspired by the R implementation at [codingfinance.com](https://www.codingfinance.com/post/2020-01-06-web-scraping-dataroma_r/)

## Contact

If you have any questions or suggestions, please open an issue in the GitHub repository.
