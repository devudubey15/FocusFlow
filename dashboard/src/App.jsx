import React, { useState, useEffect } from 'react';
import { supabase } from './lib/supabase';
import { 
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Cell,
  PieChart, Pie
} from 'recharts';
import { Monitor, Phone, Clock, Book, Play, Hash, LogOut, LayoutDashboard, Settings } from 'lucide-react';

const COLORS = ['#2E75B6', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6'];

function App() {
  const [session, setSession] = useState(null);
  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    supabase.auth.getSession().then(({ data: { session } }) => {
      setSession(session);
    });

    const {
      data: { subscription },
    } = supabase.auth.onAuthStateChange((_event, session) => {
      setSession(session);
    });

    return () => subscription.unsubscribe();
  }, []);

  if (!session) {
    return <Login />;
  }

  return (
    <div className="flex h-screen bg-slate-950 text-slate-50 overflow-hidden font-sans">
      {/* Sidebar */}
      <aside className="w-64 border-r border-slate-800 flex flex-col p-6 space-y-8">
        <div className="flex items-center space-x-3">
          <div className="w-10 h-10 bg-blue-600 rounded-xl flex items-center justify-center">
            <LayoutDashboard size={24} />
          </div>
          <span className="text-xl font-bold tracking-tight">FocusFlow</span>
        </div>

        <nav className="flex-1 space-y-2">
          <NavItem icon={<LayoutDashboard size={20} />} label="Overview" active />
          <NavItem icon={<Clock size={20} />} label="Timeline" />
          <NavItem icon={<Book size={20} />} label="Categories" />
          <NavItem icon={<Settings size={20} />} label="Settings" />
        </nav>

        <button 
          onClick={() => supabase.auth.signOut()}
          className="flex items-center space-x-3 text-slate-400 hover:text-white transition-colors p-2"
        >
          <LogOut size={20} />
          <span>Logout</span>
        </button>
      </aside>

      {/* Main Content */}
      <main className="flex-1 overflow-y-auto p-10">
        <header className="flex justify-between items-start mb-10">
          <div>
            <h1 className="text-3xl font-bold mb-2 text-slate-50">Pulse Check</h1>
            <p className="text-slate-400">Here's how your day is looking across all devices.</p>
          </div>
          <div className="flex space-x-4">
             <div className="bg-slate-900 border border-slate-800 p-4 rounded-2xl flex items-center space-x-4">
                <div className="p-3 bg-blue-500/10 text-blue-500 rounded-lg">
                  <Monitor size={24} />
                </div>
                <div>
                  <p className="text-xs text-slate-500 font-medium uppercase tracking-wider">Linux Agent</p>
                  <p className="text-sm font-semibold">Online</p>
                </div>
             </div>
             <div className="bg-slate-900 border border-slate-800 p-4 rounded-2xl flex items-center space-x-4">
                <div className="p-3 bg-emerald-500/10 text-emerald-500 rounded-lg">
                  <Phone size={24} />
                </div>
                <div>
                  <p className="text-xs text-slate-500 font-medium uppercase tracking-wider">Android App</p>
                  <p className="text-sm font-semibold">Active</p>
                </div>
             </div>
          </div>
        </header>

        {/* Dashboard Grid */}
        <div className="grid grid-cols-12 gap-8">
          {/* Summary Stats */}
          <div className="col-span-8 space-y-8">
            <div className="bg-slate-900/50 border border-slate-800 rounded-3xl p-8 backdrop-blur-sm">
              <h2 className="text-lg font-semibold mb-6">Productive Hours</h2>
              <div className="h-64">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={dummyData}>
                    <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#1e293b" />
                    <XAxis dataKey="name" axisLine={false} tickLine={false} tick={{fill: '#64748b', fontSize: 12}} />
                    <YAxis axisLine={false} tickLine={false} tick={{fill: '#64748b', fontSize: 12}} />
                    <Tooltip 
                      contentStyle={{ backgroundColor: '#0f172a', border: '1px solid #1e293b', borderRadius: '12px' }}
                      itemStyle={{ color: '#f8fafc' }}
                    />
                    <Bar dataKey="minutes" radius={[6, 6, 0, 0]}>
                      {dummyData.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} fillOpacity={0.8} />
                      ))}
                    </Bar>
                  </BarChart>
                </ResponsiveContainer>
              </div>
            </div>

            <div className="bg-slate-900/50 border border-slate-800 rounded-3xl p-8 backdrop-blur-sm">
              <h2 className="text-lg font-semibold mb-6">Recent Activity</h2>
              <div className="space-y-4">
                <ActivityItem app="VS Code" category="Work" duration="45m" time="Just now" icon={<LayoutDashboard size={18} />} />
                <ActivityItem app="Chrome" title="GitHub" category="Work" duration="12m" time="10m ago" icon={<Monitor size={18} />} />
                <ActivityItem app="YouTube" category="Entertainment" duration="24m" time="30m ago" icon={<Play size={18} />} />
                <ActivityItem app="Obsidian" category="Study" duration="1h 10m" time="1h ago" icon={<Book size={18} />} />
              </div>
            </div>
          </div>

          {/* Side Panels */}
          <div className="col-span-4 space-y-8">
            <div className="bg-slate-900/50 border border-slate-800 rounded-3xl p-8 backdrop-blur-sm flex flex-col items-center">
              <h2 className="text-lg font-semibold mb-6 w-full text-left">Category Split</h2>
              <div className="h-48 w-full">
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={dummyPie}
                      innerRadius={60}
                      outerRadius={80}
                      paddingAngle={8}
                      dataKey="value"
                    >
                      {dummyPie.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                      ))}
                    </Pie>
                    <Tooltip />
                  </PieChart>
                </ResponsiveContainer>
              </div>
              <div className="mt-4 grid grid-cols-2 gap-4 w-full">
                {dummyPie.map((item, idx) => (
                  <div key={item.name} className="flex items-center space-x-2">
                    <div className="w-3 h-3 rounded-full" style={{backgroundColor: COLORS[idx % COLORS.length]}} />
                    <span className="text-sm text-slate-400">{item.name}</span>
                  </div>
                ))}
              </div>
            </div>

            <div className="bg-gradient-to-br from-blue-600 to-indigo-700 rounded-3xl p-8 text-white">
              <h3 className="text-lg font-semibold mb-2">Deep Work Streak</h3>
              <p className="text-blue-100 text-sm mb-6">You've been focused for 2.5 hours today. 80% of your target.</p>
              <div className="w-full bg-blue-900/40 rounded-full h-3 mb-4">
                <div className="bg-white rounded-full h-3 w-4/5" />
              </div>
              <span className="text-xs font-medium uppercase tracking-widest text-blue-200">Keep it up!</span>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}

