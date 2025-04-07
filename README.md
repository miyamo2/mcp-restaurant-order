# mcp-restaurant-order

MCP(Model Context Protocol) server implementation. Just a toy.

## Demo

### With English Prompt

![Demo](./.assets/demo.gif)

### With Japanese Prompt

![Demo_ja](./.assets/demo_ja.gif)

## Quick Start

### Install

```sh
go install github.com/miyamo2/mcp-restaurant-order@latest
```

### Use in Claude Desktop

```json
{
    "mcpServers": {
        "mcp-restaurant-order": {
          "command": "mcp-restaurant-order"
        }
    }
}
```

### License

**mcp-restaurant-order** released under the [MIT License](https://github.com/miyamo2/mcp-restaurant-order/blob/main/LICENSE)
