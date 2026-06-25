export default function StatsCard({ title, value, icon, color = 'rescue' }) {
  const colors = {
    rescue: 'bg-rescue-50 text-rescue-700',
    red: 'bg-red-50 text-red-700',
    green: 'bg-green-50 text-green-700',
    amber: 'bg-amber-50 text-amber-700',
  };

  return (
    <div className="card flex items-center gap-4">
      <div className={`flex h-12 w-12 items-center justify-center rounded-xl text-2xl ${colors[color]}`}>
        {icon}
      </div>
      <div>
        <p className="text-sm font-medium text-slate-500">{title}</p>
        <p className="text-2xl font-bold text-slate-900">{value ?? 0}</p>
      </div>
    </div>
  );
}
