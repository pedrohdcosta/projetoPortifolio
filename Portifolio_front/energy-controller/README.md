# energy-controller

This template should help get you started developing with Vue 3 in Vite.

## UX Feedback Components

This project now includes consistent UX feedback components for better user experience:

### New Components

- **LoadingSpinner.vue**: A reusable loading spinner with accessibility features
  - Props: `small?: boolean` - Use small variant for inline buttons
  - Includes `aria-busy` and screen reader text
  
- **SkeletonCard.vue**: A pulsing placeholder card for loading states
  - Use when content is loading to provide visual feedback
  - Fully accessible with `role="status"` and screen reader text

- **TelemetryTable.vue**: Table component for displaying device telemetry data
  - Props: `data: TelemetryData[]`, `loading?: boolean`
  - Emits: `refresh` event
  
- **ConsumptionChart.vue**: Chart component for displaying consumption data
  - Props: `labels: string[]`, `series: number[]`, `title?: string`, `loading?: boolean`
  - Shows skeleton when loading or no data available

### Testing UX Feedback

To test the loading and error states:

1. **Dashboard**: 
   - On initial load, you'll see skeleton cards while data loads
   - Simulated network delay shows loading states
   - Error states include "Try Again" buttons

2. **Devices**: 
   - Adding a device shows a loading spinner in the button
   - Form inputs are disabled during submission
   - Errors are displayed in a prominent error card

3. **Profile**: 
   - Shows skeleton card when loading user data
   - Refresh button shows inline spinner during update
   - Error handling with retry functionality

4. **Login/Register**: 
   - Consistent error message styling
   - Loading spinner in submit button
   - Form disabled during submission

All components follow accessibility best practices with proper ARIA attributes and support for reduced motion preferences.

## Recommended IDE Setup

[VSCode](https://code.visualstudio.com/) + [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur).

## Customize configuration

See [Vite Configuration Reference](https://vite.dev/config/).

## Project Setup

```sh
npm install
```

### Compile and Hot-Reload for Development

```sh
npm run dev
```

### Compile and Minify for Production

```sh
npm run build
```
