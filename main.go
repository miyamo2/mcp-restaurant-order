package main

import (
	"encoding/json"
	"github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"log/slog"
)

type GetMenuArguments struct {
	Name         *string `json:"name" jsonschema_description:"The name of the menu."`
	LowestPrice  *uint32 `json:"lowest_price" jsonschema_description:"The lowest price of the menu. in japanese yen."`
	HighestPrice *uint32 `json:"highest_price" jsonschema_description:"The highest price of the menu. in japanese yen."`
	Category     *string `json:"category" jsonschema_description:"The category of the menu. Only the following values are valid. main | tapas | dessert | beverage"`
}

type OrderArguments struct {
	Orders []Order `json:"orders" jsonschema:"description=The list of orders"`
}

type Order struct {
	Name     string `json:"name" jsonschema:"required" jsonschema_description:"The name of the menu. must be exist in the menu. also, case sensitive."`
	Quantity uint32 `json:"quantity" jsonschema:"required" jsonschema_description:"The quantity per menu"`
}
type Menu struct {
	Name     string `json:"name" jsonschema:"required" jsonschema_description:"The name of the menu."`
	Price    uint32 `json:"price" jsonschema:"required" jsonschema_description:"The price of the menu. in japanese yen."`
	Category string `json:"category" jsonschema:"required" jsonschema_description:"The category of the menu. only the following values are provide. main | tapas | dessert | beverage"`
}

type BillingDetail struct {
	Name     string `json:"name" jsonschema:"required" jsonschema_description:"The name of the menu."`
	Quantity uint32 `json:"quantity" jsonschema:"required" jsonschema_description:"The quantity per menu"`
	Amount   uint32 `json:"amount" jsonschema:"required" jsonschema_description:"The amount per menu. in japanese yen."`
}

type Billing struct {
	Details     []BillingDetail `json:"billing_details" jsonschema:"required" jsonschema_description:"The list of billing details."`
	TotalAmount uint32          `json:"total_amount" jsonschema:"required" jsonschema_description:"The total amount of the order. in japanese yen."`
}

func main() {
	done := make(chan struct{})
	server := mcp_golang.NewServer(stdio.NewStdioServerTransport())
	err := server.RegisterTool("miyamo2_diner_search_menu", "searchs miyamo2 diner's menu for dishes matching your criteria.", func(arguments GetMenuArguments) (*mcp_golang.ToolResponse, error) {
		filteredMenu := make([]Menu, 0, len(menus))
		for _, menu := range menus {
			if arguments.Name != nil && *arguments.Name != menu.Name {
				continue
			}
			if arguments.LowestPrice != nil && *arguments.LowestPrice > menu.Price {
				continue
			}
			if arguments.HighestPrice != nil && *arguments.HighestPrice < menu.Price {
				continue
			}
			if arguments.Category != nil && *arguments.Category != menu.Category {
				continue
			}
			filteredMenu = append(filteredMenu, menu)
		}
		result, err := json.Marshal(filteredMenu)
		if err != nil {
			return nil, err
		}
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(result))), nil
	})
	if err != nil {
		panic(err)
	}
	err = server.RegisterTool("miyamo2_diner_accept_order", "miyamo2 diner will accept your order and return the total amount in Japanese yen.", func(arguments OrderArguments) (*mcp_golang.ToolResponse, error) {
		var (
			details     []BillingDetail
			totalAmount uint32
		)
		for _, order := range arguments.Orders {
			menu, ok := menus[order.Name]
			if !ok {
				continue
			}
			if order.Quantity == 0 {
				continue
			}
			amount := menu.Price * order.Quantity
			totalAmount += amount
			details = append(details, BillingDetail{
				Name:     menu.Name,
				Quantity: order.Quantity,
				Amount:   amount,
			})
		}
		billing := Billing{
			Details:     details,
			TotalAmount: totalAmount,
		}
		result, err := json.Marshal(billing)
		if err != nil {
			return nil, err
		}
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(result))), nil
	})

	err = server.Serve()
	if err != nil {
		panic(err)
	}
	<-done
	slog.Info("Server stopped")
}

var menus = map[string]Menu{
	"Pizza": {
		Name:     "Pizza",
		Price:    1500,
		Category: "main",
	},
	"Burger": {
		Name:     "Burger",
		Price:    1200,
		Category: "main",
	},
	"Pasta": {
		Name:     "Pasta",
		Price:    1200,
		Category: "main",
	},
	"Steak": {
		Name:     "Steak",
		Price:    2500,
		Category: "main",
	},
	"Salad": {
		Name:     "Salad",
		Price:    800,
		Category: "tapas",
	},
	"French Fries": {
		Name:     "French Fries",
		Price:    600,
		Category: "tapas",
	},
	"Nachos": {
		Name:     "Nachos",
		Price:    800,
		Category: "tapas",
	},
	"Ice Cream": {
		Name:     "Ice Cream",
		Price:    500,
		Category: "dessert",
	},
	"Cheesecake": {
		Name:     "Cheesecake",
		Price:    700,
		Category: "dessert",
	},
	"Soda": {
		Name:     "Soda",
		Price:    300,
		Category: "beverage",
	},
	"Hot Coffee": {
		Name:     "Hot Coffee",
		Price:    300,
		Category: "beverage",
	},
	"Iced Coffee": {
		Name:     "Iced Coffee",
		Price:    300,
		Category: "beverage",
	},
	"Beer": {
		Name:     "Beer",
		Price:    500,
		Category: "beverage",
	},
}
