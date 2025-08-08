# Meeting Room Booking System - Planning Document

## Design Analysis
Based on the reference image (landing-date.png), the system features:

### UI/UX Elements
- **Header**: Pink/magenta navigation bar with date picker, navigation links (HELP, ROOMS, REPORT), search bar, and user login
- **Sidebar**: Floor selection (1st Floor, 2nd Floor, 3rd Floor, 4th Floor)
- **Main Calendar**: Three-month view (July, August, September 2025) with date navigation
- **Day View**: Time slots from 07:00 to 17:00+ with room columns
- **Color Coding**: 
  - Pink header and time slots
  - Purple/blue booking blocks
  - Clean white background with subtle borders

### Key Features to Implement
1. **Authentication System**: User login/logout
2. **Date Navigation**: Calendar picker and navigation controls
3. **Floor/Area Selection**: Multi-floor room management
4. **Room Booking Grid**: Time-based booking interface
5. **Booking Management**: Create, view, edit bookings
6. **Responsive Design**: Mobile and desktop compatibility

### Technical Requirements
- Next.js 14+ with App Router
- TypeScript for type safety
- Tailwind CSS for styling
- Shadcn/ui components for consistent UI
- Date management with date-fns or similar
- State management (Context API or Zustand)
- Mock data for rooms and bookings

### Component Architecture
```
src/app/
├── page.tsx (Main dashboard)
├── layout.tsx (Root layout with header)
├── globals.css
├── components/
│   ├── ui/ (shadcn components)
│   ├── Header/
│   ├── Sidebar/
│   ├── Calendar/
│   ├── BookingGrid/
│   ├── BookingModal/
│   └── Auth/
├── lib/
│   ├── utils.ts
│   ├── data.ts (mock data)
│   └── types.ts
└── hooks/
    └── use-bookings.ts
```

### Color Scheme (Based on Reference)
- Primary: Pink/Magenta (#E91E63 or similar)
- Secondary: Purple gradients (#9C27B0 to #673AB7)
- Background: Clean white (#FFFFFF)
- Text: Dark gray (#374151)
- Borders: Light gray (#E5E7EB)

### Development Phases
1. **Phase 1**: Setup project structure and basic layout
2. **Phase 2**: Implement header and navigation
3. **Phase 3**: Create calendar and date navigation
4. **Phase 4**: Build booking grid and time slots
5. **Phase 5**: Add booking functionality and modals
6. **Phase 6**: Implement responsive design and polish

### Mock Data Structure
- Rooms: Floor-based room categorization
- Bookings: Time-based booking entries
- Users: Basic user authentication
- Time slots: 30-minute intervals from 7:00 AM to 6:00 PM+