# energy-controller

Frontend application for energy consumption monitoring and device management.

## Features

- **Dashboard**: Real-time energy consumption monitoring with interactive charts and telemetry data
- **Device Management**: Select and monitor individual devices
- **Telemetry Table**: Paginated view of telemetry data (time, power, voltage, current)
- **Consumption Chart**: Line chart visualization using Chart.js

## Recommended IDE Setup

[VSCode](https://code.visualstudio.com/) + [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur).

## Customize configuration

See [Vite Configuration Reference](https://vite.dev/config/).

## Project Setup

```sh
npm install
```

### Dependencies

The project includes the following key dependencies:

- **Vue 3**: Progressive JavaScript framework
- **Vue Router**: Official router for Vue.js
- **Pinia**: State management library
- **Axios**: HTTP client for API requests
- **Chart.js 4**: JavaScript charting library
- **vue-chartjs 5**: Vue wrapper for Chart.js

### Compile and Hot-Reload for Development

```sh
npm run dev
```

### Type-Check

```sh
npm run typecheck
```

### Compile and Minify for Production

```sh
npm run build
```

## API Integration

The application integrates with a REST API with the following endpoints:

- `GET /api/devices` - List all devices
- `GET /api/telemetry?device_id=<id>&limit=<num>` - Fetch telemetry data for a specific device

## Components

### ConsumptionChart.vue
Renders a line chart for energy consumption visualization. Accepts `labels` and `series` props.

### TelemetryTable.vue
Displays telemetry data in a paginated table (10 rows per page) with refresh functionality.

## API Clients

### devices.ts
- `listDevices()`: Fetches list of devices

### telemetry.ts
- `fetchTelemetry(params)`: Fetches telemetry data with query parameters

