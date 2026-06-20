# FocusFlow - Phase 1-4 Walkthrough

I have completed the core infrastructure and the initial version of the Linux Agent and Web Dashboard.

## Key Accomplishments

### 1. Database Schema (Supabase)
- Defined the core schema in [schema.sql](file:///home/devu/Desktop/FocusFlow/supabase/schema.sql).
- Tables: `categories`, `activity_logs`, `app_category_mappings`.
- Integrated Row Level Security (RLS) for data privacy.

### 2. Linux Go Agent
- Created a robust daemon in the [agent/](file:///home/devu/Desktop/FocusFlow/agent/) directory.
- **Features:**
    - Active window polling (X11).
    - Idle detection (using `xprintidle`).
    - Smart session batching (minimizes DB writes).
    - Offline buffering (using SQLite, no internet required).
    - Background sync to Supabase.
    - systemd service configuration for reliability.

### 3. Web Dashboard
- Built a premium, dark-mode React dashboard in the [dashboard/](file:///home/devu/Desktop/FocusFlow/dashboard/) directory.
- **Features:**
    - Secure login via Supabase Auth.
    - Real-time productivity charts (Pie and Bar charts).
    - Detailed activity timeline (dummy data for preview).
    - Responsive design using Tailwind CSS and Lucide icons.

## How to Run

### Linux Agent
1. Install dependencies: `sudo apt install xdotool xprintidle`
2. Create `~/.config/focusflow/config.json` with your Supabase credentials:
   ```json
   {
     "supabase_url": "YOUR_URL",
     "supabase_key": "YOUR_KEY",
     "device_id": "laptop-1"
   }
   ```
3. Build and run:
   ```bash
   cd agent
   go build -o focusflow-agent
   ./focusflow-agent
   ```

### Web Dashboard
1. Create `dashboard/.env` based on [dashboard/.env.example](file:///home/devu/Desktop/FocusFlow/dashboard/.env.example).
2. Run development server:
   ```bash
   cd dashboard
   npm run dev
   ```

## Next Steps
- **Phase 5:** Implement the Browser Extension for URL-level tracking.
- **Phase 6:** Develop the Android App.
