# Go Todo Application

A modern todo application demonstrating clean architecture using Go, HTMX, Templ, and React components. This project serves as an example of how to structure a web application using these technologies.

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
│ ├── components/ # Reusable UI components
│ ├── views/ # Page templates
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
- Air (for live reload) - `go install github.com/cosmtrek/air@latest`
- Templ - `go install github.com/a-h/templ/cmd/templ@latest`
- Bun (for running commands) - `brew install bun`

### Installation

1. Clone the repository

    ```bash
    git clone https://github.com/Neofox/go-todo.git
    cd go-todo
    ```

2. Install dependencies

    ```bash
    go mod download
    bun install
    ```

3. Build the project

    ```bash
    make build
    ```

### Development

For development with live reload:

```bash
make live
```

This will start:

- Go server with hot reload
- Templ template generation
- TailwindCSS compilation
- Asset synchronization

### Building for Production

```bash
make build
```

This will build the project for production and output the binary to the `tmp` directory.

## 🎯 React Components Integration

### Adding New React Components

1. Create your component in `static/js/components/`:

    ```tsx
    import React from 'react';

    export function MyComponent(props) {
        return <div>My Component</div>;
    }
    ```

2. Register it in `static/js/index.ts`:

    ```tsx
    const componentRegistry = {
        MyComponent: MyComponent,
    // Add more components here
    };
    ```

3. Use it in your Templ templates:

    ```go
    <div data-react-component="MyComponent" data-react-props={ templ.JSONString(map[string]string{
        "prop1": "value1",
        "prop2": "value2",
    }) }></div>
    ```

## 🧪 Testing

### Running Tests

```bash
go test ./...
```

## 🔧 Available Make Commands

- `make live` - Start development server with live reload
- `make build` - Build the project for production
- `make templ-generate` - Generate Templ templates
- `make javascript-build` - Build JavaScript/React components
- `make javascript-watch` - Watch JavaScript changes

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