function NavItem({ icon, label, active = false }) {
  return (
    <div className={`flex items-center space-x-3 p-3 rounded-xl transition-all cursor-pointer ${
      active ? 'bg-blue-600/10 text-blue-400 font-semibold' : 'text-slate-400 hover:bg-slate-900 hover:text-slate-200'
    }`}>
      {icon}
      <span>{label}</span>
    </div>
  );
}

function ActivityItem({ app, title, category, duration, time, icon }) {
  return (
    <div className="flex items-center justify-between p-4 bg-slate-800/30 rounded-2xl border border-slate-800/50 hover:border-slate-700 transition-all">
      <div className="flex items-center space-x-4">
        <div className="p-3 bg-slate-900 rounded-xl text-slate-400">
          {icon}
        </div>
        <div>
          <p className="font-semibold text-slate-200">{app}</p>
          <p className="text-xs text-slate-500">{title ? title : category}</p>
        </div>
      </div>
      <div className="text-right">
        <p className="text-sm font-medium text-slate-300">{duration}</p>
        <p className="text-xs text-slate-500">{time}</p>
      </div>
    </div>
  );
}

function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);

  const handleLogin = async (e) => {
    e.preventDefault();
    setLoading(true);
    const { error } = await supabase.auth.signInWithPassword({ email, password });
    if (error) alert(error.message);
    setLoading(false);
  };

  return (
    <div className="min-h-screen bg-slate-950 flex items-center justify-center p-6">
      <div className="w-full max-w-md bg-slate-900 border border-slate-800 rounded-3xl p-10 shadow-2xl">
        <div className="flex justify-center mb-8">
           <div className="w-16 h-16 bg-blue-600 rounded-2xl flex items-center justify-center shadow-lg shadow-blue-600/20">
             <LayoutDashboard size={32} className="text-white" />
           </div>
        </div>
        <h1 className="text-3xl font-bold text-center mb-2 text-white">FocusFlow</h1>
        <p className="text-slate-500 text-center mb-10">Sign in to track your flow.</p>
        
        <form onSubmit={handleLogin} className="space-y-6">
          <div className="space-y-2">
            <label className="text-sm font-medium text-slate-400 ml-1">Email</label>
            <input 
              type="email" 
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full bg-slate-950 border border-slate-800 rounded-2xl p-4 text-white focus:outline-none focus:ring-2 focus:ring-blue-600/50 transition-all font-mono"
              placeholder="you@example.com"
              required
            />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium text-slate-400 ml-1">Password</label>
            <input 
              type="password" 
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full bg-slate-950 border border-slate-800 rounded-2xl p-4 text-white focus:outline-none focus:ring-2 focus:ring-blue-600/50 transition-all font-mono"
              placeholder="••••••••"
              required
            />
          </div>
          <button 
            type="submit"
            disabled={loading}
            className="w-full bg-blue-600 hover:bg-blue-500 text-white font-bold py-4 rounded-2xl transition-all shadow-lg shadow-blue-600/20 active:scale-[0.98]"
          >
            {loading ? 'Authenticating...' : 'Enter Dashboard'}
          </button>
        </form>

        <p className="mt-8 text-center text-xs text-slate-600 uppercase tracking-widest font-semibold">
          Secure Cloud Sync Active
        </p>
      </div>
    </div>
  );
}

const dummyData = [
  { name: '8am', minutes: 12 },
  { name: '9am', minutes: 45 },
  { name: '10am', minutes: 55 },
  { name: '11am', minutes: 30 },
  { name: '12pm', minutes: 20 },
  { name: '1pm', minutes: 10 },
  { name: '2pm', minutes: 40 },
  { name: '3pm', minutes: 50 },
];

const dummyPie = [
  { name: 'Work', value: 400 },
  { name: 'Study', value: 300 },
  { name: 'Entertainment', value: 150 },
  { name: 'Other', value: 100 },
];

export default App;
