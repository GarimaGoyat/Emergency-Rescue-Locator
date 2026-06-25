import { useCallback, useEffect, useRef, useState } from 'react';
import Layout from '../components/Layout';
import EmergencyMap from '../components/EmergencyMap';
import LoadingSpinner from '../components/LoadingSpinner';
import SOSButton from '../components/SOSButton';
import { useNotification } from '../context/NotificationContext';
import { emergencyAPI } from '../services/api';

function getCurrentPosition() {
  return new Promise((resolve, reject) => {
    if (!navigator.geolocation) {
      reject(new Error('Geolocation is not supported by your browser'));
      return;
    }

    navigator.geolocation.getCurrentPosition(
      (position) => resolve(position),
      (error) => reject(error),
      { enableHighAccuracy: true, timeout: 15000, maximumAge: 0 }
    );
  });
}

export default function EmergencyPage() {
  const [emergency, setEmergency] = useState(null);
  const [loading, setLoading] = useState(true);
  const [sosLoading, setSosLoading] = useState(false);
  const [description, setDescription] = useState('');
  const [latestLocation, setLatestLocation] = useState(null);
  const [locationError, setLocationError] = useState(null);
  const intervalRef = useRef(null);
  const { success, error, info } = useNotification();

  const fetchActive = useCallback(async () => {
    try {
      const res = await emergencyAPI.getActive();
      setEmergency(res.data.data);
      if (res.data.data) {
        const locRes = await emergencyAPI.getLatestLocation(res.data.data.id);
        setLatestLocation(locRes.data.data);
      }
    } catch {
      setEmergency(null);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchActive();
  }, [fetchActive]);

  const sendLocationUpdate = useCallback(
    async (emergencyId) => {
      try {
        const position = await getCurrentPosition();
        const { latitude, longitude, accuracy } = position.coords;

        const res = await emergencyAPI.updateLocation(emergencyId, {
          latitude,
          longitude,
          accuracy,
        });

        setLatestLocation(res.data.data);
        setLocationError(null);
      } catch (err) {
        const message =
          err.message || err.response?.data?.error || 'Failed to update location';
        setLocationError(message);
      }
    },
    []
  );

  useEffect(() => {
    if (!emergency || emergency.status !== 'active') {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
      return;
    }

    sendLocationUpdate(emergency.id);
    intervalRef.current = setInterval(() => {
      sendLocationUpdate(emergency.id);
    }, 10000);

    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    };
  }, [emergency, sendLocationUpdate]);

  const handleSOS = async () => {
    setSosLoading(true);
    try {
      const position = await getCurrentPosition();
      const { latitude, longitude, accuracy } = position.coords;

      const res = await emergencyAPI.create({
        description: description || 'Emergency SOS alert',
        latitude,
        longitude,
        address: `${latitude.toFixed(6)}, ${longitude.toFixed(6)}`,
      });

      setEmergency(res.data.data);
      setLatestLocation({
        latitude,
        longitude,
        accuracy,
        recorded_at: new Date().toISOString(),
      });
      success('Emergency SOS sent! Help is on the way.');
      info('Your location is being tracked every 10 seconds.');
    } catch (err) {
      error(
        err.response?.data?.error ||
          err.message ||
          'Failed to send SOS. Please check location permissions.'
      );
    } finally {
      setSosLoading(false);
    }
  };

  const handleCancel = async () => {
    if (!emergency) return;
    try {
      await emergencyAPI.cancel(emergency.id);
      setEmergency(null);
      setLatestLocation(null);
      success('Emergency cancelled successfully.');
    } catch (err) {
      error(err.response?.data?.error || 'Failed to cancel emergency');
    }
  };

  const mapCenter = latestLocation
    ? [latestLocation.latitude, latestLocation.longitude]
    : emergency
    ? [emergency.latitude, emergency.longitude]
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
          <h1 className="text-2xl font-bold text-slate-900 sm:text-3xl">Emergency SOS</h1>
          <p className="mt-1 text-slate-600">
            Send an emergency alert with your real-time location
          </p>
        </div>

        {!emergency ? (
          <div className="card flex flex-col items-center py-8">
            <div className="mb-6 w-full max-w-md">
              <label htmlFor="description" className="mb-1.5 block text-sm font-medium text-slate-700">
                Description (optional)
              </label>
              <textarea
                id="description"
                rows={3}
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                className="input-field resize-none"
                placeholder="Describe your emergency situation..."
              />
            </div>
            <SOSButton onSOS={handleSOS} loading={sosLoading} />
          </div>
        ) : (
          <div className="space-y-6">
            <div className="card border-l-4 border-l-red-500">
              <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
                <div>
                  <span className="inline-flex animate-pulse items-center rounded-full bg-red-100 px-3 py-1 text-xs font-bold text-red-700">
                    ● ACTIVE EMERGENCY
                  </span>
                  <h2 className="mt-3 text-xl font-bold text-slate-900">
                    Help is being dispatched
                  </h2>
                  <p className="mt-1 text-slate-600">
                    {emergency.description || 'Emergency SOS alert'}
                  </p>
                  <p className="mt-2 text-sm text-slate-500">
                    Started: {new Date(emergency.created_at).toLocaleString()}
                  </p>
                  {latestLocation && (
                    <p className="mt-1 text-sm text-slate-500">
                      Last update: {new Date(latestLocation.recorded_at).toLocaleTimeString()}
                      {latestLocation.accuracy &&
                        ` (±${Math.round(latestLocation.accuracy)}m)`}
                    </p>
                  )}
                  {locationError && (
                    <p className="mt-2 text-sm text-red-600">{locationError}</p>
                  )}
                </div>
                <button onClick={handleCancel} className="btn-secondary shrink-0">
                  Cancel Emergency
                </button>
              </div>
            </div>

            <EmergencyMap
              emergencies={[emergency]}
              center={mapCenter}
              zoom={15}
              height="350px"
            />
          </div>
        )}
      </div>
    </Layout>
  );
}
