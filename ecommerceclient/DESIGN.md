# DESIGN.md

## 1. System Overview
This document defines the core design tokens and visual hierarchy for the e-commerce UI. It is formatted for direct consumption by AI agents and UI component generators. The system enforces a corporate modern aesthetic with specific conversion-oriented accents.

## 2. Theme Tokens (TypeScript)

```typescript
export const theme = {
  palette: {
    primary: {
      main: '#1E3A8A', // Structural baseline, headers, primary actions
      light: '#2563EB',
      dark: '#172554',
      contrastText: '#FFFFFF',
    },
    secondary: {
      main: '#F97316', // High-conversion CTA, limited-time offers
      light: '#FB923C',
      dark: '#EA580C',
      contrastText: '#FFFFFF',
    },
    info: {
      main: '#38BDF8', // Metadata chips, tags, informational alerts
      contrastText: '#121C2A',
    },
    background: {
      default: '#F3F4F6', // Flat canvas layer (Level 0)
      paper: '#FFFFFF',   // Card & surface layer (Level 1)
      footer: '#0B0F17',  // Deep grounding footer
    },
    text: {
      primary: '#374151',
      secondary: '#6B7280',
      disabled: '#9CA3AF',
    },
    border: {
      light: '#E5E7EB',
      main: '#D1D5DB',
      focus: '#2563EB',
    }
  },
  typography: {
    fontFamily: {
      heading: '"Outfit", sans-serif',
      body: '"Inter", sans-serif',
    },
    fontWeight: {
      regular: 400,
      medium: 500,
      semibold: 600,
      bold: 700,
    }
  },
  shape: {
    borderRadius: {
      sm: '0.625rem', // 10px (Inputs, secondary buttons)
      md: '0.75rem',  // 12px (Cards, panels)
      pill: '9999px', // Hero search, badges
    }
  },
  spacing: {
    xs: '0.25rem', // 4px
    sm: '0.5rem',  // 8px
    md: '1rem',    // 16px
    lg: '1.5rem',  // 24px
    xl: '2rem',    // 32px
    xxl: '3rem',   // 48px
  }
} as const;

export type AppTheme = typeof theme;
```

3. Structural & Layout Constraints
3.1. Grid System

    Desktop (>= 1024px): 12-column grid. 24px (1.5rem) gutters, 32px (2rem) page margins. Max container width: 1280px.

    Tablet (640px - 1023px): 8-column grid. 16px (1rem) gutters.

    Mobile (< 640px): 4-column grid. 12px (0.75rem) gutters, 16px (1rem) outer margins.

3.2. Elevation (Box Shadow)

    Level 1 (Default Surface): 0 1px 3px 0 rgba(0, 0, 0, 0.05), 0 1px 2px -1px rgba(0, 0, 0, 0.05). Applied with a 1px solid rgba(229, 231, 235, 0.8) border.

    Level 2 (Hover/Active Surface): 0 10px 15px -3px rgba(30, 58, 138, 0.06), 0 4px 6px -4px rgba(30, 58, 138, 0.04).

    Level 3 (Overlays/Modals): 0 20px 25px -5px rgba(15, 23, 42, 0.1), 0 8px 10px -6px rgba(15, 23, 42, 0.06).

4. Component Rules
4.1. Typography Application

    h1 through h6 must use theme.typography.fontFamily.heading (Outfit).

    Structural headings must map to theme.palette.primary.main (#1E3A8A).

    Body paragraphs, tables, and numeric values must use theme.typography.fontFamily.body (Inter).

    Interactive controls must set text-transform: none and fontWeight: 600.

4.2. Buttons

    Primary: Background #1E3A8A, Text #FFFFFF, Border Radius 10px.

    Conversion/CTA: Background #F97316, Text #FFFFFF, Drop Shadow 0 4px 12px rgba(249, 115, 22, 0.25).

    Secondary/Outlined: Background #FFFFFF, Border 1px solid #D1D5DB, Text #374151.

4.3. Form Inputs

    Base: Background #FFFFFF, Border 1px solid #D1D5DB, Border Radius 10px, Min-Height 44px.

    Focus State: Border changes to #2563EB. Must apply a focus ring: box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.15).

4.4. Header & Footer Elements

    Header Top Nav: Height 72px, Background #FFFFFF, Border-Bottom 1px solid #E5E7EB.

    Footer Container: Background #0B0F17. Flat rendering (no elevation). Text elements must use high-contrast variants (e.g., #38BDF8 for headers, #9CA3AF for descriptive text).