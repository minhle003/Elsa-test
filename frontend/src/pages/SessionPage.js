import React from 'react';
import Session from '../components/Session';
import { useLocation } from 'react-router-dom';

const useQuery = () => {
  return new URLSearchParams(useLocation().search);
};

const SessionPage = () => {
  const query = useQuery();
  const sessionId = query.get('session');
  const userId = localStorage.getItem('User-ID');

  return (
    <div>
      <Session sessionId={sessionId} userId={userId} />
    </div>
  );
};

export default SessionPage;