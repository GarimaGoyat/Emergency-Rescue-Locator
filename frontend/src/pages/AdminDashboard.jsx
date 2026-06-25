import { useCallback, useEffect, useState } from 'react';
import Layout from '../components/Layout';
import EmergencyMap from '../components/EmergencyMap';
import LoadingSpinner from '../components/LoadingSpinner';
import StatsCard from '../components/StatsCard';
import { useNotification } from '../context/NotificationContext';
import { adminAPI } from '../services/api';

export default function AdminDashboard() {
  const [emergencies, setEmergencies] = useState([]);
  const [stats, setStats] = useState(null);
  const [search, setSearch] = useState('');
  const [selected, setSelected] = useState(null);
  const [details, setDetails] = useState(null);
  const [loading, setLoading] = useState(true);
  const [searchLoading, setSearchLoading] = useState(false);
  const [resolving, setResolving] = useState(false);
  const { success, error } = useNotification();

  const fetchData = useCallback(async () => {
    try {
      const [emergenciesRes, statsRes] = await Promise.all([
        adminAPI.listActive(),
        adminAPI.stats(),
      ]);
      setEmergencies(emergenciesRes.data.data || []);
      setStats(statsRes.data.data);
    } catch {
      error('Failed to load admin data');
    } finally {
      setLoading(false);
    }
  }, [error]);

  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 15000);
    return () => clearInterval(interval);
  }, [fetchData]);

  const handleSearch = async (e) => {
    e.preventDefault();
    setSearchLoading(true);
    try {
      const res = await adminAPI.search({ q: search, status: 'active' });
      setEmergencies(res.data.data || []);
    } catch {
      error('Search failed');
    } finally {
      setSearchLoading(false);
    }
  };

  const handleSelect = async (emergency) => {
    setSelected(emergency.id);
    try {
      const res = await adminAPI.getDetails(emergency.id);
      setDetails(res.data);
    } catch {
      error('Failed to load emergency details');
    }
  };

  const handleResolve = async (id) => {
    setResolving(true);
    try {
      await adminAPI.resolve(id);
      success('Emergency marked as resolved');
      setSelected(null);
      setDetails(null);
      fetchData();
    } catch (err) {
      error(err.response?.data?.error || 'Failed to resolve emergency');
    } finally {
      setResolving(false);
    }
  };

  const mapCenter =
    emergencies.length > 0
      ? [emergencies[0].latitude, emergencies[0].longitude]
      : [20.5937, 78.9629];

  if (loading) {
    return (
      <Layout>
        <LoadingSpinner className="py-24" size="lg" />
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="space-y-8">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 sm:text-3xl">Admin Dashboard</h1>
          <p className="mt-1 text-slate-600">Monitor and manage active emergency requests</p>
        </div>

        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <StatsCard title="Total Emergencies" value={stats?.total_emergencies} icon="📊" />
          <StatsCard
            title="Active Now"
            value={stats?.active_emergencies}
            icon="🚨"
            color="red"
          />
          <StatsCard
            title="Resolved"
            value={stats?.resolved_emergencies}
            icon="✅"
            color="green"
          />
          <StatsCard
            title="Today"
            value={stats?.today_emergencies}
            icon="📅"
            color="amber"
          />
        </div>

        <form onSubmit={handleSearch} className="flex gap-3">
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search by name, email, description..."
            className="input-field flex-1"
          />
          <button type="submit" disabled={searchLoading} className="btn-primary shrink-0">
            {searchLoading ? 'Searching...' : 'Search'}
          </button>
          <button
            type="button"
            onClick={() => {
              setSearch('');
              fetchData();
            }}
            className="btn-secondary shrink-0"
          >
            Reset
          </button>
        </form>

        <div className="grid gap-6 lg:grid-cols-2">
          <div className="space-y-4">
            <h2 className="text-lg font-semibold text-slate-900">
              Active Emergencies ({emergencies.length})
            </h2>

            {emergencies.length === 0 ? (
              <div className="card text-center text-slate-500">
                No active emergencies at this time
              </div>
            ) : (
              <div className="space-y-3">
                {emergencies.map((emergency) => (
                  <button
                    key={emergency.id}
                    onClick={() => handleSelect(emergency)}
                    className={`card w-full text-left transition hover:shadow-md ${
                      selected === emergency.id ? 'ring-2 ring-rescue-500' : ''
                    }`}
                  >
                    <div className="flex items-start justify-between gap-3">
                      <div>
                        <p className="font-semibold text-slate-900">
                          {emergency.user?.name || 'Unknown User'}
                        </p>
                        <p className="text-sm text-slate-600">
                          {emergency.description || 'No description'}
                        </p>
                        <p className="mt-1 text-xs text-slate-500">
                          {new Date(emergency.created_at).toLocaleString()}
                        </p>
                      </div>
                      <span className="shrink-0 rounded-full bg-red-100 px-2 py-0.5 text-xs font-semibold text-red-700">
                        Active
                      </span>
                    </div>
                  </button>
                ))}
              </div>
            )}
          </div>

          <div className="space-y-4">
            <EmergencyMap
              emergencies={emergencies}
              center={mapCenter}
              zoom={emergencies.length ? 12 : 5}
              height="300px"
              onMarkerClick={handleSelect}
              selectedId={selected}
            />

            {details && (
              <div className="card space-y-4">
                <h3 className="text-lg font-semibold text-slate-900">Emergency Details</h3>
                <div className="grid gap-2 text-sm">
                  <p>
                    <span className="font-medium text-slate-700">User:</span>{' '}
                    {details.data.user?.name} ({details.data.user?.email})
                  </p>
                  <p>
                    <span className="font-medium text-slate-700">Phone:</span>{' '}
                    {details.data.user?.phone || 'N/A'}
                  </p>
                  <p>
                    <span className="font-medium text-slate-700">Description:</span>{' '}
                    {details.data.description || 'N/A'}
                  </p>
                  <p>
                    <span className="font-medium text-slate-700">Location:</span>{' '}
                    {details.latest_location
                      ? `${details.latest_location.latitude.toFixed(6)}, ${details.latest_location.longitude.toFixed(6)}`
                      : `${details.data.latitude.toFixed(6)}, ${details.data.longitude.toFixed(6)}`}
                  </p>
                  <p>
                    <span className="font-medium text-slate-700">Location updates:</span>{' '}
                    {details.location_history?.length || 0}
                  </p>
                </div>
                <button
                  onClick={() => handleResolve(details.data.id)}
                  disabled={resolving}
                  className="btn-primary w-full"
                >
                  {resolving ? 'Resolving...' : 'Mark as Resolved'}
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </Layout>
  );
}
