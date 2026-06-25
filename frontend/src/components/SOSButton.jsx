import { useState } from 'react';
import LoadingSpinner from './LoadingSpinner';

export default function SOSButton({ onSOS, disabled, loading }) {
  const [confirming, setConfirming] = useState(false);

  const handleClick = () => {
    if (!confirming) {
      setConfirming(true);
      return;
    }
    onSOS();
    setConfirming(false);
  };

  return (
    <div className="flex flex-col items-center gap-6">
      <button
        onClick={handleClick}
        disabled={disabled || loading}
        className={`relative flex h-48 w-48 items-center justify-center rounded-full text-3xl font-black text-white shadow-2xl transition-all sm:h-56 sm:w-56 sm:text-4xl ${
          confirming
            ? 'animate-pulse bg-primary-700 ring-8 ring-primary-200'
            : 'bg-primary-600 hover:scale-105 hover:bg-primary-700'
        } disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:scale-100`}
      >
        {loading ? (
          <LoadingSpinner size="lg" />
        ) : confirming ? (
          'CONFIRM'
        ) : (
          'SOS'
        )}
        {!confirming && !loading && (
          <span className="absolute inset-0 animate-ping rounded-full bg-primary-500 opacity-20" />
        )}
      </button>

      <p className="max-w-sm text-center text-sm text-slate-600">
        {confirming
          ? 'Tap again to confirm emergency SOS. Your location will be shared with rescue teams.'
          : 'Press the SOS button to send an emergency alert with your current location.'}
      </p>

      {confirming && (
        <button
          onClick={() => setConfirming(false)}
          className="btn-secondary"
        >
          Cancel
        </button>
      )}
    </div>
  );
}
