import { useEffect, useRef } from 'react';
import { MapContainer, TileLayer, Marker, Popup, useMap } from 'react-leaflet';
import L from 'leaflet';

const defaultIcon = new L.Icon({
  iconUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
  iconRetinaUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png',
  shadowUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png',
  iconSize: [25, 41],
  iconAnchor: [12, 41],
  popupAnchor: [1, -34],
  shadowSize: [41, 41],
});

const activeIcon = new L.Icon({
  iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-red.png',
  iconRetinaUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-red.png',
  shadowUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png',
  iconSize: [25, 41],
  iconAnchor: [12, 41],
  popupAnchor: [1, -34],
  shadowSize: [41, 41],
});

function MapUpdater({ center, zoom }) {
  const map = useMap();
  const prevCenter = useRef(center);

  useEffect(() => {
    if (
      center &&
      (prevCenter.current?.[0] !== center[0] || prevCenter.current?.[1] !== center[1])
    ) {
      map.setView(center, zoom);
      prevCenter.current = center;
    }
  }, [center, zoom, map]);

  return null;
}

export default function EmergencyMap({
  emergencies = [],
  center = [20.5937, 78.9629],
  zoom = 5,
  height = '400px',
  onMarkerClick,
  selectedId,
}) {
  return (
    <div className="overflow-hidden rounded-xl border border-slate-200" style={{ height }}>
      <MapContainer center={center} zoom={zoom} scrollWheelZoom style={{ height: '100%', width: '100%' }}>
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />
        <MapUpdater center={center} zoom={zoom} />

        {emergencies.map((emergency) => (
          <Marker
            key={emergency.id}
            position={[emergency.latitude, emergency.longitude]}
            icon={emergency.status === 'active' ? activeIcon : defaultIcon}
            eventHandlers={{
              click: () => onMarkerClick?.(emergency),
            }}
          >
            <Popup>
              <div className="space-y-1 text-sm">
                <p className="font-semibold">{emergency.user?.name || 'Unknown'}</p>
                <p className="text-slate-600">{emergency.description || 'No description'}</p>
                <p className="text-xs capitalize text-slate-500">Status: {emergency.status}</p>
                {emergency.address && (
                  <p className="text-xs text-slate-500">{emergency.address}</p>
                )}
              </div>
            </Popup>
          </Marker>
        ))}
      </MapContainer>
    </div>
  );
}
