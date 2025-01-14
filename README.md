# The REGoTTH Stack

A modern web application stack for the fullstack dev that wants something different from the usual stack.
REGoTTH stands for **Re**act, **Go**, **T**ailwind, **T**empl, and **H**TMX.

This project serves as a reference implementation demonstrating how to effectively integrate these technologies in a clean and opinionated way.

## 📑 Table of Contents

- [🛠 Tech Stack](#-tech-stack)
- [📁 Project Structure](#-project-structure)
- [🏗 Architecture](#-architecture)
- [🚀 Getting Started](#-getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Development](#development)
- [🎯 React Components Integration](#-react-components-integration)
  - [Adding New React Components](#adding-new-react-components)
- [🧪 Testing](#-testing)
- [🔧 Available Make Commands](#-available-make-commands)
  - [Development Commands](#development-commands)
  - [Setup Commands](#setup-commands)
- [🔍 Troubleshooting](#-troubleshooting)
  - [Common Issues](#common-issues)
  - [Still Having Issues?](#still-having-issues)
- [📝 License](#-license)

## 🛠 Tech Stack

- **Go** - Backend server and business logic
- **HTMX** - Frontend interactivity without JavaScript
- **Templ** - Type-safe HTML templating
- **TailwindCSS** - Utility-first CSS framework
- **React** - Component-based UI library
- **Bun** - JavaScript runtime and package manager

## 📁 Project Structure

```bash
.
├── internal/
│ ├── entity/ # Business entities/models
│ ├── repository/ # Data access layer
│ ├── service/ # Business logic layer
│ ├── server/ # HTTP server configuration
│ └── middleware/ # HTTP middleware
├── web/
│ ├── component/ # Reusable UI components
│ ├── view/ # Page templates
│ └── controller/ # HTTP request handlers
├── static/
│ ├── css/ # TailwindCSS files
│ ├── js/ # TypeScript & React components
│ └── build/ # Compiled assets
└── main.go # Application entry point
```

## 🏗 Architecture

The project follows a clean architecture pattern:

1. **Entity Layer** (`internal/entity/`) - Contains business models
2. **Repository Layer** (`internal/repository/`) - Handles data persistence
3. **Service Layer** (`internal/service/`) - Implements business logic
4. **Controller Layer** (`web/controller/`) - Handles HTTP requests
5. **View Layer** (`web/views/` & `web/components/`) - Manages UI templates
6. **React Components** (`static/js/components/`) - Client-side interactive components

## 🚀 Getting Started

### Prerequisites

- Go 1.23.4 or higher
- Make (for running commands)
- Air (for live reload) - `go install github.com/air-verse/air@latest`
- Templ (for generating templates) - `go install github.com/a-h/templ/cmd/templ@latest`
- Bun (for running the JavaScript/React components) - `brew install bun`
- jq (optional: used for renaming the module) - `brew install jq`

### Installation

1. Clone the repository

    ```bash
    git clone https://github.com/Neofox/regotth.git my-project
    cd my-project
    make rename-module NEW_NAME=github.com/username/my-project
    rm -rf .git
    git init
    ```

2. Install dependencies

```bash
go mod download
bun install
```

To see all available commands:

```bash
make help
```

### Development

For development with live reload:

```bash
make live
```

This will start:

- Go server with hot reload
- Templ template generation with hot reload
- TailwindCSS compilation
- React components compilation with hot reload

> 💡 **Note:** Always use the proxy URL (http://localhost:7331) during development to ensure all live-reload features work correctly.

### Building for Production

```bash
make build
```

This will build the project for production and output the binary to the `tmp` directory.

## 🎯 React Components Integration

### Adding New React Components

1. Create your component in `static/js/components/`:

    ```tsx
    // As we also use preact/compat we can use the normal react syntax
    // /!\ React has to be imported even if we don't use it!
    import React, { useState, type ReactNode } from "react";

    export function MyComponent(): ReactNode {
        const [count, setCount] = useState(0);
        return <div>My Component {count}</div>;
    }
    ```

2. Use it in your Templ templates:

    ```html
    <div data-react-component="MyComponent" />
    ```

The component will be automatically loaded when it appears in the DOM.

### Type Generation

The project automatically generates Go types from React component props, ensuring type safety between React components and Templ templates. When you create or modify a React component with props:

1. Define your props interface in your React component:

    ```tsx
    export interface MyComponentProps {
        prop1: string;
        prop2?: number;
    }
    ```

2. The types will be automatically generated when:

    - Running `make build`
    - During development with `make live`
    - Manually with `bun run generate-props`
  The generated types will be available in `web/generated/react_component_props.go`

3. Use the generated types in your Templ templates:

    ```html
    <div 
      data-react-component="Button" 
      data-react-props={ props.ButtonProps{
          Id: "123",
          Label: "Click me",
      }.String() }
    />
    ```

## 🧪 Testing

### Running Tests

```bash
go test ./...
```

## 🔧 Available Make Commands

Run `make help` to see all available commands. Here are the main ones:

### Development Commands

- `make live` - Start development server with live reload
- `make build` - Build the project for production

### Setup Commands

- `make rename-module NEW_NAME=your-module-name` - Rename the Go module and package name

## 🔍 Troubleshooting

### Common Issues

#### Module rename issues

1. If `make rename-module` fails:

   ```bash
   # Clean up any partial changes
   git checkout go.mod
   git checkout package.json
   # Ensure jq is installed
   brew install jq
   # Try again with the full path
   make rename-module NEW_NAME=github.com/username/project-name
   ```

#### Build errors

1. "templ command not found":

   ```bash
   go install github.com/a-h/templ/cmd/templ@latest
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

2. "bun command not found":

   ```bash
   curl -fsSL https://bun.sh/install | bash
   source ~/.bashrc  # or source ~/.zshrc
   ```

3. Missing dependencies:

   ```bash
   # Reset Go modules
   rm -rf go.sum
   go mod tidy
   
   # Reset npm modules
   rm -rf node_modules
   bun install
   ```

### Still Having Issues?

1. Verify all prerequisites are installed:

   ```bash
   go version      # Should be 1.23.4 or higher
   templ version
   bun --version
   ```

2. Clean and rebuild:

   ```bash
   # Full cleanup
   rm -rf tmp/ node_modules/ static/build/*
   
   # Fresh install
   go mod download
   bun install
   make build
   ```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
