# CMDB Lite Frontend

This is the frontend for CMDB Lite, built with Vue 3, Vite, and TailwindCSS.

## Features

- Responsive web interface
- Configuration Item management
- CI Type management
- Relationship visualization
- User authentication

## Getting Started

### Prerequisites

- Node.js 16 or higher
- npm or yarn

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/cmdb-lite.git
   cd cmdb-lite/frontend
   ```

2. Copy the environment file:
   ```bash
   cp .env.example .env
   ```

3. Install dependencies:
   ```bash
   npm install
   ```

4. Run the development server:
   ```bash
   npm run serve
   ```

The application will be available at http://localhost:3000.

### Building for Production

1. Build the application:
   ```bash
   npm run build
   ```

2. Preview the production build:
   ```bash
   npm run preview
   ```

### Running with Docker

1. Build the Docker image:
   ```bash
   docker build -t cmdb-lite-frontend .
   ```

2. Run the container:
   ```bash
   docker run -p 80:80 cmdb-lite-frontend
   ```

## Project Structure

```
frontend/
├── public/                # Static assets
├── src/                   # Source code
│   ├── assets/            # CSS and other assets
│   │   └── main.css       # Main CSS file with Tailwind
│   ├── components/        # Reusable Vue components
│   ├── router/            # Vue router configuration
│   │   └── index.js       # Router setup
│   ├── stores/            # Pinia stores for state management
│   ├── utils/             # Utility functions
│   ├── views/             # Vue views
│   │   ├── DashboardView.vue
│   │   ├── CIListView.vue
│   │   ├── CIDetailView.vue
│   │   ├── SettingsView.vue
│   │   └── NotFoundView.vue
│   ├── App.vue            # Root Vue component
│   └── main.js            # Application entry point
├── .env.example           # Environment variables example
├── package.json           # Node.js dependencies
├── vite.config.js         # Vite configuration
├── tailwind.config.js     # TailwindCSS configuration
├── postcss.config.js      # PostCSS configuration
├── index.html             # HTML template
├── nginx.conf             # Nginx configuration for Docker
├── Dockerfile             # Docker configuration for production
└── Dockerfile.dev         # Docker configuration for development
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| VITE_API_BASE_URL | Base URL for the backend API | http://localhost:8080/api |
| VITE_APP_TITLE | Application title | CMDB Lite |
| VITE_APP_VERSION | Application version | 0.1.0 |

## Available Scripts

| Script | Description |
|--------|-------------|
| `npm run serve` | Start the development server |
| `npm run build` | Build for production |
| `npm run preview` | Preview production build |
| `npm run lint` | Run ESLint to fix code issues |

## Styling

This project uses TailwindCSS for styling. The main CSS file is located at `src/assets/main.css` and includes:

- Tailwind directives
- Custom component classes
- Utility classes

## State Management

This project uses Pinia for state management. Stores are located in the `src/stores` directory.

## Routing

This project uses Vue Router for routing. The router configuration is located in `src/router/index.js`.

## API Integration

The frontend integrates with the backend API using Axios. API calls are centralized in utility functions located in the `src/utils` directory.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Commit your changes
6. Push to the branch
7. Create a Pull Request