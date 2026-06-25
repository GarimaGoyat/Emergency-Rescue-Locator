import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import Layout from '../components/Layout';
import LoadingSpinner from '../components/LoadingSpinner';
import { useAuth } from '../context/AuthContext';
import { emergencyAPI } from '../services/api';

export default function Dashboard() {
  const { user, isAdmin } = useAuth();
  const [activeEmergency, setActiveEmergency] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    emergencyAPI
      .getActive()
      .then((res) => setActiveEmergency(res.data.data))
      .catch(() => setActiveEmergency(null))
      .finally(() => setLoading(false));
  }, []);

  return (
    <Layout>
      <div className="space-y-8">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 sm:text-3xl">
            Welcome, {user?.name}
          </h1>
          <p className="mt-1 text-slate-600">
            Emergency Rescue Locator — stay safe and connected
          </p>
        </div>

        {loading ? (
          <LoadingSpinner className="py-12" />
        ) : activeEmergency ? (
          <div className="card border-l-4 border-l-primary-500">
            <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
              <div>
                <span className="inline-flex items-center rounded-full bg-red-100 px-2.5 py-0.5 text-xs font-semibold text-red-700">
                  Active Emergency
                </span>
                <h2 className="mt-2 text-lg font-semibold text-slate-900">
                  Emergency in progress
                </h2>
                <p className="text-sm text-slate-600">
                  {activeEmergency.description || 'SOS alert active — location being tracked'}
                </p>
                <p className="mt-1 text-xs text-slate-500">
                  Started: {new Date(activeEmergency.created_at).toLocaleString()}
                </p>
              </div>
              <Link to="/emergency" className="btn-danger">
                View Emergency
              </Link>
            </div>
          </div>
        ) : (
          <div className="card text-center">
            <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-green-100 text-2xl">
              ✓
            </div>
            <h2 className="text-lg font-semibold text-slate-900">No active emergencies</h2>
            <p className="mt-1 text-sm text-slate-600">
              You are safe. Press SOS if you need immediate help.
            </p>
            <Link to="/emergency" className="btn-danger mt-4 inline-flex">
              Go to SOS
            </Link>
          </div>
        )}

        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <Link to="/emergency" className="card transition hover:shadow-md">
            <span className="text-2xl">🆘</span>
            <h3 className="mt-2 font-semibold text-slate-900">Emergency SOS</h3>
            <p className="mt-1 text-sm text-slate-600">
              Send an emergency alert with your live location
            </p>
          </Link>

          <div className="card">
            <span className="text-2xl">📍</span>
            <h3 className="mt-2 font-semibold text-slate-900">Location Tracking</h3>
            <p className="mt-1 text-sm text-slate-600">
              Your location is shared every 10 seconds during active emergencies
            </p>
          </div>

          {isAdmin && (
            <Link to="/admin" className="card transition hover:shadow-md">
              <span className="text-2xl">🛡️</span>
              <h3 className="mt-2 font-semibold text-slate-900">Admin Dashboard</h3>
              <p className="mt-1 text-sm text-slate-600">
                Manage and resolve active emergency requests
              </p>
            </Link>
          )}
        </div>
      </div>
    </Layout>
  );
}
